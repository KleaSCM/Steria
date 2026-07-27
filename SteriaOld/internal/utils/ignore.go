package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// IgnorePattern represents a pattern to ignore
type IgnorePattern struct {
	Pattern string
	IsDir   bool
}

// LoadIgnorePatterns loads patterns from .steriaignore file
func LoadIgnorePatterns(repoPath string) ([]IgnorePattern, error) {
	ignoreFile := filepath.Join(repoPath, ".steriaignore")

	file, err := os.Open(ignoreFile)
	if err != nil {
		if os.IsNotExist(err) {
			// No .steriaignore file, return empty patterns
			return []IgnorePattern{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []IgnorePattern

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Determine if it's a directory pattern
		isDir := strings.HasSuffix(line, "/")
		if isDir {
			line = strings.TrimSuffix(line, "/")
		}

		patterns = append(patterns, IgnorePattern{
			Pattern: line,
			IsDir:   isDir,
		})
	}

	return patterns, scanner.Err()
}

// ShouldIgnore checks if a file or directory should be ignored
func ShouldIgnore(path string, patterns []IgnorePattern) bool {
	// Always ignore .steria directory
	if filepath.Base(path) == ".steria" {
		return true
	}

	// Check against ignore patterns
	for _, pattern := range patterns {
		if matchesPattern(path, pattern) {
			return true
		}
	}

	return false
}

// matchesPattern checks if a path matches an ignore pattern
func matchesPattern(path string, pattern IgnorePattern) bool {
	// Simple pattern matching - can be enhanced with glob patterns later
	pathBase := filepath.Base(path)

	// Exact match
	if pathBase == pattern.Pattern {
		return true
	}

	// Directory pattern
	if pattern.IsDir && filepath.Dir(path) == pattern.Pattern {
		return true
	}

	// Suffix match (e.g., "*.log")
	if strings.HasPrefix(pattern.Pattern, "*.") {
		ext := strings.TrimPrefix(pattern.Pattern, "*")
		if strings.HasSuffix(pathBase, ext) {
			return true
		}
	}

	// Prefix match (e.g., "temp*")
	if strings.HasSuffix(pattern.Pattern, "*") {
		prefix := strings.TrimSuffix(pattern.Pattern, "*")
		if strings.HasPrefix(pathBase, prefix) {
			return true
		}
	}

	return false
}

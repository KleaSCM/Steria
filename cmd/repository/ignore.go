// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: ignore.go
// Description: Implements the 'steria ignore' command for interactive editing of .steriaignore file.

package repository

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"steria/internal/metrics"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewIgnoreCmd creates the 'ignore' command for Steria
func NewIgnoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ignore [pattern]",
		Short: "Manage .steriaignore file",
		Long:  "Interactive command to add, remove, and view ignore patterns in .steriaignore file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var pattern string
			if len(args) > 0 {
				pattern = args[0]
			}
			return runIgnore(pattern)
		},
	}
	return cmd
}

// runIgnore manages the .steriaignore file interactively
func runIgnore(pattern string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	fmt.Printf("%s Managing .steriaignore file with interactive interface...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	ignorePath := filepath.Join(cwd, ".steriaignore")

	// If a pattern is provided, add it directly
	if pattern != "" {
		return addIgnorePattern(ignorePath, pattern)
	}

	// Show current ignore patterns
	fmt.Printf("\n%s Current .steriaignore patterns:\n", magenta("📋"))
	patterns, err := loadIgnorePatterns(ignorePath)
	if err != nil {
		fmt.Printf("%s No .steriaignore file found or error reading file\n", yellow("⚠️"))
	} else if len(patterns) == 0 {
		fmt.Printf("%s No ignore patterns found\n", yellow("⚠️"))
	} else {
		for i, pattern := range patterns {
			fmt.Printf("  %d. %s\n", i+1, pattern)
		}
	}

	// Show interactive menu
	fmt.Printf("\n%s Interactive .steriaignore Manager\n", cyan("🎛️"))
	fmt.Printf("  1. %s Add a new ignore pattern\n", green("➕"))
	fmt.Printf("  2. %s Remove an ignore pattern\n", red("🗑️"))
	fmt.Printf("  3. %s View all patterns\n", cyan("👁️"))
	fmt.Printf("  4. %s Exit\n", yellow("🚪"))

	// For now, we'll show a helpful message since we can't do interactive input
	fmt.Printf("\n%s Interactive mode requires terminal input handling.\n", yellow("⚠️"))
	fmt.Printf("%s For now, use: 'steria ignore <pattern>' to add patterns directly.\n", cyan("💡"))
	fmt.Printf("%s Example: 'steria ignore *.log'\n", cyan("💡"))
	fmt.Printf("%s Example: 'steria ignore temp/'\n", cyan("💡"))

	return nil
}

// addIgnorePattern adds a new pattern to the .steriaignore file
func addIgnorePattern(ignorePath, pattern string) error {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Load existing patterns
	patterns, err := loadIgnorePatterns(ignorePath)
	if err != nil {
		patterns = []string{}
	}

	// Check if pattern already exists
	for _, existingPattern := range patterns {
		if existingPattern == pattern {
			fmt.Printf("%s Pattern '%s' already exists in .steriaignore\n", yellow("⚠️"), pattern)
			return nil
		}
	}

	// Add new pattern
	patterns = append(patterns, pattern)

	// Write back to file
	err = writeIgnorePatterns(ignorePath, patterns)
	if err != nil {
		return fmt.Errorf("failed to write .steriaignore file: %w", err)
	}

	fmt.Printf("%s Successfully added pattern: %s\n", green("✅"), pattern)
	fmt.Printf("%s Total patterns in .steriaignore: %d\n", cyan("📊"), len(patterns))

	// Update performance metrics
	metrics.GlobalMetrics.IncrementFilesProcessed(1)

	return nil
}

// loadIgnorePatterns loads patterns from .steriaignore file
func loadIgnorePatterns(ignorePath string) ([]string, error) {
	file, err := os.Open(ignorePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return patterns, scanner.Err()
}

// writeIgnorePatterns writes patterns to .steriaignore file
func writeIgnorePatterns(ignorePath string, patterns []string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(ignorePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(ignorePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header comment
	_, err = fmt.Fprintf(file, "# Steria ignore patterns\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "# Lines starting with # are comments\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "# Add patterns to ignore files and directories\n\n")
	if err != nil {
		return err
	}

	// Write patterns
	for _, pattern := range patterns {
		_, err = fmt.Fprintf(file, "%s\n", pattern)
		if err != nil {
			return err
		}
	}

	return nil
}

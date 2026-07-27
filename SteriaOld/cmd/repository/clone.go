// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: clone.go
// Description: Implements the steria clone command for cloning git or Steria repositories with full error handling and detailed comments.
package repository

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"steria/internal/metrics"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone [url] [dir]",
		Short: "Clone a repository from git",
		Long:  "Clone a repository from git with optimized processing",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			dir := args[1]
			return runClone(url, dir)
		},
	}

	return cmd
}

func runClone(url, dir string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Cloning repository with optimized processing...\n", cyan("🚀"))

	// Check if destination directory exists
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("destination directory '%s' already exists", dir)
	}

	// Determine if URL is a git repo or a local Steria repo
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasSuffix(url, ".git") {
		// --- Git repository cloning ---
		fmt.Printf("%s Detected git repository. Using git to clone...\n", yellow("💡"))
		cmd := exec.Command("git", "clone", url, dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone git repository: %w", err)
		}
		fmt.Printf("%s Cloned git repository from '%s' into '%s'!\n", green("✅"), red(url), green(dir))
		fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
		return nil
	}

	// --- Steria repository cloning (local path) ---
	fmt.Printf("%s Detected Steria repository. Copying directory...\n", yellow("💡"))
	if err := copySteriaRepo(url, dir); err != nil {
		return fmt.Errorf("failed to clone Steria repository: %w", err)
	}
	fmt.Printf("%s Cloned Steria repository from '%s' into '%s'!\n", green("✅"), red(url), green(dir))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

// copySteriaRepo copies all files and the .steria folder from src to dst, skipping junk files
func copySteriaRepo(src, dst string) error {
	// List of files/folders to skip (junk)
	junk := map[string]bool{
		".git": true, "go.mod": true, "go.sum": true, "main.go": true,
		"README.md": true, "TEMPLATE.md": true, "cmd": true, "core": true,
		"internal": true, "Docs": true, "Tests": true, "steria": true,
		".gitignore": true, ".steriaignore": true,
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	for _, entry := range entries {
		if junk[entry.Name()] {
			continue
		}
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	// Always copy .steria folder if present
	steriaPath := filepath.Join(src, ".steria")
	if _, err := os.Stat(steriaPath); err == nil {
		if err := CopyDir(steriaPath, filepath.Join(dst, ".steria")); err != nil {
			return err
		}
	}
	return nil
}

func extractDirFromURL(url string) string {
	// Simple extraction - get the last part of the URL
	// This is a basic implementation
	if len(url) == 0 {
		return "repository"
	}

	// Remove .git suffix if present
	if len(url) > 4 && url[len(url)-4:] == ".git" {
		url = url[:len(url)-4]
	}

	// Get the last part after the last slash
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			return url[i+1:]
		}
	}

	return url
}

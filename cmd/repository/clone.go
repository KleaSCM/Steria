package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone [url] [directory]",
		Short: "Clone a repository from git",
		Long:  "Clone a git repository and convert it to Steria format",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			dir := ""
			if len(args) > 1 {
				dir = args[1]
			}
			return runClone(url, dir)
		},
	}

	return cmd
}

func runClone(url, dir string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("%s Cloning %s...\n", cyan("🔄"), url)

	// For now, this is a placeholder
	// We'll implement actual git cloning later
	if dir == "" {
		// Extract directory name from URL
		dir = extractDirFromURL(url)
	}

	// Create directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Initialize KleaSCM repository
	repo, err := storage.LoadOrInitRepo(dir)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	// Set remote URL
	repo.RemoteURL = url
	remotePath := filepath.Join(dir, ".steria", "remote")
	if err := os.WriteFile(remotePath, []byte(url), 0644); err != nil {
		return fmt.Errorf("failed to save remote URL: %w", err)
	}

	fmt.Printf("%s Successfully cloned to %s\n", green("✅"), dir)
	fmt.Printf("%s Repository initialized with Steria\n", green("✨"))

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

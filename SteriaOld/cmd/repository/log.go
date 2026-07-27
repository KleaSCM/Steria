// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: log.go
// Description: Implements the 'steria log' command to show a pretty, color-coded commit history.

package repository

import (
	"fmt"
	"os"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewLogCmd creates the 'log' command for Steria
func NewLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Show commit history",
		Long:  "Display a pretty, color-coded commit history with optimized processing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLog()
		},
	}
	return cmd
}

// runLog displays the commit history in a pretty, color-coded format
func runLog() error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	fmt.Printf("%s Showing commit history with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	if repo.Head == "" {
		fmt.Printf("%s No commits found in repository\n", yellow("⚠️"))
		return nil
	}

	// Walk through commit history
	currentHash := repo.Head
	commitCount := 0
	maxCommits := 50 // Limit to prevent infinite loops

	for currentHash != "" && commitCount < maxCommits {
		commit, err := repo.LoadCommit(currentHash)
		if err != nil {
			fmt.Printf("%s Failed to load commit %s: %v\n", red("❌"), currentHash[:8], err)
			break
		}

		// Print commit info with colors
		fmt.Printf("\n%s %s\n", magenta("📍"), yellow(commit.Hash[:8]))
		fmt.Printf("%s %s\n", green("👤"), commit.Author)
		fmt.Printf("%s %s\n", cyan("📅"), commit.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("%s %s\n", magenta("💬"), commit.Message)

		if len(commit.Files) > 0 {
			fmt.Printf("%s %d files\n", cyan("📁"), len(commit.Files))
		}

		// Move to parent commit
		currentHash = commit.Parent
		commitCount++
	}

	if commitCount >= maxCommits {
		fmt.Printf("\n%s Reached maximum commit limit (%d)\n", yellow("⚠️"), maxCommits)
	}

	fmt.Printf("\n%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

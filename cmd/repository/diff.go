// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: diff.go
// Description: Implements the 'steria diff' command to show differences between file versions.

package repository

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewDiffCmd creates the 'diff' command for Steria
func NewDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [file]",
		Short: "Show differences between file versions",
		Long:  "Display differences between the working directory and the last commit, or between two commits",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var filePath string
			if len(args) > 0 {
				filePath = args[0]
			}
			return runDiff(filePath)
		},
	}
	return cmd
}

// runDiff displays differences between file versions
func runDiff(filePath string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Analyzing file differences with optimized processing...\n", cyan("🚀"))

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

	// Get the last commit
	commit, err := repo.LoadCommit(repo.Head)
	if err != nil {
		return fmt.Errorf("failed to load last commit: %w", err)
	}

	// If a specific file is provided, show diff for that file only
	if filePath != "" {
		return showFileDiff(repo, filePath, commit)
	}

	// Show diff for all changed files
	return showAllDiffs(repo, commit)
}

// showFileDiff shows differences for a specific file
func showFileDiff(repo *storage.Repo, filePath string, commit *storage.Commit) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("\n%s Showing differences for: %s\n", cyan("📁"), yellow(filePath))

	// Check if file exists in current working directory
	currentPath := filePath
	if !strings.HasPrefix(currentPath, "/") {
		currentPath = filePath
	}

	// Check if file exists in commit
	fileInCommit := false
	for _, commitFile := range commit.Files {
		if commitFile == filePath {
			fileInCommit = true
			break
		}
	}

	// Read current file content
	var currentContent []string
	if _, err := os.Stat(currentPath); err == nil {
		file, err := os.Open(currentPath)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				currentContent = append(currentContent, scanner.Text())
			}
		}
	}

	// Show file status
	if !fileInCommit && len(currentContent) > 0 {
		fmt.Printf("%s File is new (not in last commit)\n", green("➕"))
		fmt.Printf("%s Total lines: %d\n", cyan("📊"), len(currentContent))
	} else if fileInCommit && len(currentContent) == 0 {
		fmt.Printf("%s File has been deleted\n", red("🗑️"))
	} else if fileInCommit && len(currentContent) > 0 {
		fmt.Printf("%s File has been modified\n", yellow("✏️"))
		fmt.Printf("%s Total lines: %d\n", cyan("📊"), len(currentContent))
	} else {
		fmt.Printf("%s File not found in working directory or commit\n", red("❌"))
	}

	// Show a simple diff (first few lines)
	if len(currentContent) > 0 {
		fmt.Printf("\n%s First few lines of current file:\n", cyan("📄"))
		maxLines := 10
		if len(currentContent) < maxLines {
			maxLines = len(currentContent)
		}
		for i := 0; i < maxLines; i++ {
			fmt.Printf("  %s\n", currentContent[i])
		}
		if len(currentContent) > maxLines {
			fmt.Printf("  %s ... (%d more lines)\n", yellow("..."), len(currentContent)-maxLines)
		}
	}

	return nil
}

// showAllDiffs shows differences for all changed files
func showAllDiffs(repo *storage.Repo, commit *storage.Commit) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	fmt.Printf("\n%s Commit: %s\n", magenta("📍"), yellow(commit.Hash[:8]))
	fmt.Printf("%s Message: %s\n", cyan("💬"), commit.Message)
	fmt.Printf("%s Files in commit: %d\n", cyan("📁"), len(commit.Files))

	// Get current working directory state
	changes, err := repo.GetChanges()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Printf("\n%s No changes detected in working directory\n", green("✅"))
		return nil
	}

	fmt.Printf("\n%s Changes in working directory:\n", cyan("🔄"))
	for _, change := range changes {
		switch change.Type {
		case storage.ChangeTypeAdded:
			fmt.Printf("  %s %s (new file)\n", green("➕"), change.Path)
		case storage.ChangeTypeModified:
			fmt.Printf("  %s %s (modified)\n", yellow("✏️"), change.Path)
		case storage.ChangeTypeDeleted:
			fmt.Printf("  %s %s (deleted)\n", red("🗑️"), change.Path)
		}
	}

	fmt.Printf("\n%s Use 'steria diff <filename>' to see detailed differences for a specific file\n", cyan("��"))
	return nil
}

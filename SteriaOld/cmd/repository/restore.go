// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: restore.go
// Description: Implements the 'steria restore' command to restore deleted or previous versions of files.

package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewRestoreCmd creates the 'restore' command for Steria
func NewRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore <file> [commit-hash]",
		Short: "Restore files from previous commits",
		Long:  "Restore deleted or previous versions of files from specific commits or the last commit",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			var commitHash string
			if len(args) > 1 {
				commitHash = args[1]
			}
			return runRestore(filePath, commitHash)
		},
	}
	return cmd
}

// runRestore restores a file from a specific commit or the last commit
func runRestore(filePath, commitHash string) error {
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

	fmt.Printf("%s Starting file restoration with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	if repo.Head == "" {
		return fmt.Errorf("no commits found in repository")
	}

	// Determine which commit to restore from
	targetCommit := repo.Head
	if commitHash != "" {
		targetCommit = commitHash
	}

	fmt.Printf("%s Restoring from commit: %s\n", magenta("📍"), yellow(targetCommit[:8]))

	// Load the target commit
	commit, err := repo.LoadCommit(targetCommit)
	if err != nil {
		return fmt.Errorf("failed to load commit %s: %w", targetCommit[:8], err)
	}

	// Check if file exists in the commit
	fileExists := false
	for _, commitFile := range commit.Files {
		if commitFile == filePath {
			fileExists = true
			break
		}
	}

	if !fileExists {
		return fmt.Errorf("file '%s' not found in commit %s", filePath, targetCommit[:8])
	}

	// Check if file exists in current working directory
	currentPath := filepath.Join(cwd, filePath)
	fileExistsCurrent := false
	if _, err := os.Stat(currentPath); err == nil {
		fileExistsCurrent = true
	}

	// Show restoration preview
	fmt.Printf("\n%s File: %s\n", cyan("📁"), yellow(filePath))
	if fileExistsCurrent {
		fmt.Printf("%s Current status: File exists in working directory\n", green("✅"))
		fmt.Printf("%s Action: Will overwrite with version from commit %s\n", yellow("⚠️"), targetCommit[:8])
	} else {
		fmt.Printf("%s Current status: File not found in working directory\n", red("❌"))
		fmt.Printf("%s Action: Will restore from commit %s\n", green("🔄"), targetCommit[:8])
	}

	// Restore the file from the commit's blob
	blobHash, ok := commit.FileBlobs[filePath]
	if !ok {
		return fmt.Errorf("file blob for '%s' not found in commit %s", filePath, targetCommit[:8])
	}
	blobPath := filepath.Join(cwd, ".steria", "objects", "blobs", blobHash)
	blobData, err := os.ReadFile(blobPath)
	if err != nil {
		return fmt.Errorf("failed to read blob for '%s': %w", filePath, err)
	}
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(currentPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	if err := os.WriteFile(currentPath, blobData, 0644); err != nil {
		return fmt.Errorf("failed to write restored file: %w", err)
	}

	fmt.Printf("%s File '%s' restored from commit %s\n", green("✅"), filePath, targetCommit[:8])

	metrics.GlobalMetrics.IncrementFilesProcessed(1)

	return nil
}

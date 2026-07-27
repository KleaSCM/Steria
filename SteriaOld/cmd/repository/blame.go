// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: blame.go
// Description: Implements blame functionality for tracking line-by-line history.

package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type BlameLine struct {
	LineNumber int
	Commit     string
	Author     string
	Timestamp  time.Time
	Content    string
}

func NewBlameCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "blame <file>",
		Short: "Show line-by-line history of a file",
		Long:  "Blame shows who changed each line and when",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return blameFile(args[0])
		},
	}
}

func blameFile(filePath string) error {
	repoPath, _ := os.Getwd()

	// Load repository
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Check if file exists
	fullPath := filepath.Join(repoPath, filePath)
	if _, err := os.Stat(fullPath); err != nil {
		return fmt.Errorf("file %s not found: %w", filePath, err)
	}

	// Get current file content
	currentContent, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get blame information
	blameLines, err := calculateBlame(repo, filePath, string(currentContent))
	if err != nil {
		return fmt.Errorf("failed to calculate blame: %w", err)
	}

	// Display blame
	displayBlame(blameLines)

	return nil
}

func calculateBlame(repo *storage.Repo, filePath, currentContent string) ([]BlameLine, error) {
	lines := strings.Split(currentContent, "\n")
	blameLines := make([]BlameLine, len(lines))

	// Initialize all lines with current commit info
	currentCommit, err := repo.LoadCommit(repo.Head)
	if err != nil {
		return nil, err
	}

	for i, line := range lines {
		blameLines[i] = BlameLine{
			LineNumber: i + 1,
			Commit:     currentCommit.Hash,
			Author:     currentCommit.Author,
			Timestamp:  currentCommit.Timestamp,
			Content:    line,
		}
	}

	// Walk through commit history to find when each line was last modified
	hash := repo.Head
	seen := make(map[string]bool)

	for hash != "" && !seen[hash] {
		seen[hash] = true
		commit, err := repo.LoadCommit(hash)
		if err != nil {
			break
		}

		// Check if this commit modified the file
		if !hasFileInCommit(commit, filePath) {
			hash = commit.Parent
			continue
		}

		// Get the file content from this commit
		commitContent, err := getFileContentFromCommit(repo, commit, filePath)
		if err != nil {
			hash = commit.Parent
			continue
		}

		// Update blame lines for lines that were modified in this commit
		updateBlameLines(blameLines, commit, commitContent)

		hash = commit.Parent
	}

	return blameLines, nil
}

func hasFileInCommit(commit *storage.Commit, filePath string) bool {
	for _, file := range commit.Files {
		if file == filePath {
			return true
		}
	}
	return false
}

func getFileContentFromCommit(repo *storage.Repo, commit *storage.Commit, filePath string) (string, error) {
	// Find the blob hash for this file in this commit
	var blobRef string
	for _, file := range commit.Files {
		if file == filePath {
			blobRef = commit.FileBlobs[file]
			break
		}
	}

	if blobRef == "" {
		return "", fmt.Errorf("file not found in commit")
	}

	// Get blob content
	blobDir := filepath.Join(repo.Path, ".steria", "objects", "blobs")
	store := &storage.LocalBlobStore{Dir: blobDir}
	data, err := storage.ReadFileBlobDecompressed(store, blobRef)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func updateBlameLines(blameLines []BlameLine, commit *storage.Commit, commitContent string) {
	commitLines := strings.Split(commitContent, "\n")

	// Simple line-by-line comparison
	// In a more sophisticated implementation, you'd use diff algorithms
	for i := 0; i < len(blameLines) && i < len(commitLines); i++ {
		if blameLines[i].Content == commitLines[i] {
			// Line content matches, update blame info
			blameLines[i].Commit = commit.Hash
			blameLines[i].Author = commit.Author
			blameLines[i].Timestamp = commit.Timestamp
		}
	}
}

func displayBlame(blameLines []BlameLine) {
	fmt.Println("Blame for file:")
	fmt.Println(strings.Repeat("=", 80))

	for _, line := range blameLines {
		// Format: commit_hash author timestamp line_number: content
		fmt.Printf("%s %s %s %d: %s\n",
			line.Commit[:8],
			line.Author,
			line.Timestamp.Format("2006-01-02 15:04:05"),
			line.LineNumber,
			line.Content)
	}
}

// Alternative display with more detailed information
func displayDetailedBlame(blameLines []BlameLine) {
	fmt.Println("Detailed blame for file:")
	fmt.Println(strings.Repeat("=", 80))

	// Group by commit for better readability
	commitGroups := make(map[string][]BlameLine)
	for _, line := range blameLines {
		commitGroups[line.Commit] = append(commitGroups[line.Commit], line)
	}

	for commitHash, lines := range commitGroups {
		if len(lines) == 0 {
			continue
		}

		// Show commit info
		fmt.Printf("\nCommit: %s\n", commitHash[:8])
		fmt.Printf("Author: %s\n", lines[0].Author)
		fmt.Printf("Date: %s\n", lines[0].Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("Lines: %d-%d\n", lines[0].LineNumber, lines[len(lines)-1].LineNumber)
		fmt.Println(strings.Repeat("-", 40))

		// Show lines
		for _, line := range lines {
			fmt.Printf("%4d: %s\n", line.LineNumber, line.Content)
		}
	}
}

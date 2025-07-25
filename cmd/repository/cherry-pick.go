// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: cherry-pick.go
// Description: Implements cherry-pick functionality for applying specific commits to current branch.

package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"

	"github.com/spf13/cobra"
)

func NewCherryPickCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cherry-pick <commit>",
		Short: "Apply a specific commit to the current branch",
		Long:  "Cherry-pick applies the changes from a specific commit to the current branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cherryPick(args[0])
		},
	}
}

func cherryPick(commitHash string) error {
	repoPath, _ := os.Getwd()

	// Load repository
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Load the commit to cherry-pick
	sourceCommit, err := repo.LoadCommit(commitHash)
	if err != nil {
		return fmt.Errorf("commit %s not found: %w", commitHash, err)
	}

	// Check if commit is already in current branch
	if isCommitInBranch(repo, sourceCommit.Hash) {
		return fmt.Errorf("commit %s is already in current branch", commitHash[:8])
	}

	// Note: We don't need current state for cherry-pick since we're applying changes directly

	// Get the parent commit of the source commit
	var parentCommit *storage.Commit
	if sourceCommit.Parent != "" {
		parentCommit, err = repo.LoadCommit(sourceCommit.Parent)
		if err != nil {
			return fmt.Errorf("failed to load parent commit: %w", err)
		}
	}

	// Calculate the changes in the source commit
	changes := calculateCommitChanges(parentCommit, sourceCommit)

	// Apply changes to working directory
	if err := applyChangesToWorkingDir(repoPath, changes); err != nil {
		return fmt.Errorf("failed to apply changes: %w", err)
	}

	// Create new commit with cherry-picked changes
	newCommit, err := repo.CreateCommit(
		fmt.Sprintf("cherry-pick: %s", sourceCommit.Message),
		"KleaSCM",
	)
	if err != nil {
		return fmt.Errorf("failed to create cherry-pick commit: %w", err)
	}

	fmt.Printf("Cherry-picked commit %s (%s)\n", commitHash[:8], sourceCommit.Message)
	fmt.Printf("New commit: %s\n", newCommit.Hash[:8])
	return nil
}

func isCommitInBranch(repo *storage.Repo, commitHash string) bool {
	// Walk from HEAD to root to check if commit is in current branch
	hash := repo.Head
	for hash != "" {
		if hash == commitHash {
			return true
		}
		commit, err := repo.LoadCommit(hash)
		if err != nil {
			break
		}
		hash = commit.Parent
	}
	return false
}

func calculateCommitChanges(parent, commit *storage.Commit) map[string]string {
	changes := make(map[string]string)

	if parent == nil {
		// Initial commit - all files are new
		for file, blob := range commit.FileBlobs {
			changes[file] = blob
		}
		return changes
	}

	// Find files that were added or modified
	for file, blob := range commit.FileBlobs {
		if parentBlob, exists := parent.FileBlobs[file]; !exists || parentBlob != blob {
			changes[file] = blob
		}
	}

	// Find files that were deleted
	for file := range parent.FileBlobs {
		if _, exists := commit.FileBlobs[file]; !exists {
			changes[file] = "" // Empty string indicates deletion
		}
	}

	return changes
}

func applyChangesToWorkingDir(repoPath string, changes map[string]string) error {
	for file, blobRef := range changes {
		filePath := filepath.Join(repoPath, file)

		if blobRef == "" {
			// File was deleted
			if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to delete file %s: %w", file, err)
			}
			continue
		}

		// Ensure directory exists
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", file, err)
		}

		// Get blob data
		blobDir := filepath.Join(repoPath, ".steria", "objects", "blobs")
		store := &storage.LocalBlobStore{Dir: blobDir}
		data, err := storage.ReadFileBlobDecompressed(store, blobRef)
		if err != nil {
			return fmt.Errorf("failed to read blob for %s: %w", file, err)
		}

		// Write file
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", file, err)
		}
	}

	return nil
}

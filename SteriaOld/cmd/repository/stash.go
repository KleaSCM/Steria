// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: stash.go
// Description: Implements stash management for temporarily saving uncommitted changes.

package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"
	"time"

	"github.com/spf13/cobra"
)

type Stash struct {
	ID        string            `json:"id"`
	Message   string            `json:"message"`
	Branch    string            `json:"branch"`
	Timestamp time.Time         `json:"timestamp"`
	Files     map[string]string `json:"files"` // file path -> blob hash
}

func NewStashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stash",
		Short: "Manage stashes",
		Long:  "Save, list, apply, and drop stashes for temporary storage of changes",
	}
	cmd.AddCommand(newStashSaveCmd())
	cmd.AddCommand(newStashListCmd())
	cmd.AddCommand(newStashApplyCmd())
	cmd.AddCommand(newStashDropCmd())
	cmd.AddCommand(newStashPopCmd())
	return cmd
}

func newStashSaveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save [message]",
		Short: "Save current changes to stash",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := "WIP"
			if len(args) > 0 {
				message = args[0]
			}
			return stashSave(message)
		},
	}
}

func newStashListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all stashes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return stashList()
		},
	}
}

func newStashApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "apply [stash]",
		Short: "Apply a stash (keeps the stash)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stashID := "0" // Default to latest
			if len(args) > 0 {
				stashID = args[0]
			}
			return stashApply(stashID)
		},
	}
}

func newStashDropCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drop [stash]",
		Short: "Drop a stash",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stashID := "0" // Default to latest
			if len(args) > 0 {
				stashID = args[0]
			}
			return stashDrop(stashID)
		},
	}
}

func newStashPopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pop [stash]",
		Short: "Apply and drop a stash",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			stashID := "0" // Default to latest
			if len(args) > 0 {
				stashID = args[0]
			}
			return stashPop(stashID)
		},
	}
}

func stashSave(message string) error {
	repoPath, _ := os.Getwd()

	// Load repository
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Get current changes
	changes, err := repo.GetChanges()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if len(changes) == 0 {
		return fmt.Errorf("no changes to stash")
	}

	// Create stash
	stash := &Stash{
		ID:        fmt.Sprintf("%d", time.Now().Unix()),
		Message:   message,
		Branch:    repo.Branch,
		Timestamp: time.Now(),
		Files:     make(map[string]string),
	}

	// Save changed files to blobs
	for _, change := range changes {
		if change.Type != storage.ChangeTypeDeleted {
			// Use the hash from the change object
			stash.Files[change.Path] = change.Hash
		} else {
			stash.Files[change.Path] = "" // Empty string for deleted files
		}
	}

	// Save stash
	if err := saveStash(repoPath, stash); err != nil {
		return fmt.Errorf("failed to save stash: %w", err)
	}

	// Revert working directory to clean state
	if err := revertWorkingDirectory(repoPath, repo); err != nil {
		return fmt.Errorf("failed to revert working directory: %w", err)
	}

	fmt.Printf("Saved stash %s: %s\n", stash.ID, message)
	return nil
}

func stashList() error {
	repoPath, _ := os.Getwd()
	stashes, err := loadStashes(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load stashes: %w", err)
	}

	if len(stashes) == 0 {
		fmt.Println("No stashes found.")
		return nil
	}

	fmt.Println("Stashes:")
	for i, stash := range stashes {
		fmt.Printf("  stash@{%d}: %s on %s (%s)\n",
			i, stash.Message, stash.Branch, stash.Timestamp.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func stashApply(stashID string) error {
	repoPath, _ := os.Getwd()

	// Load stash
	stash, err := loadStashByID(repoPath, stashID)
	if err != nil {
		return fmt.Errorf("failed to load stash: %w", err)
	}

	// Apply stash files
	for file, blobRef := range stash.Files {
		filePath := filepath.Join(repoPath, file)

		if blobRef == "" {
			// File was deleted in stash
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

	fmt.Printf("Applied stash %s\n", stashID)
	return nil
}

func stashDrop(stashID string) error {
	repoPath, _ := os.Getwd()

	// Load stashes
	stashes, err := loadStashes(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load stashes: %w", err)
	}

	// Find stash by ID
	var targetStash *Stash
	for _, stash := range stashes {
		if stash.ID == stashID {
			targetStash = stash
			break
		}
	}

	if targetStash == nil {
		return fmt.Errorf("stash %s not found", stashID)
	}

	// Remove stash file
	stashPath := filepath.Join(repoPath, ".steria", "stashes", targetStash.ID)
	if err := os.Remove(stashPath); err != nil {
		return fmt.Errorf("failed to remove stash: %w", err)
	}

	fmt.Printf("Dropped stash %s\n", stashID)
	return nil
}

func stashPop(stashID string) error {
	// Apply stash
	if err := stashApply(stashID); err != nil {
		return err
	}

	// Drop stash
	return stashDrop(stashID)
}

func saveStash(repoPath string, stash *Stash) error {
	stashesDir := filepath.Join(repoPath, ".steria", "stashes")
	if err := os.MkdirAll(stashesDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(stash, "", "  ")
	if err != nil {
		return err
	}

	stashPath := filepath.Join(stashesDir, stash.ID)
	return os.WriteFile(stashPath, data, 0644)
}

func loadStashes(repoPath string) ([]*Stash, error) {
	stashesDir := filepath.Join(repoPath, ".steria", "stashes")
	entries, err := os.ReadDir(stashesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Stash{}, nil
		}
		return nil, err
	}

	var stashes []*Stash
	for _, entry := range entries {
		if !entry.IsDir() {
			stash, err := loadStash(repoPath, entry.Name())
			if err != nil {
				continue // Skip corrupted stashes
			}
			stashes = append(stashes, stash)
		}
	}

	return stashes, nil
}

func loadStash(repoPath, id string) (*Stash, error) {
	stashPath := filepath.Join(repoPath, ".steria", "stashes", id)
	data, err := os.ReadFile(stashPath)
	if err != nil {
		return nil, err
	}

	var stash Stash
	if err := json.Unmarshal(data, &stash); err != nil {
		return nil, err
	}

	return &stash, nil
}

func loadStashByID(repoPath, stashID string) (*Stash, error) {
	stashes, err := loadStashes(repoPath)
	if err != nil {
		return nil, err
	}

	for _, stash := range stashes {
		if stash.ID == stashID {
			return stash, nil
		}
	}

	return nil, fmt.Errorf("stash %s not found", stashID)
}

func revertWorkingDirectory(repoPath string, repo *storage.Repo) error {
	// Get current HEAD commit
	if repo.Head == "" {
		// No commits yet, just remove all files
		entries, err := os.ReadDir(repoPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.Name() != ".steria" {
				entryPath := filepath.Join(repoPath, entry.Name())
				if err := os.RemoveAll(entryPath); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Load HEAD commit
	commit, err := repo.LoadCommit(repo.Head)
	if err != nil {
		return err
	}

	// Restore files from HEAD
	for _, file := range commit.Files {
		filePath := filepath.Join(repoPath, file)
		blobRef := commit.FileBlobs[file]

		// Ensure directory exists
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Get blob data
		blobDir := filepath.Join(repoPath, ".steria", "objects", "blobs")
		store := &storage.LocalBlobStore{Dir: blobDir}
		data, err := storage.ReadFileBlobDecompressed(store, blobRef)
		if err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

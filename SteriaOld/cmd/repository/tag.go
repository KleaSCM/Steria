// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: tag.go
// Description: Implements tag management commands for Steria version control.

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

type Tag struct {
	Name      string    `json:"name"`
	Commit    string    `json:"commit"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags for releases",
		Long:  "Create, list, delete, and checkout tags for version management",
	}
	cmd.AddCommand(newTagCreateCmd())
	cmd.AddCommand(newTagListCmd())
	cmd.AddCommand(newTagDeleteCmd())
	cmd.AddCommand(newTagCheckoutCmd())
	return cmd
}

func newTagCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <name> [commit] [message]",
		Short: "Create a new tag",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			commit := ""
			message := ""

			if len(args) > 1 {
				commit = args[1]
			}
			if len(args) > 2 {
				message = args[2]
			}

			return createTag(name, commit, message)
		},
	}
}

func newTagListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listTags()
		},
	}
}

func newTagDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteTag(args[0])
		},
	}
}

func newTagCheckoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkout <name>",
		Short: "Checkout a tag (creates detached HEAD)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return checkoutTag(args[0])
		},
	}
}

func createTag(name, commit, message string) error {
	repoPath, _ := os.Getwd()

	// Load repository
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Use HEAD if no commit specified
	if commit == "" {
		commit = repo.Head
	}

	// Verify commit exists
	if _, err := repo.LoadCommit(commit); err != nil {
		return fmt.Errorf("commit %s not found: %w", commit, err)
	}

	// Create tag
	tag := &Tag{
		Name:      name,
		Commit:    commit,
		Message:   message,
		Author:    "KleaSCM",
		Timestamp: time.Now(),
	}

	// Save tag
	if err := saveTag(repoPath, tag); err != nil {
		return fmt.Errorf("failed to save tag: %w", err)
	}

	fmt.Printf("Created tag '%s' pointing to commit %s\n", name, commit[:8])
	return nil
}

func listTags() error {
	repoPath, _ := os.Getwd()
	tags, err := loadTags(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load tags: %w", err)
	}

	if len(tags) == 0 {
		fmt.Println("No tags found.")
		return nil
	}

	fmt.Println("Tags:")
	for _, tag := range tags {
		fmt.Printf("  %s -> %s (%s)\n", tag.Name, tag.Commit[:8], tag.Timestamp.Format("2006-01-02 15:04:05"))
		if tag.Message != "" {
			fmt.Printf("    %s\n", tag.Message)
		}
	}
	return nil
}

func deleteTag(name string) error {
	repoPath, _ := os.Getwd()
	tagPath := filepath.Join(repoPath, ".steria", "refs", "tags", name)

	if err := os.Remove(tagPath); err != nil {
		return fmt.Errorf("failed to delete tag '%s': %w", name, err)
	}

	fmt.Printf("Deleted tag '%s'\n", name)
	return nil
}

func checkoutTag(name string) error {
	repoPath, _ := os.Getwd()

	// Load tag
	tag, err := loadTag(repoPath, name)
	if err != nil {
		return fmt.Errorf("failed to load tag '%s': %w", name, err)
	}

	// Update HEAD to point to tagged commit
	headPath := filepath.Join(repoPath, ".steria", "HEAD")
	if err := os.WriteFile(headPath, []byte(tag.Commit), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	fmt.Printf("Checked out tag '%s' (commit %s)\n", name, tag.Commit[:8])
	fmt.Println("You are now in 'detached HEAD' state.")
	return nil
}

func saveTag(repoPath string, tag *Tag) error {
	tagsDir := filepath.Join(repoPath, ".steria", "refs", "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(tag, "", "  ")
	if err != nil {
		return err
	}

	tagPath := filepath.Join(tagsDir, tag.Name)
	return os.WriteFile(tagPath, data, 0644)
}

func loadTag(repoPath, name string) (*Tag, error) {
	tagPath := filepath.Join(repoPath, ".steria", "refs", "tags", name)
	data, err := os.ReadFile(tagPath)
	if err != nil {
		return nil, err
	}

	var tag Tag
	if err := json.Unmarshal(data, &tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

func loadTags(repoPath string) ([]*Tag, error) {
	tagsDir := filepath.Join(repoPath, ".steria", "refs", "tags")
	entries, err := os.ReadDir(tagsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Tag{}, nil
		}
		return nil, err
	}

	var tags []*Tag
	for _, entry := range entries {
		if !entry.IsDir() {
			tag, err := loadTag(repoPath, entry.Name())
			if err != nil {
				continue // Skip corrupted tags
			}
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

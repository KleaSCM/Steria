package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/core"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "done \"message\" - signer",
		Short: "Done! Commit, sign, and sync everything",
		Long: `The magical "done" command. When you're finished working:
- Automatically detects changes
- Creates a smart commit message
- Signs with your identity
- Syncs everything up
- Out of sight, out of mind!

Example: steria done "feat - added new feature" - KleaSCM`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := args[0]
			signer := args[1]
			return runDone(signer, message)
		},
	}

	return cmd
}

func runDone(signer, message string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Starting Steria done process...\n", cyan("🚀"))

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Initialize or load repo
	repo, err := core.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Check for changes
	changes, err := repo.GetChanges()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Printf("%s No changes detected. Everything is clean!\n", green("✨"))
		return nil
	}

	fmt.Printf("%s Found %d changed files\n", yellow("📝"), len(changes))

	// Generate commit message if not provided
	if message == "" {
		message = generateSmartMessage(changes)
	}

	// Create commit
	commit, err := repo.CreateCommit(message, signer)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	fmt.Printf("%s Created commit: %s\n", green("✅"), commit.Hash[:8])

	// Sync with remote if available
	if repo.HasRemote() {
		fmt.Printf("%s Syncing with remote...\n", cyan("🔄"))
		if err := repo.Sync(); err != nil {
			fmt.Printf("%s Warning: sync failed: %v\n", yellow("⚠️"), err)
		} else {
			fmt.Printf("%s Successfully synced!\n", green("🎉"))
		}
	}

	fmt.Printf("%s Done! Everything is committed and synced.\n", green("🎯"))
	fmt.Printf("%s You can now forget about it - out of sight, out of mind!\n", cyan("💫"))

	return nil
}

func generateSmartMessage(changes []core.FileChange) string {
	if len(changes) == 0 {
		return "Empty commit"
	}

	if len(changes) == 1 {
		change := changes[0]
		action := "Updated"
		if change.Type == core.ChangeTypeAdded {
			action = "Added"
		} else if change.Type == core.ChangeTypeDeleted {
			action = "Removed"
		}
		return fmt.Sprintf("%s %s", action, filepath.Base(change.Path))
	}

	// Count by type
	added, modified, deleted := 0, 0, 0
	for _, change := range changes {
		switch change.Type {
		case core.ChangeTypeAdded:
			added++
		case core.ChangeTypeModified:
			modified++
		case core.ChangeTypeDeleted:
			deleted++
		}
	}

	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", modified))
	}
	if deleted > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", deleted))
	}

	return fmt.Sprintf("Updated %s files", strings.Join(parts, ", "))
}

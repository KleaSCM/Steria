package workflow

import (
	"fmt"
	"os"

	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync with remote repository",
		Long:  "Sync changes with the remote repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync()
		},
	}

	return cmd
}

func runSync() error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	if !repo.HasRemote() {
		fmt.Printf("%s No remote configured. Use 'clone' to set up a remote.\n", yellow("⚠️"))
		return nil
	}

	fmt.Printf("%s Syncing with %s...\n", cyan("🔄"), repo.RemoteURL)

	// For now, this is a placeholder
	// In a full implementation, we'd handle actual syncing
	if err := repo.Sync(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Printf("%s Successfully synced!\n", green("✅"))
	return nil
}

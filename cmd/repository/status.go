package repository

import (
	"fmt"
	"os"

	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show repository status",
		Long:  "Show the current status of the repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus()
		},
	}

	return cmd
}

func runStatus() error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	fmt.Printf("%s Repository: %s\n", cyan("📁"), repo.Config.Name)
	fmt.Printf("%s Branch: %s\n", cyan("🌿"), green(repo.Branch))

	if repo.Head != "" {
		fmt.Printf("%s HEAD: %s\n", cyan("📍"), yellow(repo.Head[:8]))
	} else {
		fmt.Printf("%s HEAD: %s\n", cyan("📍"), red("no commits"))
	}

	if repo.RemoteURL != "" {
		fmt.Printf("%s Remote: %s\n", cyan("🌐"), repo.RemoteURL)
	} else {
		fmt.Printf("%s Remote: %s\n", cyan("🌐"), red("none"))
	}

	// Check for changes
	changes, err := repo.GetChanges()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Printf("%s Working directory is clean\n", green("✨"))
	} else {
		fmt.Printf("\n%s Changes:\n", yellow("📝"))
		for _, change := range changes {
			icon := "📄"
			color := green
			switch change.Type {
			case storage.ChangeTypeAdded:
				icon = "➕"
				color = green
			case storage.ChangeTypeModified:
				icon = "✏️"
				color = yellow
			case storage.ChangeTypeDeleted:
				icon = "🗑️"
				color = red
			}
			fmt.Printf("  %s %s\n", icon, color(change.Path))
		}
	}

	return nil
}

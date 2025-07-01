package repository

import (
	"fmt"
	"os"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show repository status",
		Long:  "Show the current status of the repository with optimized processing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus()
		},
	}

	return cmd
}

func runStatus() error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Checking status with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Create optimized repository
	optRepo := storage.NewOptimizedRepo(repo)

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

	// Check for changes with optimized method
	endOp := metrics.GlobalMetrics.StartOperation("get_changes")
	changes, err := optRepo.GetChangesOptimized()
	endOp()

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

	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

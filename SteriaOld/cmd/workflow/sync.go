package workflow

import (
	"fmt"
	"os"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync with remote repository",
		Long:  "Sync with remote repository using optimized processing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync()
		},
	}

	return cmd
}

func runSync() error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	_ = color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Starting optimized sync process...\n", cyan("🚀"))

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

	if !optRepo.HasRemote() {
		return fmt.Errorf("no remote configured for this repository")
	}

	fmt.Printf("%s Syncing with remote: %s\n", cyan("🔄"), optRepo.RemoteURL)

	// Perform sync with optimized method
	endOp := metrics.GlobalMetrics.StartOperation("sync")
	err = optRepo.Sync()
	endOp()

	if err != nil {
		fmt.Printf("%s Sync failed: %v\n", red("❌"), err)
		return err
	}

	fmt.Printf("%s Successfully synced with remote!\n", green("✅"))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

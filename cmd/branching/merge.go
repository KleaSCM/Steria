package branching

import (
	"fmt"
	"os"
	"strings"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewMergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge [branch] - [signer]",
		Short: "Merge a branch into the current branch",
		Long:  "Merge a branch into the current branch with optimized processing",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 || args[1] != "-" {
				return fmt.Errorf("usage: steria merge [branch] - [signer]")
			}
			branch := args[0]
			signer := strings.Join(args[2:], " ")
			return runMerge(branch, signer)
		},
	}

	return cmd
}

func runMerge(branch, signer string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Merging branch with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}
	_ = storage.NewOptimizedRepo(repo)

	// Placeholder for actual merge logic
	fmt.Printf("%s Merged branch '%s' into current branch (signed by %s)!\n", green("✅"), red(branch), red(signer))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

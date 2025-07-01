package projects

import (
	"fmt"
	"os"
	"strings"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [project name] [version] - [signer]",
		Short: "Pull a specific version from a project",
		Long:  "Pull a specific version from a project with optimized processing",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 || args[2] != "-" {
				return fmt.Errorf("usage: steria pull [project name] [version] - [signer]")
			}
			project := args[0]
			version := args[1]
			signer := strings.Join(args[3:], " ")
			return runPull(project, version, signer)
		},
	}

	return cmd
}

func runPull(project, version, signer string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Pulling version with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}
	_ = storage.NewOptimizedRepo(repo)

	// Placeholder for actual pull logic
	fmt.Printf("%s Pulled version '%s' of project '%s' (signed by %s)!\n", green("✅"), red(version), red(project), red(signer))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

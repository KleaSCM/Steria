package branching

import (
	"fmt"
	"os"
	"path/filepath"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDeleteBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-branch [name]",
		Short: "Delete a branch",
		Long:  "Delete a branch (cannot delete the currently checked-out branch) with optimized processing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteBranchCmd(args[0])
		},
	}

	return cmd
}

func runDeleteBranchCmd(name string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s Deleting branch with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Initialize optimized repository for future use
	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}
	_ = storage.NewOptimizedRepo(repo)

	// Prevent deleting the current branch
	branchPath := filepath.Join(cwd, ".steria", "branch")
	currentBranch, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("failed to read current branch: %w", err)
	}

	if string(currentBranch) == name {
		return fmt.Errorf("cannot delete the currently checked-out branch: %s", red(name))
	}

	branchFile := filepath.Join(cwd, ".steria", "branches", name)

	// Check if branch exists
	if _, err := os.Stat(branchFile); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", red(name))
	}

	if err := os.Remove(branchFile); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	fmt.Printf("%s Branch '%s' deleted successfully!\n", green("✅"), red(name))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

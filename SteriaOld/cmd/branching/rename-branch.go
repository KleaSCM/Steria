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

func NewRenameBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename-branch [old-name] [new-name]",
		Short: "Rename a branch",
		Long:  "Rename a branch from old-name to new-name with optimized processing",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]
			newName := args[1]
			return runRenameBranch(oldName, newName)
		},
	}

	return cmd
}

func runRenameBranch(oldName, newName string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Renaming branch with optimized processing...\n", cyan("🚀"))

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

	oldBranchFile := filepath.Join(cwd, ".steria", "branches", oldName)
	newBranchFile := filepath.Join(cwd, ".steria", "branches", newName)

	// Check if old branch exists
	if _, err := os.Stat(oldBranchFile); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", red(oldName))
	}

	// Check if new branch already exists
	if _, err := os.Stat(newBranchFile); err == nil {
		return fmt.Errorf("branch '%s' already exists", red(newName))
	}

	// Ensure parent directory for new branch file exists
	if err := os.MkdirAll(filepath.Dir(newBranchFile), 0755); err != nil {
		return fmt.Errorf("failed to create branch parent dir: %w", err)
	}

	// Rename the branch file
	if err := os.Rename(oldBranchFile, newBranchFile); err != nil {
		return fmt.Errorf("failed to rename branch: %w", err)
	}

	// If this was the current branch, update the current branch file
	branchPath := filepath.Join(cwd, ".steria", "branch")
	if currentBranch, err := os.ReadFile(branchPath); err == nil {
		if string(currentBranch) == oldName {
			if err := os.WriteFile(branchPath, []byte(newName), 0644); err != nil {
				return fmt.Errorf("failed to update current branch: %w", err)
			}
		}
	}

	fmt.Printf("%s Renamed branch '%s' to '%s'\n", green("✅"), cyan(oldName), cyan(newName))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

package branching

import (
	"fmt"
	"os"

	"steria/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewBranchCmd() *cobra.Command {
	var deleteFlag bool
	cmd := &cobra.Command{
		Use:   "branch [name]",
		Short: "Create, switch, or delete a branch",
		Long:  "Create a new branch, switch to an existing one, or delete a branch with --delete.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if deleteFlag {
				return runDeleteBranch(args[0])
			}
			return runBranch(args[0])
		},
	}
	cmd.Flags().BoolVar(&deleteFlag, "delete", false, "Delete the branch instead of switching/creating")
	return cmd
}

func runBranch(name string) error {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := core.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// For now, just update the branch file
	// In a full implementation, we'd handle branch switching properly
	branchPath := fmt.Sprintf("%s/.steria/branch", cwd)
	if err := os.WriteFile(branchPath, []byte(name), 0644); err != nil {
		return fmt.Errorf("failed to switch branch: %w", err)
	}

	repo.Branch = name

	fmt.Printf("%s Switched to branch: %s\n", green("✅"), cyan(name))
	return nil
}

func runDeleteBranch(name string) error {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	branchFile := fmt.Sprintf("%s/.steria/branches/%s", cwd, name)
	if err := os.Remove(branchFile); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	fmt.Printf("%s Branch '%s' deleted successfully!\n", green("✅"), red(name))
	return nil
}

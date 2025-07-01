package branching

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDeleteBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-branch [name]",
		Short: "Delete a branch",
		Long:  "Delete a branch (cannot delete the currently checked-out branch)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteBranchCmd(args[0])
		},
	}

	return cmd
}

func runDeleteBranchCmd(name string) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

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
	return nil
}

package branching

import (
	"fmt"
	"os"
	"path/filepath"

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

	branchesDir := filepath.Join(cwd, ".steria", "branches")
	branchFile := filepath.Join(branchesDir, name)

	// Create branches dir if not exists
	if err := os.MkdirAll(branchesDir, 0755); err != nil {
		return fmt.Errorf("failed to create branches dir: %w", err)
	}

	// If branch file doesn't exist, create it with current HEAD
	if _, err := os.Stat(branchFile); os.IsNotExist(err) {
		head := repo.Head
		if head == "" {
			head = ""
		}
		if err := os.WriteFile(branchFile, []byte(head), 0644); err != nil {
			return fmt.Errorf("failed to create branch: %w", err)
		}
	}

	// Switch branch: update .steria/branch and .steria/HEAD
	if err := os.WriteFile(filepath.Join(cwd, ".steria", "branch"), []byte(name), 0644); err != nil {
		return fmt.Errorf("failed to switch branch: %w", err)
	}
	branchHead, _ := os.ReadFile(branchFile)
	if err := os.WriteFile(filepath.Join(cwd, ".steria", "HEAD"), branchHead, 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	repo.Branch = name
	switchMsg := fmt.Sprintf("%s Switched to branch: %s\n", green("✅"), cyan(name))
	fmt.Print(switchMsg)
	return nil
}

func runDeleteBranch(name string) error {
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Prevent deleting the current branch
	branchPath := filepath.Join(cwd, ".steria", "branch")
	currentBranch, _ := os.ReadFile(branchPath)
	if string(currentBranch) == name {
		return fmt.Errorf("cannot delete the currently checked-out branch: %s", name)
	}

	branchFile := filepath.Join(cwd, ".steria", "branches", name)
	if err := os.Remove(branchFile); err != nil {
		return fmt.Errorf("failed to delete branch: %w", err)
	}

	fmt.Printf("%s Branch '%s' deleted successfully!\n", green("✅"), red(name))
	return nil
}

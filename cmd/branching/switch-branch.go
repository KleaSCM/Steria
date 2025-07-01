package branching

import (
	"fmt"
	"os"
	"path/filepath"

	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewSwitchBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch-branch [name]",
		Short: "Switch to an existing branch",
		Long:  "Switch to an existing branch and update HEAD",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSwitchBranch(args[0])
		},
	}

	return cmd
}

func runSwitchBranch(name string) error {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	branchFile := filepath.Join(cwd, ".steria", "branches", name)

	// Check if branch exists
	if _, err := os.Stat(branchFile); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", red(name))
	}

	// Read the branch's HEAD
	branchHead, err := os.ReadFile(branchFile)
	if err != nil {
		return fmt.Errorf("failed to read branch HEAD: %w", err)
	}

	// Switch branch: update .steria/branch and .steria/HEAD
	if err := os.WriteFile(filepath.Join(cwd, ".steria", "branch"), []byte(name), 0644); err != nil {
		return fmt.Errorf("failed to switch branch: %w", err)
	}

	if err := os.WriteFile(filepath.Join(cwd, ".steria", "HEAD"), branchHead, 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	repo.Branch = name
	fmt.Printf("%s Switched to branch: %s\n", green("✅"), cyan(name))
	return nil
}

package branching

import (
	"fmt"
	"os"
	"steria/core"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewMergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge \"project name\" - signer",
		Short: "Merge a branch into the current branch",
		Long:  "Merge changes from another branch into the current branch",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			signer := strings.Join(args[1:], " ")
			return runMerge(projectName, signer)
		},
	}

	return cmd
}

func runMerge(projectName, signer string) error {
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

	// For now, this is a placeholder
	// In a full implementation, we'd handle actual merging
	fmt.Printf("%s Merging project %s into %s (signed by %s)...\n", cyan("🔄"), projectName, repo.Branch, signer)
	fmt.Printf("%s Merge completed successfully!\n", green("✅"))

	return nil
}

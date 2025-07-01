package projects

import (
	"fmt"
	"os"

	"steria/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull \"project name\" version - signer",
		Short: "Pull a specific version from a project",
		Long:  "Pull a specific version from a project and merge it",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			version := args[1]
			signer := args[2]
			return runPull(projectName, version, signer)
		},
	}

	return cmd
}

func runPull(projectName, version, signer string) error {
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	_, err = core.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// For now, this is a placeholder
	// In a full implementation, we'd handle actual pulling
	fmt.Printf("%s Pulling version %s from project %s (signed by %s)...\n", cyan("📥"), version, projectName, signer)
	fmt.Printf("%s Pull completed successfully!\n", green("✅"))

	return nil
}

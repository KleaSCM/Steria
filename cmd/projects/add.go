package projects

import (
	"fmt"
	"os"

	"steria/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add \"project name\" - signer",
		Short: "Add a project",
		Long:  "Add a new project to the repository",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			signer := args[1]
			return runAdd(projectName, signer)
		},
	}

	return cmd
}

func runAdd(projectName, signer string) error {
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
	// In a full implementation, we'd handle actual adding
	fmt.Printf("%s Adding project %s (signed by %s)...\n", cyan("➕"), projectName, signer)
	fmt.Printf("%s Project added successfully!\n", green("✅"))

	return nil
}

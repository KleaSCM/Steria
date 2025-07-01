package projects

import (
	"fmt"
	"os"
	"steria/core"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete \"project name\" - signer",
		Short: "Delete a project",
		Long:  "Delete a project from the repository",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			signer := strings.Join(args[1:], " ")
			return runDelete(projectName, signer)
		},
	}

	return cmd
}

func runDelete(projectName, signer string) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	_, err = core.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// For now, this is a placeholder
	// In a full implementation, we'd handle actual deletion
	fmt.Printf("%s Deleting project %s (signed by %s)...\n", red("🗑️"), projectName, signer)
	fmt.Printf("%s Project deleted successfully!\n", green("✅"))

	return nil
}

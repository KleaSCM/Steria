package main

import (
	"fmt"
	"os"

	"steria/cmd/branching"
	"steria/cmd/projects"
	"steria/cmd/repository"
	"steria/cmd/workflow"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "steria",
		Short: "Steria - Get out of the way version control",
		Long: `Steria is a fast, version control system that just works.
When you're done working, just type "done" and sign it. That's it.`,
		Version: "0.1.0",
	}

	// Add repository commands
	rootCmd.AddCommand(repository.NewCloneCmd())
	rootCmd.AddCommand(repository.NewStatusCmd())

	// Add workflow commands
	rootCmd.AddCommand(workflow.NewDoneCmd())
	rootCmd.AddCommand(workflow.NewCommitCmd())
	rootCmd.AddCommand(workflow.NewSyncCmd())

	// Add branching commands
	rootCmd.AddCommand(branching.NewBranchCmd())
	rootCmd.AddCommand(branching.NewMergeCmd())

	// Add project commands
	rootCmd.AddCommand(projects.NewPullCmd())
	rootCmd.AddCommand(projects.NewDeleteCmd())
	rootCmd.AddCommand(projects.NewAddCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

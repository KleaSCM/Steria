// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: main.go
// Description: Main entry point for Steria CLI with all command registrations.

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
	rootCmd := &cobra.Command{
		Use:   "steria",
		Short: "Steria - A modern version control system",
		Long:  "Steria is a fast, efficient version control system with advanced features.",
	}

	// Add all command groups
	rootCmd.AddCommand(branching.NewAddBranchCmd())
	rootCmd.AddCommand(branching.NewBranchCmd())
	rootCmd.AddCommand(branching.NewDeleteBranchCmd())
	rootCmd.AddCommand(branching.NewMergeCmd())
	rootCmd.AddCommand(branching.NewRenameBranchCmd())
	rootCmd.AddCommand(branching.NewSwitchBranchCmd())
	rootCmd.AddCommand(branching.NewBranchGraphCmd())

	rootCmd.AddCommand(projects.NewAddCmd())
	rootCmd.AddCommand(projects.NewDeleteCmd())
	rootCmd.AddCommand(projects.NewPullCmd())

	rootCmd.AddCommand(repository.NewCloneCmd())
	rootCmd.AddCommand(repository.NewStatusCmd())
	rootCmd.AddCommand(repository.NewDiffCmd())
	rootCmd.AddCommand(repository.NewSearchCmd())
	rootCmd.AddCommand(repository.NewRestoreCmd())
	rootCmd.AddCommand(repository.NewIgnoreCmd())
	rootCmd.AddCommand(repository.NewRemoteCmd())
	rootCmd.AddCommand(repository.NewPushCmd())
	rootCmd.AddCommand(repository.NewPullCmd())
	rootCmd.AddCommand(repository.NewTagCmd())
	rootCmd.AddCommand(repository.NewCherryPickCmd())
	rootCmd.AddCommand(repository.NewStashCmd())
	rootCmd.AddCommand(repository.NewBlameCmd())
	rootCmd.AddCommand(repository.NewRebaseCmd())
	rootCmd.AddCommand(repository.NewConflictsCmd())
	rootCmd.AddCommand(repository.NewResolveCmd())

	rootCmd.AddCommand(workflow.NewCommitCmd())
	rootCmd.AddCommand(workflow.NewDoneCmd())
	rootCmd.AddCommand(workflow.NewSyncCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

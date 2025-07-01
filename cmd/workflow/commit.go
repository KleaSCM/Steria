package workflow

import (
	"fmt"
	"os"

	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit \"message\" - signer",
		Short: "Create a commit",
		Long:  "Create a commit with the current changes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := args[0]
			signer := args[1]
			return runCommit(message, signer)
		},
	}

	return cmd
}

func runCommit(message, signer string) error {
	green := color.New(color.FgGreen).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	commit, err := repo.CreateCommit(message, signer)
	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	fmt.Printf("%s Created commit: %s\n", green("✅"), commit.Hash[:8])
	return nil
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: conflicts.go
// Description: Implements the 'steria conflicts' CLI command to list unresolved merge/rebase conflicts in the repository.

package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewConflictsCmd returns the Cobra command for 'steria conflicts'
func NewConflictsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conflicts",
		Short: "List unresolved merge/rebase conflicts in the repository",
		Long:  "Shows all files with unresolved merge or rebase conflicts, including line numbers and details.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConflicts()
		},
	}
	return cmd
}

// runConflicts lists all unresolved conflicts in the current repository
func runConflicts() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repoRoot := findRepoRoot(cwd)
	if repoRoot == "" {
		return fmt.Errorf("not inside a Steria repository")
	}

	conflicts, err := storage.ListUnresolvedConflicts(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load conflicts: %w", err)
	}

	if len(conflicts) == 0 {
		color.New(color.FgGreen).Printf("\nNo unresolved conflicts! Your repository is clean.\n\n")
		return nil
	}

	magenta := color.New(color.FgMagenta).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Printf("\n%s\n", bold("Unresolved Conflicts:"))
	for _, c := range conflicts {
		fmt.Printf("%s %s\n", red("✖"), magenta(c.File))
		if c.Type == "line" && len(c.Lines) > 0 {
			fmt.Printf("    Lines: %v\n", c.Lines)
		}
		if c.Details != "" {
			fmt.Printf("    Details: %s\n", c.Details)
		}
		fmt.Println()
	}
	return nil
}

// findRepoRoot walks up from cwd to find the Steria repo root
func findRepoRoot(start string) string {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, ".steria")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

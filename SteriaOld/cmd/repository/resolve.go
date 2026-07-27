// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: resolve.go
// Description: Implements the 'steria resolve <file>' CLI command for resolving merge/rebase conflicts interactively.

package repository

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"steria/internal/storage"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewResolveCmd returns the Cobra command for 'steria resolve <file>'
func NewResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve <file>",
		Short: "Resolve a conflicted file interactively",
		Long:  "Opens the conflicted file in your editor. After editing, you can mark it as resolved.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResolve(args[0])
		},
	}
	return cmd
}

// runResolve opens the conflicted file in $EDITOR and marks as resolved if user confirms
func runResolve(file string) error {
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

	found := false
	for _, c := range conflicts {
		if c.File == file {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("file '%s' is not marked as conflicted or is already resolved", file)
	}

	filePath := filepath.Join(repoRoot, file)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file '%s' does not exist in the repository", file)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}
	color.New(color.FgCyan).Printf("\nOpening %s in %s...\n", file, editor)
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	color.New(color.FgYellow).Printf("\nDid you resolve all conflicts in %s? (y/n): ", file)
	reader := bufio.NewReader(os.Stdin)
	resp, _ := reader.ReadString('\n')
	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp != "y" && resp != "yes" {
		color.New(color.FgRed).Println("Aborted. File not marked as resolved.")
		return nil
	}

	user := os.Getenv("STERIA_USER")
	if user == "" {
		user = os.Getenv("USER")
	}
	if user == "" {
		user = "unknown"
	}

	if err := storage.ResolveConflict(repoRoot, file, user); err != nil {
		return fmt.Errorf("failed to mark conflict as resolved: %w", err)
	}
	color.New(color.FgGreen).Printf("\nConflict in %s marked as resolved!\n\n", file)
	return nil
}

package branching

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"steria/internal/storage"

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

func runSwitchBranch(branch string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	repoRoot := findRepoRoot(cwd)
	if repoRoot == "" {
		return fmt.Errorf("not inside a Steria repository")
	}

	branchesDir := filepath.Join(repoRoot, ".steria", "branches")
	branchPath := filepath.Join(branchesDir, branch)
	// Always allow switching to 'Stem' if it exists (default branch)
	if branch == "Stem" {
		if _, err := os.Stat(branchPath); err == nil {
			return doSwitchBranch(repoRoot, branch, branchPath)
		}
		return fmt.Errorf("default branch 'Stem' does not exist")
	}

	if _, err := os.Stat(branchPath); err != nil {
		return fmt.Errorf("branch '%s' does not exist", branch)
	}
	return doSwitchBranch(repoRoot, branch, branchPath)
}

func doSwitchBranch(repoRoot, branch, branchPath string) error {
	head, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("failed to read branch ref: %w", err)
	}
	// Update HEAD and branch
	headPath := filepath.Join(repoRoot, ".steria", "HEAD")
	branchFile := filepath.Join(repoRoot, ".steria", "branch")
	if err := os.WriteFile(headPath, head, 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}
	if err := os.WriteFile(branchFile, []byte(branch), 0644); err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	// Restore working directory to match HEAD commit of the target branch
	repo, err := storage.LoadOrInitRepo(repoRoot)
	if err != nil {
		return fmt.Errorf("failed to load repo for restore: %w", err)
	}
	commit, err := repo.LoadCommit(strings.TrimSpace(string(head)))
	if err != nil {
		return fmt.Errorf("failed to load HEAD commit for restore: %w", err)
	}
	blobDir := filepath.Join(repoRoot, ".steria", "objects", "blobs")
	store := &storage.LocalBlobStore{Dir: blobDir}
	for _, file := range commit.Files {
		blob := commit.FileBlobs[file]
		if blob == "" {
			continue
		}
		data, err := storage.ReadFileBlobDecompressed(store, blob)
		if err != nil {
			continue
		}
		os.MkdirAll(filepath.Dir(filepath.Join(repoRoot, file)), 0755)
		os.WriteFile(filepath.Join(repoRoot, file), data, 0644)
	}

	fmt.Printf("\nSwitched to branch '%s'\n\n", branch)
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

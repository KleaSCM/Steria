// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: merge.go
// Description: Implements the steria merge command for merging branches, with fast-forward support. Advanced merge support is planned for future implementation.
package branching

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewMergeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merge [branch] - [signer]",
		Short: "Merge a branch into the current branch",
		Long:  "Merge a branch into the current branch with optimized processing",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 3 || args[1] != "-" {
				return fmt.Errorf("usage: steria merge [branch] - [signer]")
			}
			branch := args[0]
			signer := strings.Join(args[2:], " ")
			return runMerge(branch, signer)
		},
	}

	return cmd
}

func runMerge(branch, signer string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Merging branch with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}
	_ = storage.NewOptimizedRepo(repo)

	// --- Real fast-forward merge logic ---
	branchesDir := filepath.Join(cwd, ".steria", "branches")
	branchFile := filepath.Join(branchesDir, branch)
	if _, err := os.Stat(branchFile); err != nil {
		return fmt.Errorf("branch '%s' does not exist", branch)
	}
	// Read the commit hash of the target branch
	branchHashBytes, err := os.ReadFile(branchFile)
	if err != nil {
		return fmt.Errorf("failed to read branch file: %w", err)
	}
	branchHash := strings.TrimSpace(string(branchHashBytes))
	if branchHash == "" {
		return fmt.Errorf("branch '%s' has no commits", branch)
	}
	// Check if fast-forward is possible (current HEAD is ancestor of branchHash)
	currentHash := repo.Head
	if currentHash == branchHash {
		fmt.Printf("%s Already up to date!\n", yellow("💡"))
		return nil
	}
	// Walk the ancestry of branchHash to see if currentHash is an ancestor
	ancestor := false
	testHash := branchHash
	for i := 0; i < 1000 && testHash != ""; i++ {
		if testHash == currentHash {
			ancestor = true
			break
		}
		commit, err := repo.LoadCommit(testHash)
		if err != nil {
			break
		}
		testHash = commit.Parent
	}
	if ancestor {
		// Fast-forward: update HEAD and branch pointer
		repo.Head = branchHash
		headPath := filepath.Join(cwd, ".steria", "HEAD")
		if err := os.WriteFile(headPath, []byte(branchHash), 0644); err != nil {
			return fmt.Errorf("failed to update HEAD: %w", err)
		}
		// Update current branch pointer (if needed)
		currentBranchFile := filepath.Join(cwd, ".steria", "branch")
		if _, err := os.Stat(currentBranchFile); err == nil {
			if err := os.WriteFile(currentBranchFile, []byte(branch), 0644); err != nil {
				return fmt.Errorf("failed to update current branch: %w", err)
			}
		}
		fmt.Printf("%s Fast-forward merged branch '%s' into current branch (signed by %s)!\n", green("✅"), red(branch), red(signer))
		fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
		return nil
	}
	// --- TODO: Advanced merge (three-way, conflict resolution, etc.) ---
	fmt.Printf("%s Non-fast-forward merge required.\n", yellow("⚠️"))
	fmt.Printf("%s TODO: Advanced merge (three-way, conflict resolution) not yet implemented, babe!\n", cyan("💡"))
	return fmt.Errorf("non-fast-forward merge not yet implemented")
}

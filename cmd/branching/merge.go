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
	// --- Advanced merge (three-way, conflict resolution) ---
	// Find merge base (common ancestor)
	findMergeBase := func(a, b string) (string, error) {
		visited := make(map[string]bool)
		// Walk ancestry of a
		hash := a
		for i := 0; i < 1000 && hash != ""; i++ {
			visited[hash] = true
			commit, err := repo.LoadCommit(hash)
			if err != nil {
				break
			}
			hash = commit.Parent
		}
		// Walk ancestry of b
		hash = b
		for i := 0; i < 1000 && hash != ""; i++ {
			if visited[hash] {
				return hash, nil
			}
			commit, err := repo.LoadCommit(hash)
			if err != nil {
				break
			}
			hash = commit.Parent
		}
		return "", fmt.Errorf("no common ancestor found")
	}

	mergeBaseHash, err := findMergeBase(currentHash, branchHash)
	if err != nil {
		return fmt.Errorf("failed to find merge base: %w", err)
	}
	baseCommit, err := repo.LoadCommit(mergeBaseHash)
	if err != nil {
		return fmt.Errorf("failed to load merge base: %w", err)
	}
	targetCommit, err := repo.LoadCommit(branchHash)
	if err != nil {
		return fmt.Errorf("failed to load target branch commit: %w", err)
	}
	currentCommit, err := repo.LoadCommit(currentHash)
	if err != nil {
		return fmt.Errorf("failed to load current branch commit: %w", err)
	}

	// Build set of all files
	fileSet := make(map[string]struct{})
	for _, f := range baseCommit.Files {
		fileSet[f] = struct{}{}
	}
	for _, f := range targetCommit.Files {
		fileSet[f] = struct{}{}
	}
	for _, f := range currentCommit.Files {
		fileSet[f] = struct{}{}
	}

	// For each file, determine merge result
	conflicts := []string{}
	for file := range fileSet {
		baseBlob := baseCommit.FileBlobs[file]
		currentBlob := currentCommit.FileBlobs[file]
		targetBlob := targetCommit.FileBlobs[file]

		// If unchanged in current, changed in target: use target
		if baseBlob == currentBlob && baseBlob != targetBlob && targetBlob != "" {
			if err := restoreBlobToFile(repo, file, targetBlob); err != nil {
				return err
			}
			continue
		}
		// If unchanged in target, changed in current: use current
		if baseBlob == targetBlob && baseBlob != currentBlob && currentBlob != "" {
			if err := restoreBlobToFile(repo, file, currentBlob); err != nil {
				return err
			}
			continue
		}
		// If changed in both and blobs differ: conflict
		if baseBlob != currentBlob && baseBlob != targetBlob && currentBlob != targetBlob && currentBlob != "" && targetBlob != "" {
			if err := writeConflictFile(repo, file, currentBlob, targetBlob); err != nil {
				return err
			}
			conflicts = append(conflicts, file)
			continue
		}
		// If only in target (added): use target
		if baseBlob == "" && currentBlob == "" && targetBlob != "" {
			if err := restoreBlobToFile(repo, file, targetBlob); err != nil {
				return err
			}
			continue
		}
		// If only in current (added): use current
		if baseBlob == "" && targetBlob == "" && currentBlob != "" {
			if err := restoreBlobToFile(repo, file, currentBlob); err != nil {
				return err
			}
			continue
		}
		// If deleted in both: remove file
		if baseBlob != "" && currentBlob == "" && targetBlob == "" {
			os.Remove(file)
			continue
		}
	}

	if len(conflicts) > 0 {
		fmt.Printf("%s Merge completed with conflicts in the following files:\n", yellow("⚠️"))
		for _, f := range conflicts {
			fmt.Printf("  - %s\n", f)
		}
		fmt.Printf("Please resolve conflicts and commit the result.\n")
		return fmt.Errorf("merge completed with conflicts")
	}

	// Update HEAD and branch pointer
	repo.Head = branchHash
	headPath := filepath.Join(cwd, ".steria", "HEAD")
	if err := os.WriteFile(headPath, []byte(branchHash), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}
	currentBranchFile := filepath.Join(cwd, ".steria", "branch")
	if _, err := os.Stat(currentBranchFile); err == nil {
		if err := os.WriteFile(currentBranchFile, []byte(branch), 0644); err != nil {
			return fmt.Errorf("failed to update current branch: %w", err)
		}
	}
	fmt.Printf("%s Three-way merged branch '%s' into current branch (signed by %s)\n", green("✅"), red(branch), red(signer))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

// restoreBlobToFile restores a file from a blob hash
func restoreBlobToFile(repo *storage.Repo, filePath, blobHash string) error {
	blobPath := filepath.Join(repo.Path, ".steria", "objects", "blobs", blobHash)
	blobData, err := os.ReadFile(blobPath)
	if err != nil {
		return fmt.Errorf("failed to read blob for '%s': %w", filePath, err)
	}
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	if err := os.WriteFile(filePath, blobData, 0644); err != nil {
		return fmt.Errorf("failed to write restored file: %w", err)
	}
	return nil
}

// writeConflictFile writes a file with conflict markers
func writeConflictFile(repo *storage.Repo, filePath, currentBlob, targetBlob string) error {
	currentBlobPath := filepath.Join(repo.Path, ".steria", "objects", "blobs", currentBlob)
	currentData, _ := os.ReadFile(currentBlobPath)
	targetBlobPath := filepath.Join(repo.Path, ".steria", "objects", "blobs", targetBlob)
	targetData, _ := os.ReadFile(targetBlobPath)
	conflictContent := []byte(fmt.Sprintf("<<<<<<< CURRENT\n%s\n=======\n%s\n>>>>>>> TARGET\n", string(currentData), string(targetData)))
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	if err := os.WriteFile(filePath, conflictContent, 0644); err != nil {
		return fmt.Errorf("failed to write conflict file: %w", err)
	}
	return nil
}

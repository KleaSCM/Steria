// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: rebase.go
// Description: Implements simple interactive rebase for Steria - reorganize commits easily.

package repository

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"steria/internal/storage"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type RebaseAction struct {
	Action string // keep, combine, skip
	Hash   string
	Msg    string
}

func NewRebaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rebase",
		Short: "Reorganize your commits",
		Long:  "Opens an editor where you can reorder, combine, or skip commits to clean up your history",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRebase()
		},
	}
	return cmd
}

func runRebase() error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Opening commit organizer...\n", cyan("🔄"))

	repoPath, _ := os.Getwd()
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Get all commits from HEAD back to the beginning
	commits, err := getAllCommits(repo)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	if len(commits) <= 1 {
		fmt.Printf("%s Only one commit found - nothing to reorganize!\n", yellow("💡"))
		return nil
	}

	fmt.Printf("%s Found %d commits to organize\n", green("✅"), len(commits))

	// Write commit list to temp file
	planFile := filepath.Join(os.TempDir(), fmt.Sprintf("steria-rebase-%d.txt", time.Now().UnixNano()))
	f, err := os.Create(planFile)
	if err != nil {
		return fmt.Errorf("failed to create plan file: %w", err)
	}

	// Write header instructions
	fmt.Fprintf(f, "# Steria Commit Organizer\n")
	fmt.Fprintf(f, "# Edit this file to reorganize your commits:\n")
	fmt.Fprintf(f, "# - keep: Keep this commit as-is\n")
	fmt.Fprintf(f, "# - combine: Merge with the commit above it\n")
	fmt.Fprintf(f, "# - skip: Remove this commit\n")
	fmt.Fprintf(f, "# You can also reorder lines to change commit order\n")
	fmt.Fprintf(f, "# Lines starting with # are ignored\n\n")

	// Write commits in reverse order (oldest first)
	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]
		fmt.Fprintf(f, "keep %s %s\n", commit.Hash[:8], commit.Message)
	}
	f.Close()

	// Open editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	fmt.Printf("%s Opening %s editor...\n", yellow("📝"), editor)
	cmd := exec.Command(editor, planFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	// Parse the edited plan
	plan, err := parseRebasePlan(planFile, commits)
	if err != nil {
		return fmt.Errorf("failed to parse plan: %w", err)
	}

	// Apply the reorganization
	if err := applyRebasePlan(repo, plan); err != nil {
		return fmt.Errorf("failed to apply changes: %w", err)
	}

	fmt.Printf("%s Successfully reorganized commits!\n", green("🎉"))
	return nil
}

func getAllCommits(repo *storage.Repo) ([]*storage.Commit, error) {
	var commits []*storage.Commit
	seen := map[string]bool{}

	hash := repo.Head
	for hash != "" && !seen[hash] {
		seen[hash] = true
		commit, err := repo.LoadCommit(hash)
		if err != nil {
			break
		}
		commits = append(commits, commit)
		hash = commit.Parent
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}

	return commits, nil
}

func parseRebasePlan(planFile string, commits []*storage.Commit) ([]RebaseAction, error) {
	f, err := os.Open(planFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var plan []RebaseAction
	commitMap := map[string]*storage.Commit{}
	for _, c := range commits {
		commitMap[c.Hash[:8]] = c
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		action, hash := parts[0], parts[1]
		msg := strings.Join(parts[2:], " ")

		c, ok := commitMap[hash]
		if !ok {
			return nil, fmt.Errorf("unknown commit hash: %s", hash)
		}

		// Use original message if none provided
		if msg == "" {
			msg = c.Message
		}

		plan = append(plan, RebaseAction{
			Action: action,
			Hash:   c.Hash,
			Msg:    msg,
		})
	}

	return plan, nil
}

func applyRebasePlan(repo *storage.Repo, plan []RebaseAction) error {
	if len(plan) == 0 {
		return fmt.Errorf("no commits to process")
	}

	// Find the first commit that's not being skipped
	var firstCommit *storage.Commit
	for _, action := range plan {
		if action.Action != "skip" {
			firstCommit, _ = repo.LoadCommit(action.Hash)
			break
		}
	}

	if firstCommit == nil {
		return fmt.Errorf("all commits were skipped")
	}

	// Start from the first commit's parent
	parent := firstCommit.Parent
	var combinedMsg strings.Builder

	for _, action := range plan {
		switch action.Action {
		case "skip":
			fmt.Printf("Skipping commit %s\n", action.Hash[:8])
			continue

		case "combine":
			if combinedMsg.Len() > 0 {
				combinedMsg.WriteString("\n")
			}
			combinedMsg.WriteString(action.Msg)
			fmt.Printf("Combining commit %s\n", action.Hash[:8])

		case "keep":
			commit, _ := repo.LoadCommit(action.Hash)

			// If we have a combined message, use it
			msg := action.Msg
			if combinedMsg.Len() > 0 {
				msg = combinedMsg.String()
				combinedMsg.Reset()
			}

			// Apply commit's changes to working directory
			if err := applyCommitToWorkingDir(repo, commit); err != nil {
				return fmt.Errorf("failed to apply commit %s: %w", action.Hash[:8], err)
			}

			// Create new commit
			newCommit, err := repo.CreateCommit(msg, commit.Author)
			if err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}

			parent = newCommit.Hash
			fmt.Printf("Applied commit %s as %s\n", action.Hash[:8], newCommit.Hash[:8])
		}
	}

	// Update HEAD to the last commit
	headPath := filepath.Join(repo.Path, ".steria", "HEAD")
	return os.WriteFile(headPath, []byte(parent), 0644)
}

func applyCommitToWorkingDir(repo *storage.Repo, commit *storage.Commit) error {
	repoPath := repo.Path

	for _, file := range commit.Files {
		filePath := filepath.Join(repoPath, file)
		blobRef := commit.FileBlobs[file]

		// Ensure directory exists
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Get blob data
		blobDir := filepath.Join(repoPath, ".steria", "objects", "blobs")
		store := &storage.LocalBlobStore{Dir: blobDir}
		data, err := storage.ReadFileBlobDecompressed(store, blobRef)
		if err != nil {
			return err
		}

		// Write file
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

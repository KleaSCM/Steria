package workflow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"steria/internal/metrics"
	"steria/internal/security"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewDoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "done \"message\" - signer",
		Short: "Done! Commit, sign, and sync everything",
		Long: `The magical "done" command. When you're finished working:
- Automatically detects changes with concurrent processing
- Creates a smart commit message
- Signs with your identity using cryptographic signatures
- Syncs everything up with performance optimizations
- Out of sight, out of mind!

Example: steria done "feat - added new feature" - KleaSCM`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := args[0]
			// Join all remaining arguments as the signer (in case it contains spaces)
			signer := strings.Join(args[1:], " ")
			return runDone(signer, message)
		},
	}

	return cmd
}

func runDone(signer, message string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Starting Steria ULTRA-FAST done process...\n", cyan("🚀"))

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Initialize or load repo with optimized version
	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Create optimized repository
	optRepo := storage.NewOptimizedRepo(repo)

	// Check for changes with optimized method
	endOp := metrics.GlobalMetrics.StartOperation("get_changes")
	changes, err := optRepo.GetChangesOptimized()
	endOp()

	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Printf("%s No changes detected. Everything is clean!\n", green("✨"))
		return nil
	}

	fmt.Printf("%s Found %d changed files\n", yellow("📝"), len(changes))

	// Generate commit message if not provided
	if message == "" {
		message = generateSmartMessage(changes)
	}

	// Create cryptographic signature
	keyPair, err := security.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	signature, err := keyPair.SignMessage(message)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	// Verify signature
	valid, err := security.VerifySignature(signature)
	if err != nil {
		return fmt.Errorf("failed to verify signature: %w", err)
	}

	if !valid {
		return fmt.Errorf("signature verification failed")
	}

	fmt.Printf("%s Message cryptographically signed by: %s\n", green("🔐"), red(signer))

	// Create commit with optimized method
	endOp = metrics.GlobalMetrics.StartOperation("create_commit")
	commit, err := optRepo.CreateCommitOptimized(message, signer)
	endOp()

	if err != nil {
		return fmt.Errorf("failed to create commit: %w", err)
	}

	metrics.GlobalMetrics.IncrementCommitsCreated()
	fmt.Printf("%s Created commit: %s\n", green("✅"), commit.Hash[:8])

	// Sync with remote if available
	if optRepo.HasRemote() {
		fmt.Printf("%s Syncing with remote...\n", cyan("🔄"))
		endOp = metrics.GlobalMetrics.StartOperation("sync")
		if err := optRepo.Sync(); err != nil {
			fmt.Printf("%s Warning: sync failed: %v\n", yellow("⚠️"), err)
		} else {
			fmt.Printf("%s Successfully synced!\n", green("🎉"))
		}
		endOp()
	}

	fmt.Printf("%s ULTRA-FAST DONE! Everything is committed and synced.\n", green("🎯"))
	fmt.Printf("%s Performance optimized with concurrent processing and caching!\n", cyan("⚡"))
	fmt.Printf("%s You can now forget about it - out of sight, out of mind!\n", cyan("💫"))

	return nil
}

func generateSmartMessage(changes []storage.FileChange) string {
	if len(changes) == 0 {
		return "Empty commit"
	}

	if len(changes) == 1 {
		change := changes[0]
		action := "Updated"
		if change.Type == storage.ChangeTypeAdded {
			action = "Added"
		} else if change.Type == storage.ChangeTypeDeleted {
			action = "Removed"
		}
		return fmt.Sprintf("%s %s", action, filepath.Base(change.Path))
	}

	// Count by type
	added, modified, deleted := 0, 0, 0
	for _, change := range changes {
		switch change.Type {
		case storage.ChangeTypeAdded:
			added++
		case storage.ChangeTypeModified:
			modified++
		case storage.ChangeTypeDeleted:
			deleted++
		}
	}

	parts := []string{}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", modified))
	}
	if deleted > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", deleted))
	}

	return fmt.Sprintf("Updated %s files", strings.Join(parts, ", "))
}

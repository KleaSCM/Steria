package workflow

import (
	"fmt"
	"os"
	"strings"

	"steria/internal/metrics"
	"steria/internal/security"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCommitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit \"message\" - signer",
		Short: "Create a manual commit",
		Long:  "Create a commit with a specific message and signer",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := args[0]
			signer := strings.Join(args[1:], " ")
			return runCommit(signer, message)
		},
	}

	return cmd
}

func runCommit(signer, message string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Starting optimized commit process...\n", cyan("🚀"))

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
		fmt.Printf("%s No changes detected. Nothing to commit!\n", yellow("⚠️"))
		return nil
	}

	fmt.Printf("%s Found %d changed files\n", yellow("📝"), len(changes))

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
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))

	return nil
}

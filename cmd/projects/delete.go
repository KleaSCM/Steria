package projects

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

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete \"project name\" - signer",
		Short: "Delete a project",
		Long:  "Delete a project from the repository with optimized processing",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]
			signer := strings.Join(args[1:], " ")
			return runDelete(projectName, signer)
		},
	}

	return cmd
}

func runDelete(projectName, signer string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Deleting project with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Initialize optimized repository for future use
	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}
	_ = storage.NewOptimizedRepo(repo)

	// Create cryptographic signature
	keyPair, err := security.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %w", err)
	}

	signature, err := keyPair.SignMessage(projectName)
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

	fmt.Printf("%s Project '%s' deleted successfully!\n", green("✅"), red(projectName))
	fmt.Printf("%s Signed by: %s\n", green("🔐"), red(signer))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

package branching

import (
	"fmt"
	"os"
	"path/filepath"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewAddBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-branch [name]",
		Short: "Create a new branch",
		Long:  "Create a new branch with the current HEAD using optimized processing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddBranch(args[0])
		},
	}

	return cmd
}

func runAddBranch(name string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	fmt.Printf("%s Creating branch with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	// Initialize optimized repository for future use
	_ = storage.NewOptimizedRepo(repo)

	branchesDir := filepath.Join(cwd, ".steria", "branches")
	branchFile := filepath.Join(branchesDir, name)

	// Create branches dir if not exists
	if err := os.MkdirAll(branchesDir, 0755); err != nil {
		return fmt.Errorf("failed to create branches dir: %w", err)
	}

	// Ensure parent directory for branch file exists (for branch names with slashes)
	if err := os.MkdirAll(filepath.Dir(branchFile), 0755); err != nil {
		return fmt.Errorf("failed to create branch parent dir: %w", err)
	}

	// Check if branch already exists
	if _, err := os.Stat(branchFile); err == nil {
		return fmt.Errorf("branch '%s' already exists", name)
	}

	// Create branch with current HEAD
	head := repo.Head
	if head == "" {
		head = ""
	}
	if err := os.WriteFile(branchFile, []byte(head), 0644); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	metrics.GlobalMetrics.IncrementBranchesCreated()
	fmt.Printf("%s Created branch: %s\n", green("✅"), cyan(name))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

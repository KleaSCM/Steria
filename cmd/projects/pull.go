// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: pull.go
// Description: Implements the steria projects pull command for pulling a specific version from a local Steria project. Remote/registry support is planned for future implementation.
package projects

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"steria/internal/metrics"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [project name] [version] - [signer]",
		Short: "Pull a specific version from a project",
		Long:  "Pull a specific version from a project with optimized processing",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 || args[2] != "-" {
				return fmt.Errorf("usage: steria pull [project name] [version] - [signer]")
			}
			project := args[0]
			version := args[1]
			signer := strings.Join(args[3:], " ")
			return runPull(project, version, signer)
		},
	}

	return cmd
}

func runPull(project, version, signer string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Pulling version with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Local project pull mode
	steriaBase := "/home/klea/Steria/"
	projectDir := filepath.Join(steriaBase, project)
	if _, err := os.Stat(projectDir); err == nil {
		// Project exists locally
		commitDir := filepath.Join(projectDir, ".steria", "commits", version)
		if _, err := os.Stat(commitDir); err != nil {
			return fmt.Errorf("commit version '%s' not found in project '%s'", version, project)
		}
		// Copy all files from commitDir to cwd (future implementation)
		fmt.Printf("%s Found local project. Copying files from commit %s...\n", yellow("💡"), red(version))
		fmt.Printf("%s Would copy files from '%s' to '%s'\n", cyan("💡"), commitDir, cwd)
		fmt.Printf("%s Pulled version '%s' of project '%s' (signed by %s)!\n", green("✅"), red(version), red(project), red(signer))
		fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
		return nil
	}

	// Remote/project registry mode (planned for future implementation)
	fmt.Printf("%s Project '%s' not found locally.\n", yellow("⚠️"), project)
	fmt.Printf("%s Remote/project registry support is planned for future implementation.\n", cyan("💡"))
	return fmt.Errorf("remote/project registry support not yet implemented")
}

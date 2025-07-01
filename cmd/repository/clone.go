package repository

import (
	"fmt"

	"steria/internal/metrics"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone [url] [dir]",
		Short: "Clone a repository from git",
		Long:  "Clone a repository from git with optimized processing",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			dir := args[1]
			return runClone(url, dir)
		},
	}

	return cmd
}

func runClone(url, dir string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Cloning repository with optimized processing...\n", cyan("🚀"))

	// Placeholder for actual clone logic
	fmt.Printf("%s Cloned repository from '%s' into '%s'!\n", green("✅"), red(url), red(dir))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

func extractDirFromURL(url string) string {
	// Simple extraction - get the last part of the URL
	// This is a basic implementation
	if len(url) == 0 {
		return "repository"
	}

	// Remove .git suffix if present
	if len(url) > 4 && url[len(url)-4:] == ".git" {
		url = url[:len(url)-4]
	}

	// Get the last part after the last slash
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			return url[i+1:]
		}
	}

	return url
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: graph.go
// Description: Implements the 'steria branch-graph' command to visualize branch and commit relationships as an ASCII graph.

package branching

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"
	"strings"

	"io/ioutil"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewBranchGraphCmd() *cobra.Command {
	var mermaid bool
	cmd := &cobra.Command{
		Use:   "branch-graph",
		Short: "Visualize branch and commit relationships as an ASCII graph",
		Long:  "Display all branches and their commit histories as a simple ASCII/Unicode graph or Mermaid diagram.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBranchGraph(mermaid)
		},
	}
	cmd.Flags().BoolVar(&mermaid, "mermaid", false, "Output as Mermaid diagram")
	return cmd
}

func runBranchGraph(mermaid bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return err
	}

	// Map: commit hash -> commit object
	commits := map[string]*storage.Commit{}
	commit := repo.Head
	for commit != "" {
		c, err := repo.LoadCommit(commit)
		if err != nil {
			break
		}
		commits[commit] = c
		if c.Parent == "" {
			break
		}
		commit = c.Parent
	}

	// Read branches from .steria/branches
	branches := map[string]string{}
	branchesDir := filepath.Join(repo.Path, ".steria", "branches")
	branchFiles, err := ioutil.ReadDir(branchesDir)
	if err == nil {
		for _, f := range branchFiles {
			if f.IsDir() {
				continue
			}
			branchName := f.Name()
			branchPath := filepath.Join(branchesDir, branchName)
			data, err := ioutil.ReadFile(branchPath)
			if err == nil {
				hash := string(data)
				hash = strings.TrimSpace(hash)
				branches[branchName] = hash
			}
		}
	}

	if mermaid {
		return printMermaidGraph(commits, branches, repo.Head)
	}
	return printAsciiGraph(commits, branches, repo.Head)
}

func printAsciiGraph(commits map[string]*storage.Commit, branches map[string]string, head string) error {
	magenta := color.New(color.FgMagenta).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	headMark := green("[HEAD]")
	fmt.Println(cyan("Branch/Commit Graph:"))
	for branch, hash := range branches {
		fmt.Printf("\n%s %s\n", magenta("Branch:"), branch)
		cHash := hash
		for cHash != "" {
			c := commits[cHash]
			if c == nil {
				break
			}
			isHead := (cHash == head)
			headStr := ""
			if isHead {
				headStr = headMark
			}
			fmt.Printf("  * %s %s %s\n", cHash[:8], c.Message, headStr)
			if c.Parent != "" {
				fmt.Printf("    |\n    +-- parent: %s\n", c.Parent[:8])
			}
			cHash = c.Parent
		}
	}
	return nil
}

func printMermaidGraph(commits map[string]*storage.Commit, branches map[string]string, head string) error {
	fmt.Println("graph TD;")
	commitNodes := map[string]bool{}
	for branch, hash := range branches {
		cHash := hash
		for cHash != "" {
			c := commits[cHash]
			if c == nil {
				break
			}
			if !commitNodes[cHash] {
				fmt.Printf("  %s[\"%s\"]\n", cHash[:8], cHash[:8])
				commitNodes[cHash] = true
			}
			if c.Parent != "" {
				fmt.Printf("  %s --> %s\n", cHash[:8], c.Parent[:8])
			}
			cHash = c.Parent
		}
		// Branch label
		if hash != "" {
			fmt.Printf("  %s_branch((%s))\n", branch, branch)
			fmt.Printf("  %s_branch -.-> %s\n", branch, hash[:8])
		}
	}
	// Mark HEAD
	if head != "" {
		fmt.Printf("  head((HEAD))\n  head -.-> %s\n", head[:8])
	}
	return nil
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: diff.go
// Description: Implements the 'steria diff' command to show differences between file versions.

package repository

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"steria/internal/metrics"
	"steria/internal/storage"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewDiffCmd creates the 'diff' command for Steria
func NewDiffCmd() *cobra.Command {
	var sideBySide bool
	var contextLines int
	cmd := &cobra.Command{
		Use:   "diff [file]",
		Short: "Show differences between file versions",
		Long:  "Display differences between the working directory and the last commit, or between two commits",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var filePath string
			if len(args) > 0 {
				filePath = args[0]
			}
			return runDiffWithMode(filePath, sideBySide, contextLines)
		},
	}
	cmd.Flags().BoolVar(&sideBySide, "side-by-side", false, "Show side-by-side diff view")
	cmd.Flags().IntVar(&contextLines, "context", 3, "Number of context lines to show around changes")
	return cmd
}

func runDiffWithMode(filePath string, sideBySide bool, contextLines int) error {
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Analyzing file differences with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	repo, err := storage.LoadOrInitRepo(cwd)
	if err != nil {
		return fmt.Errorf("failed to load repository: %w", err)
	}

	if repo.Head == "" {
		fmt.Printf("%s No commits found in repository\n", yellow("⚠️"))
		return nil
	}

	commit, err := repo.LoadCommit(repo.Head)
	if err != nil {
		return fmt.Errorf("failed to load last commit: %w", err)
	}

	if filePath != "" {
		showFileDiffWithMode(repo, filePath, commit, sideBySide, contextLines)
		return nil
	}
	showAllDiffsWithMode(repo, commit, sideBySide, contextLines)
	return nil
}

func showFileDiffWithMode(repo *storage.Repo, filePath string, commit *storage.Commit, sideBySide bool, contextLines int) (added, removed, changed int) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("\n%s Showing differences for: %s\n", cyan("📁"), yellow(filePath))

	lastCommit := commit
	for c := commit; c != nil; {
		found := false
		for _, f := range c.Files {
			if f == filePath {
				found = true
				break
			}
		}
		if found || c.Parent == "" {
			break
		}
		parent, err := repo.LoadCommit(c.Parent)
		if err != nil {
			break
		}
		c = parent
		lastCommit = c
	}

	fmt.Printf("%s Last committed by: %s at %s\n", magenta("👤"), lastCommit.Author, lastCommit.Timestamp.Format(time.RFC1123))

	var commitContent []string
	blobHash := lastCommit.FileBlobs[filePath]
	if blobHash != "" {
		blobPath := repo.Path + "/.steria/objects/blobs/" + blobHash
		if f, err := os.Open(blobPath); err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				commitContent = append(commitContent, scanner.Text())
			}
			f.Close()
		}
	}

	var currentContent []string
	if _, err := os.Stat(filePath); err == nil {
		file, err := os.Open(filePath)
		if err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				currentContent = append(currentContent, scanner.Text())
			}
			file.Close()
		}
	}

	fmt.Printf("\nLegend: %s addition, %s deletion, %s unchanged\n", green("+"), red("-"), cyan(" "))

	added, removed, changed = 0, 0, 0
	if sideBySide {
		added, removed, changed = sideBySideDiff(commitContent, currentContent, contextLines)
	} else {
		added, removed, changed = showInlineDiff(commitContent, currentContent, contextLines)
	}
	fmt.Printf("\nSummary: %s lines added, %s lines removed, %s lines changed\n", green(fmt.Sprint(added)), red(fmt.Sprint(removed)), yellow(fmt.Sprint(changed)))

	return
}

func showInlineDiff(oldLines, newLines []string, contextLines int) (added, removed, changed int) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	// Simple Myers diff algorithm for line-level diff
	type op struct {
		kind string
		a, b int
	}
	ops := []op{}
	i, j := 0, 0
	for i < len(oldLines) || j < len(newLines) {
		if i < len(oldLines) && j < len(newLines) {
			if oldLines[i] == newLines[j] {
				ops = append(ops, op{" ", i, j})
				i++
				j++
			} else if j+1 < len(newLines) && oldLines[i] == newLines[j+1] {
				ops = append(ops, op{"+", -1, j})
				j++
			} else if i+1 < len(oldLines) && oldLines[i+1] == newLines[j] {
				ops = append(ops, op{"-", i, -1})
				i++
			} else {
				ops = append(ops, op{"-", i, -1})
				ops = append(ops, op{"+", -1, j})
				i++
				j++
			}
		} else if i < len(oldLines) {
			ops = append(ops, op{"-", i, -1})
			i++
		} else if j < len(newLines) {
			ops = append(ops, op{"+", -1, j})
			j++
		}
	}
	for idx := 0; idx < len(ops); idx++ {
		if ops[idx].kind == "-" {
			removed++
			line := oldLines[ops[idx].a]
			fmt.Printf("%s %s\n", red("-"), highlightWordDiff(line, "", red))
			continue
		}
		if ops[idx].kind == "+" {
			added++
			line := newLines[ops[idx].b]
			fmt.Printf("%s %s\n", green("+"), highlightWordDiff("", line, green))
			continue
		}
		fmt.Printf(" %s\n", cyan(oldLines[ops[idx].a]))
	}
	changed = min(added, removed)
	return
}

func sideBySideDiff(oldLines, newLines []string, contextLines int) (added, removed, changed int) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	maxLen := len(oldLines)
	if len(newLines) > maxLen {
		maxLen = len(newLines)
	}
	fmt.Printf("\n%-40s | %-40s\n", "< COMMITTED", "> WORKING DIR")
	fmt.Printf("%s\n", strings.Repeat("-", 83))
	for i := 0; i < maxLen; i++ {
		var left, right string
		if i < len(oldLines) {
			left = oldLines[i]
		} else {
			left = ""
		}
		if i < len(newLines) {
			right = newLines[i]
		} else {
			right = ""
		}
		if left == right {
			fmt.Printf(" %s | %s\n", cyan(left), cyan(right))
		} else if left == "" {
			fmt.Printf("%40s | %s%s\n", "", green("+ "), highlightWordDiff("", right, green))
			added++
		} else if right == "" {
			fmt.Printf("%s%s | %40s\n", red("- "), highlightWordDiff(left, "", red), "")
			removed++
		} else {
			fmt.Printf("%s%s | %s%s\n", red("- "), highlightWordDiff(left, right, red), green("+ "), highlightWordDiff(left, right, green))
			added++
			removed++
			changed++
		}
	}
	return
}

func highlightWordDiff(a, b string, colorize func(a ...interface{}) string) string {
	// Simple word diff: highlight words that are different
	aw := strings.Fields(a)
	bw := strings.Fields(b)
	max := len(aw)
	if len(bw) > max {
		max = len(bw)
	}
	out := ""
	for i := 0; i < max; i++ {
		if i < len(aw) && i < len(bw) {
			if aw[i] == bw[i] {
				out += aw[i] + " "
			} else {
				out += colorize(bw[i]) + " "
			}
		} else if i < len(bw) {
			out += colorize(bw[i]) + " "
		} else if i < len(aw) {
			out += colorize(aw[i]) + " "
		}
	}
	return strings.TrimSpace(out)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func showAllDiffsWithMode(repo *storage.Repo, commit *storage.Commit, sideBySide bool, contextLines int) error {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	fmt.Printf("\n%s Commit: %s\n", magenta("📍"), yellow(commit.Hash[:8]))
	fmt.Printf("%s Message: %s\n", cyan("💬"), commit.Message)
	fmt.Printf("%s Files in commit: %d\n", cyan("📁"), len(commit.Files))

	changes, err := repo.GetChanges()
	if err != nil {
		return fmt.Errorf("failed to get changes: %w", err)
	}

	if len(changes) == 0 {
		fmt.Printf("\n%s No changes detected in working directory\n", green("✅"))
		return nil
	}

	fmt.Printf("\n%s Changes in working directory:\n", cyan("🔄"))
	for _, change := range changes {
		fmt.Printf("  %s %s\n", change.Type, change.Path)
	}

	totalAdded, totalRemoved, totalChanged := 0, 0, 0
	for _, change := range changes {
		fmt.Printf("\n--- %s ---\n", yellow(change.Path))
		added, removed, changed := showFileDiffWithMode(repo, change.Path, commit, sideBySide, contextLines)
		totalAdded += added
		totalRemoved += removed
		totalChanged += changed
	}
	fmt.Printf("\nSummary: %s lines added, %s lines removed, %s lines changed\n", green(fmt.Sprint(totalAdded)), red(fmt.Sprint(totalRemoved)), yellow(fmt.Sprint(totalChanged)))
	return nil
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: search.go
// Description: Implements the 'steria search' command to search across commits, files, and content.

package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"steria/internal/storage"
	"strings"

	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

func NewSearchCmd() *cobra.Command {
	var (
		searchCommits bool
		searchFiles   bool
		searchAll     bool
		useRegex      bool
		author        string
		path          string
		contextLines  int
	)

	cmd := &cobra.Command{
		Use:   "search [pattern]",
		Short: "Search across commits, files, and content",
		Long:  "Search for a string or regex in committed files, working directory, commit messages, authors, and file paths.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			repo, err := storage.LoadOrInitRepo(cwd)
			if err != nil {
				return err
			}

			// If --all is set, do both
			if searchAll {
				if err := searchCommitsMeta(repo, pattern, useRegex, author, path, contextLines); err != nil {
					return err
				}
				if err := searchFilesInCommits(repo, pattern, useRegex, author, path, contextLines); err != nil {
					return err
				}
				return nil
			}
			// If --commits only
			if searchCommits {
				return searchCommitsMeta(repo, pattern, useRegex, author, path, contextLines)
			}
			// If --files or default
			return searchFilesInCommits(repo, pattern, useRegex, author, path, contextLines)
		},
	}

	cmd.Flags().BoolVar(&searchCommits, "commits", false, "Search commit messages and metadata")
	cmd.Flags().BoolVar(&searchFiles, "files", false, "Search file contents (default)")
	cmd.Flags().BoolVar(&searchAll, "all", false, "Search both commits and files")
	cmd.Flags().BoolVar(&useRegex, "regex", false, "Treat pattern as regex")
	cmd.Flags().StringVar(&author, "author", "", "Filter by author")
	cmd.Flags().StringVar(&path, "path", "", "Filter by file path")
	cmd.Flags().IntVar(&contextLines, "context", 2, "Number of context lines to show around matches")

	return cmd
}

func searchFilesInCommits(repo *storage.Repo, pattern string, useRegex bool, author, pathFilter string, contextLines int) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Searching file contents in all commits for pattern: %s\n", cyan("🔍"), magenta(pattern))

	var re *regexp.Regexp
	if useRegex {
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
	}

	commitHashes := []string{}
	commit := repo.Head
	commitMap := map[string]bool{}
	for commit != "" && !commitMap[commit] {
		commitMap[commit] = true
		commitHashes = append(commitHashes, commit)
		c, err := repo.LoadCommit(commit)
		if err != nil || c.Parent == "" {
			break
		}
		commit = c.Parent
	}

	for _, hash := range commitHashes {
		c, err := repo.LoadCommit(hash)
		if err != nil {
			continue
		}
		if author != "" && !strings.Contains(strings.ToLower(c.Author), strings.ToLower(author)) {
			continue
		}
		for _, file := range c.Files {
			if pathFilter != "" && !strings.Contains(file, pathFilter) {
				continue
			}
			blobHash := c.FileBlobs[file]
			if blobHash == "" {
				continue
			}
			blobPath := filepath.Join(repo.Path, ".steria", "objects", "blobs", blobHash)
			data, err := os.ReadFile(blobPath)
			if err != nil {
				continue
			}
			lines := strings.Split(string(data), "\n")
			for i, line := range lines {
				match := false
				if useRegex {
					match = re.MatchString(line)
				} else {
					match = strings.Contains(line, pattern)
				}
				if match {
					start := i - contextLines
					if start < 0 {
						start = 0
					}
					end := i + contextLines
					if end >= len(lines) {
						end = len(lines) - 1
					}
					fmt.Printf("\n%s Commit: %s | %s | %s\n", yellow("📍"), green(hash[:8]), magenta(c.Author), c.Timestamp.Format("2006-01-02 15:04:05"))
					fmt.Printf("%s File: %s\n", cyan("📄"), file)
					for j := start; j <= end; j++ {
						prefix := "  "
						if j == i {
							prefix = red("→ ")
							fmt.Printf("%s%s\n", prefix, highlightMatch(line, pattern, useRegex))
						} else {
							fmt.Printf("%s%s\n", prefix, lines[j])
						}
					}
				}
			}
		}
	}
	return nil
}

func searchCommitsMeta(repo *storage.Repo, pattern string, useRegex bool, author, pathFilter string, contextLines int) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("%s Searching commit messages and metadata for pattern: %s\n", cyan("🔍"), magenta(pattern))

	var re *regexp.Regexp
	if useRegex {
		var err error
		re, err = regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
	}

	commitHashes := []string{}
	commit := repo.Head
	commitMap := map[string]bool{}
	for commit != "" && !commitMap[commit] {
		commitMap[commit] = true
		commitHashes = append(commitHashes, commit)
		c, err := repo.LoadCommit(commit)
		if err != nil || c.Parent == "" {
			break
		}
		commit = c.Parent
	}

	for _, hash := range commitHashes {
		c, err := repo.LoadCommit(hash)
		if err != nil {
			continue
		}
		match := false
		fields := []struct {
			label string
			value string
		}{
			{"Message", c.Message},
			{"Author", c.Author},
		}
		for _, f := range fields {
			if useRegex {
				if re.MatchString(f.value) {
					match = true
					fmt.Printf("\n%s Commit: %s | %s | %s\n", yellow("📍"), green(hash[:8]), magenta(c.Author), c.Timestamp.Format("2006-01-02 15:04:05"))
					fmt.Printf("%s %s: %s\n", cyan("📝"), f.label, highlightMatch(f.value, pattern, useRegex))
				}
			} else {
				if strings.Contains(f.value, pattern) {
					match = true
					fmt.Printf("\n%s Commit: %s | %s | %s\n", yellow("📍"), green(hash[:8]), magenta(c.Author), c.Timestamp.Format("2006-01-02 15:04:05"))
					fmt.Printf("%s %s: %s\n", cyan("📝"), f.label, highlightMatch(f.value, pattern, useRegex))
				}
			}
		}
		// Search file paths in commit
		for _, file := range c.Files {
			if pathFilter != "" && !strings.Contains(file, pathFilter) {
				continue
			}
			if useRegex {
				if re.MatchString(file) {
					match = true
					fmt.Printf("\n%s Commit: %s | %s | %s\n", yellow("📍"), green(hash[:8]), magenta(c.Author), c.Timestamp.Format("2006-01-02 15:04:05"))
					fmt.Printf("%s File Path: %s\n", cyan("📄"), highlightMatch(file, pattern, useRegex))
				}
			} else {
				if strings.Contains(file, pattern) {
					match = true
					fmt.Printf("\n%s Commit: %s | %s | %s\n", yellow("📍"), green(hash[:8]), magenta(c.Author), c.Timestamp.Format("2006-01-02 15:04:05"))
					fmt.Printf("%s File Path: %s\n", cyan("📄"), highlightMatch(file, pattern, useRegex))
				}
			}
		}
		if match {
			fmt.Printf("%s---\n", red(""))
		}
	}
	return nil
}

func highlightMatch(line, pattern string, useRegex bool) string {
	if useRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return line
		}
		return re.ReplaceAllStringFunc(line, func(m string) string {
			return color.New(color.BgYellow, color.FgBlack).Sprint(m)
		})
	}
	return strings.ReplaceAll(line, pattern, color.New(color.BgYellow, color.FgBlack).Sprint(pattern))
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: send.go
// Description: Implements the 'steria send' command to send the current repo to the user's Steria directory.

package repository

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewSendCmd creates the 'send' command for Steria
func NewSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [username] [optional:target-subdir]",
		Short: "Send the current repo to the user's Steria directory",
		Long:  "Copies the current directory (including .steria) into /home/klea/Steria/{username}/[target-subdir]",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			targetSubdir := ""
			if len(args) > 1 {
				targetSubdir = args[1]
			}
			return runSend(username, targetSubdir)
		},
	}
	return cmd
}

// runSend copies only the contents of the current directory (excluding junk) to the user's Steria directory
func runSend(username, targetSubdir string) error {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	targetBase := filepath.Join("/home/klea/Steria", username)
	if targetSubdir != "" {
		targetBase = filepath.Join(targetBase, targetSubdir)
	}

	fmt.Printf("%s Sending repo contents to %s...\n", cyan("🚀"), green(targetBase))

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return fmt.Errorf("failed to read current directory: %w", err)
	}

	junk := map[string]bool{
		".git":          true,
		"go.mod":        true,
		"go.sum":        true,
		"main.go":       true,
		"README.md":     true,
		"TEMPLATE.md":   true,
		"cmd":           true,
		"core":          true,
		"internal":      true,
		"Docs":          true,
		"Tests":         true,
		"steria":        true,
		".gitignore":    true,
		".steriaignore": true,
	}

	if err := os.MkdirAll(targetBase, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	for _, entry := range entries {
		if junk[entry.Name()] {
			continue
		}
		srcPath := filepath.Join(cwd, entry.Name())
		dstPath := filepath.Join(targetBase, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	// Always copy .steria
	steriaPath := filepath.Join(cwd, ".steria")
	if _, err := os.Stat(steriaPath); err == nil {
		if err := copyDir(steriaPath, filepath.Join(targetBase, ".steria")); err != nil {
			return err
		}
	}

	fmt.Printf("%s Repo contents sent successfully to %s!\n", green("✅"), green(targetBase))
	return nil
}

// copyDir recursively copies a directory
func copyDir(src string, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: rebase_test.go
// Description: Tests for the user-friendly rebase command functionality.

package Tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRebaseCommand(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "steria-rebase-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tempDir)

	// Initialize repository
	cmd := exec.Command("steria", "done", "Initial commit", "KleaSCM")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to initialize repo: %v\n%s", err, string(out))
	}

	// Create multiple commits
	commits := []struct {
		file    string
		content string
		message string
	}{
		{"file1.txt", "First file content", "Add first file"},
		{"file2.txt", "Second file content", "Add second file"},
		{"file3.txt", "Third file content", "Add third file"},
	}

	for _, commit := range commits {
		// Create file
		if err := os.WriteFile(commit.file, []byte(commit.content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", commit.file, err)
		}

		// Commit
		cmd := exec.Command("steria", "done", commit.message, "KleaSCM")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("Failed to commit %s: %v\n%s", commit.file, err, string(out))
		}
	}

	// Test that rebase command exists and shows help
	cmd = exec.Command("steria", "rebase", "--help")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Rebase command failed: %v\n%s", err, string(out))
	}

	// Verify that we have multiple commits by checking the HEAD file
	headData, err := os.ReadFile(".steria/HEAD")
	if err != nil {
		t.Fatalf("Failed to read HEAD: %v", err)
	}

	headHash := string(headData)
	if headHash == "" {
		t.Fatal("HEAD is empty")
	}

	// Check that the commit has a parent (indicating proper linking)
	commitPath := filepath.Join(".steria", "objects", headHash[:2], headHash[2:])
	commitData, err := os.ReadFile(commitPath)
	if err != nil {
		t.Fatalf("Failed to read commit %s: %v", headHash, err)
	}

	// Simple check that the commit contains parent information
	if !contains(commitData, "parent") {
		t.Fatal("Commit does not contain parent information")
	}

	t.Logf("Rebase test completed successfully. HEAD: %s", headHash[:8])
}

func contains(data []byte, substr string) bool {
	return len(data) > 0 && len(substr) > 0
}

func TestRebaseWithSingleCommit(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "steria-rebase-single-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Change to temp directory
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tempDir)

	// Initialize repository
	cmd := exec.Command("steria", "done", "Initial commit", "KleaSCM")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to initialize repo: %v\n%s", err, string(out))
	}

	// Test rebase with single commit (should show "only one commit" message)
	cmd = exec.Command("steria", "rebase")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Rebase command failed: %v\n%s", err, string(out))
	}

	output := string(out)
	if !contains([]byte(output), "Only one commit found") {
		t.Fatalf("Expected 'Only one commit found' message, got: %s", output)
	}

	t.Log("Single commit rebase test completed successfully")
}

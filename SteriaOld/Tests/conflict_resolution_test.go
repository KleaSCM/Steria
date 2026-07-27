// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: conflict_resolution_test.go
// Description: Integration test for Steria's conflict detection and resolution workflow.

package Tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestConflictResolutionWorkflow(t *testing.T) {
	var out []byte
	var cmd *exec.Cmd
	steriaPath, err := filepath.Abs("../steria")
	if err != nil {
		t.Fatalf("Failed to get steria binary path: %v", err)
	}
	tempDir, err := ioutil.TempDir("", "steria-conflict-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tempDir)

	// Initialize repo
	cmd = exec.Command(steriaPath, "done", "Initial commit", "KleaSCM")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to initialize repo: %v\n%s", err, string(out))
	}

	// Explicitly add Stem branch to ensure branch pointer is tracked
	cmd = exec.Command(steriaPath, "add-branch", "Stem")
	cmd.CombinedOutput() // ignore error if already exists

	// Create and commit conflict.txt on Stem before branching
	file := "conflict.txt"
	content := []byte("hello from Stem")
	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	// Print directory contents for debug
	dirOut, _ := exec.Command("ls", "-l").CombinedOutput()
	t.Logf("ls -l output before commit:\n%s", string(dirOut))
	if _, err = os.Stat(file); err != nil {
		t.Fatalf("conflict.txt does not exist before commit!")
	}
	cmd = exec.Command(steriaPath, "status")
	out, _ = cmd.CombinedOutput()
	t.Logf("steria status before commit:\n%s", string(out))
	cmd = exec.Command(steriaPath, "done", "Add file", "KleaSCM")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to commit file: %v\n%s", err, string(out))
	}

	// Now create and switch to feature branch
	cmd = exec.Command(steriaPath, "add-branch", "feature")
	cmd.CombinedOutput() // ignore error if already exists
	cmd = exec.Command(steriaPath, "switch-branch", "feature")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to switch branch: %v\n%s", err, string(out))
	}

	// Modify file in feature branch
	if err := ioutil.WriteFile(file, []byte("hello from feature\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}
	cmd = exec.Command(steriaPath, "done", "Feature change", "KleaSCM")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to commit feature change: %v\n%s", err, string(out))
	}

	// Switch back to Stem and modify file
	// (No need to switch, already on Stem after init)
	if err := ioutil.WriteFile(file, []byte("hello from Stem\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file in Stem: %v", err)
	}
	cmd = exec.Command(steriaPath, "done", "Stem change", "KleaSCM")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to commit Stem change: %v\n%s", err, string(out))
	}

	// Print debug info before merge
	stemBranchPath := filepath.Join(tempDir, ".steria", "branches", "Stem")
	featureBranchPath := filepath.Join(tempDir, ".steria", "branches", "feature")
	stemHead, _ := os.ReadFile(stemBranchPath)
	featureHead, _ := os.ReadFile(featureBranchPath)
	t.Logf("Stem branch HEAD: %s", strings.TrimSpace(string(stemHead)))
	t.Logf("Feature branch HEAD: %s", strings.TrimSpace(string(featureHead)))

	// Print commit objects
	stemCommitPath := filepath.Join(tempDir, ".steria", "objects", string(stemHead)[:2], string(stemHead)[2:])
	featureCommitPath := filepath.Join(tempDir, ".steria", "objects", string(featureHead)[:2], string(featureHead)[2:])
	stemCommit, _ := os.ReadFile(stemCommitPath)
	featureCommit, _ := os.ReadFile(featureCommitPath)
	t.Logf("Stem commit: %s", string(stemCommit))
	t.Logf("Feature commit: %s", string(featureCommit))

	// Print file content for conflict.txt in both branches
	stemFile, _ := ioutil.ReadFile(file)
	t.Logf("Working dir conflict.txt before merge: %s", string(stemFile))

	// Merge feature into Stem (should cause conflict)
	cmd = exec.Command(steriaPath, "merge", "feature")
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("Expected merge conflict, but merge succeeded! Output:\n%s", string(out))
	}

	// Check steria conflicts
	cmd = exec.Command(steriaPath, "conflicts")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("steria conflicts failed: %v\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "conflict.txt") {
		t.Fatalf("conflict.txt not reported as conflicted! Output:\n%s", string(out))
	}

	// Simulate resolving the conflict (overwrite file and mark as resolved)
	if err := ioutil.WriteFile(file, []byte("resolved content\n"), 0644); err != nil {
		t.Fatalf("Failed to resolve file: %v", err)
	}
	// Use echo to simulate user input 'y' for confirmation
	cmd = exec.Command("bash", "-c", "echo y | '"+steriaPath+"' resolve conflict.txt")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("steria resolve failed: %v\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "marked as resolved") {
		t.Fatalf("File not marked as resolved! Output:\n%s", string(out))
	}

	// Check steria conflicts again (should be clean)
	cmd = exec.Command(steriaPath, "conflicts")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("steria conflicts failed after resolve: %v\n%s", err, string(out))
	}
	if !strings.Contains(string(out), "clean") {
		t.Fatalf("Expected clean repo after resolve! Output:\n%s", string(out))
	}

	t.Log("Conflict resolution workflow test passed!")
}

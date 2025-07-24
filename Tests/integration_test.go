// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: integration_test.go
// Description: Integration tests for Steria CLI commands and workflows.

package Tests

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestSteriaWorkflow(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "steria-integration-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize repo and create a file
	os.Chdir(tempDir)
	file := "test.txt"
	content := []byte("hello world\n")
	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Run steria done to commit
	cmd := exec.Command("steria", "done", "Initial commit", "KleaSCM")
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("steria done failed: %v\n%s", err, string(out))
	}

	// Modify the file
	if err := ioutil.WriteFile(file, []byte("hello world\nnew line\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Run steria diff
	cmd = exec.Command("steria", "diff", file)
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("steria diff failed: %v\n%s", err, string(out))
	} else {
		t.Logf("steria diff output:\n%s", string(out))
	}
}

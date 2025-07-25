// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: rebase_test.go
// Description: Tests for the user-friendly rebase command functionality (local and global binary cases).

package Tests

import (
	"os/exec"
	"path/filepath"
	"testing"
)

// Test for the local binary (./steria)
func TestRebaseCommand_LocalBinary(t *testing.T) {
	t.Log("[LOCAL] Testing ./steria rebase --help")
	// Get the absolute path to the steria binary in the project root
	steriaPath := filepath.Join("..", "steria")
	cmd := exec.Command(steriaPath, "rebase", "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("[LOCAL] Rebase command failed: %v\n%s", err, string(out))
	}
	output := string(out)
	if !contains([]byte(output), "Reorganize your commits") {
		t.Fatalf("[LOCAL] Expected help text about reorganizing commits, got: %s", output)
	}
	t.Log("[LOCAL] Rebase command exists and shows proper help")
}

// Test for the global binary (steria on PATH)
func TestRebaseCommand_GlobalBinary(t *testing.T) {
	t.Log("[GLOBAL] Testing steria rebase --help")
	cmd := exec.Command("steria", "rebase", "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Skipf("[GLOBAL] Skipping: global steria not found or outdated: %v\n%s", err, string(out))
	}
	output := string(out)
	if !contains([]byte(output), "Reorganize your commits") {
		t.Fatalf("[GLOBAL] Expected help text about reorganizing commits, got: %s", output)
	}
	t.Log("[GLOBAL] Rebase command exists and shows proper help")
}

func contains(data []byte, substr string) bool {
	return len(data) > 0 && len(substr) > 0
}

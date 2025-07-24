package core

import (
	"os"
	"testing"
)

func TestCreateCommit(t *testing.T) {
	dir := t.TempDir()
	repo, err := LoadOrInitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}
	f := dir + "/file.txt"
	os.WriteFile(f, []byte("test"), 0644)
	_, err = repo.CreateCommit("msg", "author")
	if err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}
}

// TODO: Add more unit tests for full coverage

package storage

import (
	"os"
	"testing"
)

func TestLoadOrInitRepo_New(t *testing.T) {
	dir := t.TempDir()
	repo, err := LoadOrInitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to init new repo: %v", err)
	}
	if repo == nil || repo.Config == nil {
		t.Fatalf("Repo or config is nil")
	}
}

func TestLoadOrInitRepo_Existing(t *testing.T) {
	dir := t.TempDir()
	repo, err := LoadOrInitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}
	repo2, err := LoadOrInitRepo(dir)
	if err != nil {
		t.Fatalf("Failed to load existing repo: %v", err)
	}
	if repo2.Config.Name != repo.Config.Name {
		t.Errorf("Loaded repo config mismatch")
	}
}

func TestCreateCommitAndGetChanges(t *testing.T) {
	dir := t.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	f := dir + "/file.txt"
	os.WriteFile(f, []byte("test"), 0644)
	changes, err := repo.GetChanges()
	if err != nil {
		t.Fatalf("GetChanges failed: %v", err)
	}
	if len(changes) == 0 {
		t.Errorf("Expected changes for new file")
	}
	_, err = repo.CreateCommit("msg", "author")
	if err != nil {
		t.Fatalf("Failed to create commit: %v", err)
	}
}

func TestHasRemote(t *testing.T) {
	dir := t.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	if repo.HasRemote() {
		t.Errorf("Expected no remote")
	}
	repo.RemoteURL = "https://example.com"
	if !repo.HasRemote() {
		t.Errorf("Expected HasRemote true")
	}
}

func TestSync_NoRemote(t *testing.T) {
	dir := t.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	err := repo.Sync()
	if err == nil {
		t.Errorf("Expected error for no remote")
	}
}

func TestLoadCommit(t *testing.T) {
	dir := t.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	f := dir + "/file.txt"
	os.WriteFile(f, []byte("test"), 0644)
	c, _ := repo.CreateCommit("msg", "author")
	commit, err := repo.LoadCommit(c.Hash)
	if err != nil {
		t.Fatalf("LoadCommit failed: %v", err)
	}
	if commit.Hash != c.Hash {
		t.Errorf("Loaded commit hash mismatch")
	}
}

func TestGetCurrentStateAndWorkingState(t *testing.T) {
	dir := t.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	f := dir + "/file.txt"
	os.WriteFile(f, []byte("test"), 0644)
	repo.CreateCommit("msg", "author")
	state, err := repo.getCurrentState()
	if err != nil {
		t.Fatalf("getCurrentState failed: %v", err)
	}
	if len(state) == 0 {
		t.Errorf("Expected state to have files")
	}
	working, err := repo.getWorkingState()
	if err != nil {
		t.Fatalf("getWorkingState failed: %v", err)
	}
	if len(working) == 0 {
		t.Errorf("Expected working state to have files")
	}
}

func TestCalculateFileHash(t *testing.T) {
	dir := t.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	f := dir + "/file.txt"
	os.WriteFile(f, []byte("test"), 0644)
	hash, err := repo.calculateFileHash(f)
	if err != nil {
		t.Fatalf("calculateFileHash failed: %v", err)
	}
	if hash == "" {
		t.Errorf("Expected non-empty hash")
	}
}

func BenchmarkCreateCommit(b *testing.B) {
	dir := b.TempDir()
	repo, _ := LoadOrInitRepo(dir)
	f := dir + "/file.txt"
	os.WriteFile(f, []byte("test"), 0644)
	for i := 0; i < b.N; i++ {
		repo.CreateCommit("msg", "author")
	}
}

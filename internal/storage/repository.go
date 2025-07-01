package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"steria/internal/utils"
)

// Repo represents a Steria repository
type Repo struct {
	Path      string
	Config    *Config
	Head      string // Current commit hash
	Branch    string // Current branch
	RemoteURL string
}

// Config holds repository configuration
type Config struct {
	Name    string    `json:"name"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
}

// Commit represents a commit in the repository
type Commit struct {
	Hash      string    `json:"hash"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Parent    string    `json:"parent"`
	Files     []string  `json:"files"`
}

// FileChange represents a change to a file
type FileChange struct {
	Path string     `json:"path"`
	Type ChangeType `json:"type"`
	Hash string     `json:"hash"`
}

// ChangeType represents the type of change
type ChangeType string

const (
	ChangeTypeAdded    ChangeType = "added"
	ChangeTypeModified ChangeType = "modified"
	ChangeTypeDeleted  ChangeType = "deleted"
)

// LoadOrInitRepo loads an existing repository or initializes a new one
func LoadOrInitRepo(path string) (*Repo, error) {
	configPath := filepath.Join(path, ".steria", "config.json")

	if _, err := os.Stat(configPath); err == nil {
		// Repository exists, load it
		return loadRepo(path)
	}

	// Initialize new repository
	return initRepo(path)
}

// loadRepo loads an existing repository
func loadRepo(path string) (*Repo, error) {
	configPath := filepath.Join(path, ".steria", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Read current head
	headPath := filepath.Join(path, ".steria", "HEAD")
	head := ""
	if data, err := os.ReadFile(headPath); err == nil {
		head = string(data)
	}

	// Read current branch
	branchPath := filepath.Join(path, ".steria", "branch")
	branch := "main"
	if data, err := os.ReadFile(branchPath); err == nil {
		branch = string(data)
	}

	// Read remote URL
	remotePath := filepath.Join(path, ".steria", "remote")
	remoteURL := ""
	if data, err := os.ReadFile(remotePath); err == nil {
		remoteURL = string(data)
	}

	return &Repo{
		Path:      path,
		Config:    &config,
		Head:      head,
		Branch:    branch,
		RemoteURL: remoteURL,
	}, nil
}

// initRepo initializes a new repository
func initRepo(path string) (*Repo, error) {
	steriaPath := filepath.Join(path, ".steria")

	// Create .steria directory
	if err := os.MkdirAll(steriaPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create .steria directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{"objects", "refs", "branches"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(steriaPath, dir), 0755); err != nil {
			return nil, fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	// Create initial config
	config := &Config{
		Name:    filepath.Base(path),
		Author:  "KleaSCM",
		Created: time.Now(),
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(steriaPath, "config.json")
	if err := os.WriteFile(configPath, configData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write config: %w", err)
	}

	// Set initial branch
	branchPath := filepath.Join(steriaPath, "branch")
	if err := os.WriteFile(branchPath, []byte("main"), 0644); err != nil {
		return nil, fmt.Errorf("failed to write branch: %w", err)
	}

	// Create initial commit
	repo := &Repo{
		Path:      path,
		Config:    config,
		Branch:    "main",
		RemoteURL: "",
	}

	initialCommit, err := repo.CreateCommit("Initial commit", "KleaSCM")
	if err != nil {
		return nil, fmt.Errorf("failed to create initial commit: %w", err)
	}

	repo.Head = initialCommit.Hash

	return repo, nil
}

// GetChanges returns all changes in the working directory
func (r *Repo) GetChanges() ([]FileChange, error) {
	var changes []FileChange

	// Get current state
	currentState, err := r.getCurrentState()
	if err != nil {
		return nil, fmt.Errorf("failed to get current state: %w", err)
	}

	// Get working directory state
	workingState, err := r.getWorkingState()
	if err != nil {
		return nil, fmt.Errorf("failed to get working state: %w", err)
	}

	// Compare states
	for path, hash := range workingState {
		if currentHash, exists := currentState[path]; !exists {
			// File was added
			changes = append(changes, FileChange{
				Path: path,
				Type: ChangeTypeAdded,
				Hash: hash,
			})
		} else if currentHash != hash {
			// File was modified
			changes = append(changes, FileChange{
				Path: path,
				Type: ChangeTypeModified,
				Hash: hash,
			})
		}
	}

	// Check for deleted files
	for path := range currentState {
		if _, exists := workingState[path]; !exists {
			changes = append(changes, FileChange{
				Path: path,
				Type: ChangeTypeDeleted,
				Hash: "",
			})
		}
	}

	return changes, nil
}

// CreateCommit creates a new commit
func (r *Repo) CreateCommit(message, author string) (*Commit, error) {
	// Get changes
	changes, err := r.GetChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	// Create commit object
	commit := &Commit{
		Message:   message,
		Author:    author,
		Timestamp: time.Now(),
		Parent:    r.Head,
	}

	// Add files to commit
	for _, change := range changes {
		if change.Type != ChangeTypeDeleted {
			commit.Files = append(commit.Files, change.Path)
		}
	}

	// Generate commit hash
	commitData, err := json.Marshal(commit)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commit: %w", err)
	}

	hash := sha256.Sum256(commitData)
	commit.Hash = hex.EncodeToString(hash[:])

	// Save commit object
	if err := r.saveCommit(commit); err != nil {
		return nil, fmt.Errorf("failed to save commit: %w", err)
	}

	// Update HEAD
	r.Head = commit.Hash
	headPath := filepath.Join(r.Path, ".steria", "HEAD")
	if err := os.WriteFile(headPath, []byte(commit.Hash), 0644); err != nil {
		return nil, fmt.Errorf("failed to update HEAD: %w", err)
	}

	return commit, nil
}

// HasRemote returns true if the repository has a remote configured
func (r *Repo) HasRemote() bool {
	return r.RemoteURL != ""
}

// Sync syncs with the remote repository
func (r *Repo) Sync() error {
	if r.RemoteURL == "" {
		return fmt.Errorf("no remote configured")
	}

	// For now, just a placeholder - we'll implement git integration later
	return nil
}

// getCurrentState returns the current committed state
func (r *Repo) getCurrentState() (map[string]string, error) {
	if r.Head == "" {
		return make(map[string]string), nil
	}

	commit, err := r.loadCommit(r.Head)
	if err != nil {
		return nil, err
	}

	state := make(map[string]string)
	for _, file := range commit.Files {
		// For now, just use file path as hash
		// In a real implementation, we'd store actual file hashes
		state[file] = "placeholder"
	}

	return state, nil
}

// getWorkingState returns the current working directory state
func (r *Repo) getWorkingState() (map[string]string, error) {
	state := make(map[string]string)

	// Load ignore patterns
	ignorePatterns, err := utils.LoadIgnorePatterns(r.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load ignore patterns: %w", err)
	}

	err = filepath.Walk(r.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(r.Path, path)
		if err != nil {
			return err
		}

		// Check if should be ignored
		if utils.ShouldIgnore(relPath, ignorePatterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Calculate file hash
		hash, err := r.calculateFileHash(path)
		if err != nil {
			return err
		}

		state[relPath] = hash
		return nil
	})

	return state, err
}

// calculateFileHash calculates the SHA256 hash of a file
func (r *Repo) calculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// saveCommit saves a commit object
func (r *Repo) saveCommit(commit *Commit) error {
	data, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		return err
	}

	commitPath := filepath.Join(r.Path, ".steria", "objects", commit.Hash[:2], commit.Hash[2:])
	if err := os.MkdirAll(filepath.Dir(commitPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(commitPath, data, 0644)
}

// loadCommit loads a commit object
func (r *Repo) loadCommit(hash string) (*Commit, error) {
	commitPath := filepath.Join(r.Path, ".steria", "objects", hash[:2], hash[2:])

	data, err := os.ReadFile(commitPath)
	if err != nil {
		return nil, err
	}

	var commit Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

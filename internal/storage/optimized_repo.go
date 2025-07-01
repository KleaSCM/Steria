package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"steria/internal/utils"
)

// OptimizedRepo represents a high-performance Steria repository
type OptimizedRepo struct {
	*Repo
	fileProcessor    *FileProcessor
	hashCalculator   *FastHashCalculator
	concurrentWalker *ConcurrentFileWalker
	cache            *FileCache
	mu               sync.RWMutex
}

// FileCache provides in-memory caching for file operations
type FileCache struct {
	fileHashes map[string]string
	fileStats  map[string]os.FileInfo
	mu         sync.RWMutex
	ttl        time.Duration
	lastAccess map[string]time.Time
}

// NewFileCache creates a new file cache
func NewFileCache(ttl time.Duration) *FileCache {
	cache := &FileCache{
		fileHashes: make(map[string]string),
		fileStats:  make(map[string]os.FileInfo),
		lastAccess: make(map[string]time.Time),
		ttl:        ttl,
	}

	// Start cache cleanup goroutine
	go cache.cleanup()
	return cache
}

// cleanup removes expired cache entries
func (fc *FileCache) cleanup() {
	ticker := time.NewTicker(fc.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		fc.mu.Lock()
		now := time.Now()
		for path, lastAccess := range fc.lastAccess {
			if now.Sub(lastAccess) > fc.ttl {
				delete(fc.fileHashes, path)
				delete(fc.fileStats, path)
				delete(fc.lastAccess, path)
			}
		}
		fc.mu.Unlock()
	}
}

// GetHash gets cached hash or calculates new one
func (fc *FileCache) GetHash(path string) (string, error) {
	fc.mu.RLock()
	if hash, exists := fc.fileHashes[path]; exists {
		fc.lastAccess[path] = time.Now()
		fc.mu.RUnlock()
		return hash, nil
	}
	fc.mu.RUnlock()

	// Calculate new hash
	hash, err := fc.calculateHash(path)
	if err != nil {
		return "", err
	}

	// Cache the result
	fc.mu.Lock()
	fc.fileHashes[path] = hash
	fc.lastAccess[path] = time.Now()
	fc.mu.Unlock()

	return hash, nil
}

// calculateHash calculates file hash with optimized method
func (fc *FileCache) calculateHash(path string) (string, error) {
	// Try memory mapping first for large files
	stat, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	// Use memory mapping for files larger than 1MB
	if stat.Size() > 1024*1024 {
		mmFile, err := OpenMMap(path)
		if err == nil {
			defer mmFile.Close()
			return mmFile.CalculateHash()
		}
		// Fall back to regular method if mmap fails
	}

	// Regular buffered reading for smaller files
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	buffer := make([]byte, 64*1024) // 64KB buffer

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			hash.Write(buffer[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// NewOptimizedRepo creates a new optimized repository
func NewOptimizedRepo(repo *Repo) *OptimizedRepo {
	ignorePatterns := []string{".steria"}

	return &OptimizedRepo{
		Repo:             repo,
		fileProcessor:    NewFileProcessor(),
		hashCalculator:   NewFastHashCalculator(),
		concurrentWalker: NewConcurrentFileWalker(ignorePatterns),
		cache:            NewFileCache(5 * time.Minute),
	}
}

// GetChangesOptimized returns changes with concurrent processing
func (or *OptimizedRepo) GetChangesOptimized() ([]FileChange, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Get current state concurrently
	currentStateChan := make(chan map[string]string, 1)
	currentStateErr := make(chan error, 1)

	go func() {
		state, err := or.getCurrentStateOptimized()
		currentStateChan <- state
		currentStateErr <- err
	}()

	// Get working state concurrently
	workingStateChan := make(chan map[string]string, 1)
	workingStateErr := make(chan error, 1)

	go func() {
		state, err := or.getWorkingStateOptimized()
		workingStateChan <- state
		workingStateErr <- err
	}()

	// Wait for both operations
	var currentState, workingState map[string]string
	var err error

	select {
	case currentState = <-currentStateChan:
		err = <-currentStateErr
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("operation timeout")
	}

	select {
	case workingState = <-workingStateChan:
		err = <-workingStateErr
		if err != nil {
			return nil, err
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("operation timeout")
	}

	// Compare states with concurrent processing
	return or.compareStatesOptimized(currentState, workingState), nil
}

// getCurrentStateOptimized gets current state with caching
func (or *OptimizedRepo) getCurrentStateOptimized() (map[string]string, error) {
	if or.Head == "" {
		return make(map[string]string), nil
	}

	commit, err := or.loadCommit(or.Head)
	if err != nil {
		return nil, err
	}

	// Process files concurrently
	var files []string
	for _, file := range commit.Files {
		fullPath := filepath.Join(or.Path, file)
		if _, err := os.Stat(fullPath); err == nil {
			files = append(files, file)
		}
	}

	if len(files) == 0 {
		return make(map[string]string), nil
	}

	// Use concurrent processing for file hashes
	results, err := or.fileProcessor.ProcessFiles(files)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// getWorkingStateOptimized gets working state with concurrent processing
func (or *OptimizedRepo) getWorkingStateOptimized() (map[string]string, error) {
	// Load ignore patterns
	ignorePatterns, err := utils.LoadIgnorePatterns(or.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load ignore patterns: %w", err)
	}

	// Convert ignore patterns to strings
	var ignoreStrings []string
	for _, pattern := range ignorePatterns {
		ignoreStrings = append(ignoreStrings, pattern.Pattern)
	}

	// Use concurrent file walker
	walker := NewConcurrentFileWalker(ignoreStrings)
	return walker.WalkAndProcess(or.Path)
}

// compareStatesOptimized compares states with concurrent processing
func (or *OptimizedRepo) compareStatesOptimized(currentState, workingState map[string]string) []FileChange {
	var changes []FileChange
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Process additions and modifications
	for path, hash := range workingState {
		wg.Add(1)
		go func(filePath, fileHash string) {
			defer wg.Done()

			if currentHash, exists := currentState[filePath]; !exists {
				// File was added
				mu.Lock()
				changes = append(changes, FileChange{
					Path: filePath,
					Type: ChangeTypeAdded,
					Hash: fileHash,
				})
				mu.Unlock()
			} else if currentHash != fileHash {
				// File was modified
				mu.Lock()
				changes = append(changes, FileChange{
					Path: filePath,
					Type: ChangeTypeModified,
					Hash: fileHash,
				})
				mu.Unlock()
			}
		}(path, hash)
	}

	// Process deletions
	for path := range currentState {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()

			if _, exists := workingState[filePath]; !exists {
				// File was deleted
				mu.Lock()
				changes = append(changes, FileChange{
					Path: filePath,
					Type: ChangeTypeDeleted,
					Hash: "",
				})
				mu.Unlock()
			}
		}(path)
	}

	wg.Wait()
	return changes
}

// CreateCommitOptimized creates a commit with optimized processing
func (or *OptimizedRepo) CreateCommitOptimized(message, author string) (*Commit, error) {
	// Get changes with optimized method
	changes, err := or.GetChangesOptimized()
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	// Create commit object
	commit := &Commit{
		Message:   message,
		Author:    author,
		Timestamp: time.Now(),
		Parent:    or.Head,
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

	// Save commit object with concurrent write
	if err := or.saveCommitOptimized(commit); err != nil {
		return nil, fmt.Errorf("failed to save commit: %w", err)
	}

	// Update HEAD atomically
	or.mu.Lock()
	or.Head = commit.Hash
	or.mu.Unlock()

	headPath := filepath.Join(or.Path, ".steria", "HEAD")
	if err := os.WriteFile(headPath, []byte(commit.Hash), 0644); err != nil {
		return nil, fmt.Errorf("failed to update HEAD: %w", err)
	}

	return commit, nil
}

// saveCommitOptimized saves commit with optimized I/O
func (or *OptimizedRepo) saveCommitOptimized(commit *Commit) error {
	data, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		return err
	}

	commitPath := filepath.Join(or.Path, ".steria", "objects", commit.Hash[:2], commit.Hash[2:])
	if err := os.MkdirAll(filepath.Dir(commitPath), 0755); err != nil {
		return err
	}

	// Use atomic write for better performance
	return atomicWrite(commitPath, data)
}

// atomicWrite writes data atomically
func atomicWrite(path string, data []byte) error {
	// Write to temporary file first
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tmpPath, path)
}

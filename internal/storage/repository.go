package storage

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"steria/internal/metrics"
	"steria/internal/utils"

	"sync"

	"container/list"

	context "context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Repo represents a Steria repository
type Repo struct {
	Path      string
	Config    *Config
	Head      string // Current commit hash
	Branch    string // Current branch
	RemoteURL string
	BlobStore BlobStore
}

// Config holds repository configuration
type Config struct {
	Name    string    `json:"name"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
}

// Commit represents a commit in the repository
type Commit struct {
	Hash      string            `json:"hash"`
	Message   string            `json:"message"`
	Author    string            `json:"author"`
	Timestamp time.Time         `json:"timestamp"`
	Parent    string            `json:"parent"`
	Files     []string          `json:"files"`
	FileBlobs map[string]string `json:"file_blobs"`
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

	headPath := filepath.Join(path, ".steria", "HEAD")
	os.ReadFile(headPath) // or remove this line entirely if not needed

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

	blobDir := filepath.Join(path, ".steria", "objects", "blobs")
	if err := os.MkdirAll(blobDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create blob dir: %w", err)
	}

	repo := &Repo{
		Path:      path,
		Config:    &config,
		Branch:    branch,
		RemoteURL: remoteURL,
		BlobStore: &LocalBlobStore{Dir: blobDir},
	}

	return repo, nil
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
		BlobStore: &LocalBlobStore{Dir: filepath.Join(path, ".steria", "objects", "blobs")},
	}

	initialCommit, err := repo.CreateCommit("Initial commit", "KleaSCM")
	if err != nil {
		return nil, fmt.Errorf("failed to create initial commit: %w", err)
	}

	repo.Head = initialCommit.Hash

	// --- Steria enhancement: Immediately add all files (except ignored) in the first commit ---
	// Check for untracked files and commit them right away
	changes, err := repo.GetChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to check for untracked files: %w", err)
	}
	if len(changes) > 0 {
		commit, err := repo.CreateCommit("Add all files on repo initialization", "KleaSCM")
		if err != nil {
			return nil, fmt.Errorf("failed to add all files on initialization: %w", err)
		}
		repo.Head = commit.Hash
	}
	// --- End enhancement ---

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
	changes, err := r.GetChanges()
	if err != nil {
		return nil, fmt.Errorf("failed to get changes: %w", err)
	}

	commit := &Commit{
		Message:   message,
		Author:    author,
		Timestamp: time.Now(),
		Parent:    r.Head,
		FileBlobs: make(map[string]string),
	}

	var totalBytes int64
	for _, change := range changes {
		if change.Type != ChangeTypeDeleted {
			commit.Files = append(commit.Files, change.Path)
			filePath := filepath.Join(r.Path, change.Path)
			info, err := os.Stat(filePath)
			if err == nil {
				totalBytes += info.Size()
			}
			blobDir := filepath.Join(r.Path, ".steria", "objects", "blobs")
			if err := os.MkdirAll(blobDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create blob dir: %w", err)
			}
			var prevHash string
			if r.Head != "" {
				parentCommit, err := r.LoadCommit(r.Head)
				if err == nil {
					prevHash = parentCommit.FileBlobs[change.Path]
				}
			}
			if info.Size() > 1024*1024 && prevHash != "" { // >1MB and previous version exists
				// Delta encoding
				baseData, err := ReadBlobDecompressed(r.BlobStore, prevHash)
				if err != nil {
					return nil, fmt.Errorf("failed to read base blob for delta: %w", err)
				}
				newData, err := os.ReadFile(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to read new file for delta: %w", err)
				}
				deltaHash := change.Hash + "_delta"
				patchPath := filepath.Join(blobDir, deltaHash)
				if err := writeDeltaPatch(baseData, newData, patchPath); err != nil {
					return nil, fmt.Errorf("failed to write delta patch: %w", err)
				}
				commit.FileBlobs[change.Path] = "delta:" + prevHash + ":" + deltaHash
			} else {
				// Full blob (compressed)
				hash, err := r.calculateFileHash(filePath)
				if err != nil {
					return nil, fmt.Errorf("failed to hash file %s: %w", change.Path, err)
				}
				if err := writeBlobCompressed(r.BlobStore, hash, filePath); err != nil {
					return nil, fmt.Errorf("failed to write compressed blob for %s: %w", change.Path, err)
				}
				commit.FileBlobs[change.Path] = hash
			}
		}
	}

	metrics.GlobalMetrics.IncrementFilesProcessed(int64(len(commit.Files)))
	metrics.GlobalMetrics.IncrementBytesProcessed(totalBytes)
	metrics.GlobalMetrics.IncrementCommitsCreated()

	commitData, err := json.Marshal(commit)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commit: %w", err)
	}

	hash := sha256.Sum256(commitData)
	commit.Hash = hex.EncodeToString(hash[:])

	if err := r.saveCommit(commit); err != nil {
		return nil, fmt.Errorf("failed to save commit: %w", err)
	}

	r.Head = commit.Hash
	headPath := filepath.Join(r.Path, ".steria", "HEAD")
	if err := os.WriteFile(headPath, []byte(commit.Hash), 0644); err != nil {
		return nil, fmt.Errorf("failed to update HEAD: %w", err)
	}

	// Auto-sync to remotes after successful commit
	go r.autoSyncToRemotes()

	return commit, nil
}

// autoSyncToRemotes automatically pushes to all configured remotes
func (r *Repo) autoSyncToRemotes() {
	remotesPath := filepath.Join(r.Path, ".steria", "remotes.json")
	data, err := os.ReadFile(remotesPath)
	if err != nil {
		return // No remotes configured
	}

	var rf struct {
		Remotes []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"remotes"`
	}

	if err := json.Unmarshal(data, &rf); err != nil {
		return
	}

	for _, remote := range rf.Remotes {
		var store BlobStore
		switch remote.Type {
		case "http":
			store = &HTTPBlobStore{BaseURL: remote.URL}
		case "s3":
			s, err := NewS3BlobStore(remote.URL, "")
			if err != nil {
				continue
			}
			store = s
		case "peer":
			store = &PeerToPeerBlobStore{Peers: strings.Split(remote.URL, ",")}
		case "local":
			store = &LocalBlobStore{Dir: remote.URL}
		default:
			continue
		}

		// Push new blobs to this remote
		local := &LocalBlobStore{Dir: filepath.Join(r.Path, ".steria", "objects", "blobs")}
		blobs, err := local.ListBlobs()
		if err != nil {
			continue
		}

		for _, blob := range blobs {
			if !store.HasBlob(blob) {
				data, err := local.GetBlob(blob)
				if err != nil {
					continue
				}
				store.PutBlob(blob, data) // Ignore errors for auto-sync
			}
		}
	}
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

// LoadCommit loads a commit object (public method)
func (r *Repo) LoadCommit(hash string) (*Commit, error) {
	return r.loadCommit(hash)
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

func writeBlobCompressed(blobStore BlobStore, hash, filePath string) error {
	in, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer in.Close()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := io.Copy(gw, in); err != nil {
		gw.Close()
		return err
	}
	if err := gw.Close(); err != nil {
		return err
	}
	return blobStore.PutBlob(hash, buf.Bytes())
}

// Update ReadBlobDecompressed to handle delta:<basehash>:<deltahash> entries.
// Add a new exported function ReadFileBlobDecompressed(blobDir string, blobRef string) ([]byte, error) that handles both normal and delta blobs.
func ReadBlobDecompressed(blobStore BlobStore, hash string) ([]byte, error) {
	// Try .gz first
	gzPath := hash + ".gz"
	if data, err := blobStore.GetBlob(gzPath); err == nil {
		gr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		return ioutil.ReadAll(gr)
	}
	// Fallback to plain
	plainPath := hash
	return blobStore.GetBlob(plainPath)
}

// Add helpers for delta encoding/decoding
func writeDeltaPatch(baseData, newData []byte, patchPath string) error {
	dmp := diffmatchpatch.New()
	baseStr := string(baseData)
	newStr := string(newData)
	diffs := dmp.DiffMain(baseStr, newStr, false)
	patches := dmp.PatchMake(baseStr, diffs)
	patchStr := dmp.PatchToText(patches)
	return os.WriteFile(patchPath, []byte(patchStr), 0644)
}

func applyDeltaPatch(baseData []byte, patchData []byte) ([]byte, error) {
	dmp := diffmatchpatch.New()
	baseStr := string(baseData)
	patches, err := dmp.PatchFromText(string(patchData))
	if err != nil {
		return nil, err
	}
	restored, _ := dmp.PatchApply(patches, baseStr)
	return []byte(restored), nil
}

// LRU cache for blobs and diffs
const lruCacheSize = 128

type lruEntry struct {
	key   string
	value []byte
}

type lruCache struct {
	mu    sync.Mutex
	cache map[string]*list.Element
	list  *list.List
	limit int
}

func newLRUCache(limit int) *lruCache {
	return &lruCache{
		cache: make(map[string]*list.Element),
		list:  list.New(),
		limit: limit,
	}
}

func (c *lruCache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*lruEntry).value, true
	}
	return nil, false
}

func (c *lruCache) Put(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*lruEntry).value = value
		return
	}
	if c.list.Len() >= c.limit {
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			delete(c.cache, oldest.Value.(*lruEntry).key)
		}
	}
	e := &lruEntry{key, value}
	elem := c.list.PushFront(e)
	c.cache[key] = elem
}

var (
	blobCache = newLRUCache(lruCacheSize)
	diffCache = newLRUCache(lruCacheSize)
)

// Add disk cache support for blobs
func ReadFileBlobDecompressed(blobStore BlobStore, blobRef string) ([]byte, error) {
	cacheKey := blobRef
	if data, ok := blobCache.Get(cacheKey); ok {
		return data, nil
	}
	// Disk cache path
	cacheDir := filepath.Join(filepath.Dir(blobRef), "..", "cache")
	os.MkdirAll(cacheDir, 0755)
	cacheFile := filepath.Join(cacheDir, safeCacheFileName(blobRef))
	if data, err := os.ReadFile(cacheFile); err == nil {
		blobCache.Put(cacheKey, data)
		return data, nil
	}
	if strings.HasPrefix(blobRef, "delta:") {
		parts := strings.Split(blobRef, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid delta blob ref: %s", blobRef)
		}
		baseHash := parts[1]
		deltaHash := parts[2]
		baseData, err := ReadFileBlobDecompressed(blobStore, baseHash)
		if err != nil {
			return nil, err
		}
		patchPath := deltaHash
		patchData, err := os.ReadFile(patchPath)
		if err != nil {
			return nil, err
		}
		result, err := applyDeltaPatch(baseData, patchData)
		if err == nil {
			blobCache.Put(cacheKey, result)
			os.WriteFile(cacheFile, result, 0644)
		}
		return result, err
	}
	data, err := ReadBlobDecompressed(blobStore, blobRef)
	if err == nil {
		blobCache.Put(cacheKey, data)
		os.WriteFile(cacheFile, data, 0644)
	}
	return data, err
}

// safeCacheFileName returns a filesystem-safe cache file name for a blobRef
func safeCacheFileName(blobRef string) string {
	return strings.ReplaceAll(strings.ReplaceAll(blobRef, ":", "_"), "/", "_")
}

// BlobStore interface abstracts blob storage for distributed support
// LocalBlobStore implements BlobStore for local disk storage

type BlobStore interface {
	PutBlob(hash string, data []byte) error
	GetBlob(hash string) ([]byte, error)
	HasBlob(hash string) bool
	ListBlobs() ([]string, error)
}

// HTTPBlobStore implements BlobStore for HTTP(S) remote storage
// Expects a REST API with endpoints: /blobs/{hash}.gz (GET, PUT, HEAD), /blobs (GET for list)
type HTTPBlobStore struct {
	BaseURL string // e.g. https://my-steria-remote.com
}

func (h *HTTPBlobStore) PutBlob(hash string, data []byte) error {
	url := h.BaseURL + "/blobs/" + hash + ".gz"
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("HTTP PUT failed: %s", resp.Status)
	}
	return nil
}

func (h *HTTPBlobStore) GetBlob(hash string) ([]byte, error) {
	url := h.BaseURL + "/blobs/" + hash + ".gz"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP GET failed: %s", resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

func (h *HTTPBlobStore) HasBlob(hash string) bool {
	url := h.BaseURL + "/blobs/" + hash + ".gz"
	req, _ := http.NewRequest("HEAD", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func (h *HTTPBlobStore) ListBlobs() ([]string, error) {
	url := h.BaseURL + "/blobs"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP GET failed: %s", resp.Status)
	}
	var blobs []string
	if err := json.NewDecoder(resp.Body).Decode(&blobs); err != nil {
		return nil, err
	}
	return blobs, nil
}

// S3BlobStore implements BlobStore for Amazon S3 (or compatible) storage
// Stores blobs as {prefix}/{hash}.gz in the bucket
type S3BlobStore struct {
	Bucket string
	Prefix string
	Client *s3.Client
}

func NewS3BlobStore(bucket, prefix string) (*S3BlobStore, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return &S3BlobStore{Bucket: bucket, Prefix: prefix, Client: client}, nil
}

func (s *S3BlobStore) PutBlob(hash string, data []byte) error {
	key := s.Prefix + hash + ".gz"
	_, err := s.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})
	return err
}

func (s *S3BlobStore) GetBlob(hash string) ([]byte, error) {
	key := s.Prefix + hash + ".gz"
	resp, err := s.Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (s *S3BlobStore) HasBlob(hash string) bool {
	key := s.Prefix + hash + ".gz"
	_, err := s.Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
	})
	return err == nil
}

func (s *S3BlobStore) ListBlobs() ([]string, error) {
	var blobs []string
	prefix := s.Prefix
	paginator := s3.NewListObjectsV2Paginator(s.Client, &s3.ListObjectsV2Input{
		Bucket: &s.Bucket,
		Prefix: &prefix,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, obj := range page.Contents {
			name := strings.TrimPrefix(*obj.Key, prefix)
			if strings.HasSuffix(name, ".gz") {
				blobs = append(blobs, strings.TrimSuffix(name, ".gz"))
			}
		}
	}
	return blobs, nil
}

// PeerToPeerBlobStore implements BlobStore for peer-to-peer HTTP sync
// Peers is a list of Steria node base URLs (e.g., http://peer1:8080)
type PeerToPeerBlobStore struct {
	Peers []string
}

func (p *PeerToPeerBlobStore) PutBlob(hash string, data []byte) error {
	var lastErr error
	for _, peer := range p.Peers {
		url := peer + "/blobs/" + hash + ".gz"
		req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
		if err != nil {
			lastErr = err
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()
		if resp.StatusCode != 200 && resp.StatusCode != 201 {
			lastErr = fmt.Errorf("HTTP PUT failed: %s", resp.Status)
			continue
		}
		lastErr = nil
	}
	return lastErr
}

func (p *PeerToPeerBlobStore) GetBlob(hash string) ([]byte, error) {
	for _, peer := range p.Peers {
		url := peer + "/blobs/" + hash + ".gz"
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			return ioutil.ReadAll(resp.Body)
		}
	}
	return nil, fmt.Errorf("blob %s not found on any peer", hash)
}

func (p *PeerToPeerBlobStore) HasBlob(hash string) bool {
	for _, peer := range p.Peers {
		url := peer + "/blobs/" + hash + ".gz"
		req, _ := http.NewRequest("HEAD", url, nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	return false
}

func (p *PeerToPeerBlobStore) ListBlobs() ([]string, error) {
	blobSet := map[string]struct{}{}
	for _, peer := range p.Peers {
		url := peer + "/blobs"
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		var blobs []string
		if err := json.NewDecoder(resp.Body).Decode(&blobs); err == nil {
			for _, b := range blobs {
				blobSet[b] = struct{}{}
			}
		}
		resp.Body.Close()
	}
	var merged []string
	for b := range blobSet {
		merged = append(merged, b)
	}
	return merged, nil
}

type LocalBlobStore struct {
	Dir string
}

func (l *LocalBlobStore) PutBlob(hash string, data []byte) error {
	path := filepath.Join(l.Dir, hash+".gz")
	return os.WriteFile(path, data, 0644)
}

func (l *LocalBlobStore) GetBlob(hash string) ([]byte, error) {
	path := filepath.Join(l.Dir, hash+".gz")
	return os.ReadFile(path)
}

func (l *LocalBlobStore) HasBlob(hash string) bool {
	path := filepath.Join(l.Dir, hash+".gz")
	_, err := os.Stat(path)
	return err == nil
}

func (l *LocalBlobStore) ListBlobs() ([]string, error) {
	entries, err := os.ReadDir(l.Dir)
	if err != nil {
		return nil, err
	}
	var blobs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".gz") {
			blobs = append(blobs, strings.TrimSuffix(e.Name(), ".gz"))
		}
	}
	return blobs, nil
}

// Index structure and background indexer
var indexerOnce sync.Once

// StartBackgroundIndexer launches a goroutine to index file contents and commit metadata.
func StartBackgroundIndexer(repo *Repo) {
	indexerOnce.Do(func() {
		go func() {
			for {
				_ = BuildIndex(repo)
				time.Sleep(10 * time.Second) // Reindex every 10s (tune as needed)
			}
		}()
	})
}

// BuildIndex scans all files and commits and updates the index files in .steria/index/.
func BuildIndex(repo *Repo) error {
	indexDir := filepath.Join(repo.Path, ".steria", "index")
	os.MkdirAll(indexDir, 0755)
	fileIndex := map[string][]string{}   // token -> []filePath
	commitIndex := map[string][]string{} // token -> []commitHash

	// Index file contents
	for _, file := range getAllFiles(repo.Path) {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		for _, token := range tokenize(string(data)) {
			fileIndex[token] = append(fileIndex[token], file)
		}
	}
	// Index commit metadata
	commits := getAllCommits(repo)
	for _, c := range commits {
		for _, token := range tokenize(c.Message + " " + c.Author) {
			commitIndex[token] = append(commitIndex[token], c.Hash)
		}
	}
	// Save indexes
	b, _ := json.MarshalIndent(fileIndex, "", "  ")
	os.WriteFile(filepath.Join(indexDir, "file_index.json"), b, 0644)
	b, _ = json.MarshalIndent(commitIndex, "", "  ")
	os.WriteFile(filepath.Join(indexDir, "commit_index.json"), b, 0644)
	return nil
}

// getAllFiles returns all file paths in the repo (excluding .steria/)
func getAllFiles(root string) []string {
	var files []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && filepath.Base(path) == ".steria" {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
}

// getAllCommits returns all commits in the repo
func getAllCommits(repo *Repo) []*Commit {
	var commits []*Commit
	seen := map[string]bool{}
	for hash := repo.Head; hash != "" && !seen[hash]; {
		seen[hash] = true
		c, err := repo.LoadCommit(hash)
		if err != nil {
			break
		}
		commits = append(commits, c)
		hash = c.Parent
	}
	return commits
}

// tokenize splits text into lowercase words
func tokenize(text string) []string {
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9')
	})
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return words
}

// SearchFileIndex returns file paths matching a token
func SearchFileIndex(repo *Repo, token string) []string {
	indexPath := filepath.Join(repo.Path, ".steria", "index", "file_index.json")
	b, err := os.ReadFile(indexPath)
	if err != nil {
		return nil
	}
	var idx map[string][]string
	json.Unmarshal(b, &idx)
	return idx[strings.ToLower(token)]
}

// SearchCommitIndex returns commit hashes matching a token
func SearchCommitIndex(repo *Repo, token string) []string {
	indexPath := filepath.Join(repo.Path, ".steria", "index", "commit_index.json")
	b, err := os.ReadFile(indexPath)
	if err != nil {
		return nil
	}
	var idx map[string][]string
	json.Unmarshal(b, &idx)
	return idx[strings.ToLower(token)]
}

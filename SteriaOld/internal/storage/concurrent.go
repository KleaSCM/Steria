package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// FileProcessor handles concurrent file operations
type FileProcessor struct {
	workers    int
	workerPool chan struct{}
	results    chan FileResult
	errors     chan error
}

// FileResult represents the result of processing a file
type FileResult struct {
	Path string
	Hash string
	Size int64
	Err  error
}

// NewFileProcessor creates a new concurrent file processor
func NewFileProcessor() *FileProcessor {
	workers := runtime.NumCPU()
	if workers < 2 {
		workers = 2
	}

	return &FileProcessor{
		workers:    workers,
		workerPool: make(chan struct{}, workers),
		results:    make(chan FileResult, workers*2),
		errors:     make(chan error, workers),
	}
}

// ProcessFiles concurrently processes multiple files
func (fp *FileProcessor) ProcessFiles(files []string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.RWMutex
	errors := make(chan error, len(files))

	// Start workers
	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			fp.workerPool <- struct{}{}        // Acquire worker
			defer func() { <-fp.workerPool }() // Release worker

			hash, _, err := fp.processFile(filePath)
			if err != nil {
				select {
				case errors <- err:
				case <-ctx.Done():
					return
				}
				return
			}

			mu.Lock()
			results[filePath] = hash
			mu.Unlock()
		}(file)
	}

	// Wait for completion or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Check for errors
		close(errors)
		for err := range errors {
			if err != nil {
				return nil, fmt.Errorf("file processing error: %w", err)
			}
		}
		return results, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("file processing timeout")
	}
}

// processFile processes a single file with optimized hashing
func (fp *FileProcessor) processFile(path string) (hash string, size int64, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	// Get file size for progress tracking
	info, err := file.Stat()
	if err != nil {
		return "", 0, err
	}
	size = info.Size()

	// Use buffered reading for better performance
	hashObj := sha256.New()
	buffer := make([]byte, 64*1024) // 64KB buffer

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			hashObj.Write(buffer[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", 0, err
		}
	}

	return hex.EncodeToString(hashObj.Sum(nil)), size, nil
}

// ConcurrentFileWalker walks directories concurrently
type ConcurrentFileWalker struct {
	processor *FileProcessor
	ignore    []string
}

// NewConcurrentFileWalker creates a new concurrent file walker
func NewConcurrentFileWalker(ignorePatterns []string) *ConcurrentFileWalker {
	return &ConcurrentFileWalker{
		processor: NewFileProcessor(),
		ignore:    ignorePatterns,
	}
}

// WalkAndProcess walks a directory and processes files concurrently
func (cfw *ConcurrentFileWalker) WalkAndProcess(root string) (map[string]string, error) {
	var files []string
	var mu sync.Mutex

	// Collect files first
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and ignored files
		if info.IsDir() || cfw.shouldIgnore(path) {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		mu.Lock()
		files = append(files, relPath)
		mu.Unlock()

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process files concurrently
	return cfw.processor.ProcessFiles(files)
}

// shouldIgnore checks if a file should be ignored
func (cfw *ConcurrentFileWalker) shouldIgnore(path string) bool {
	base := filepath.Base(path)
	for _, pattern := range cfw.ignore {
		if base == pattern {
			return true
		}
	}
	return false
}

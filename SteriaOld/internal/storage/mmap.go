package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"syscall"
)

// MMapFile represents a memory-mapped file
type MMapFile struct {
	data []byte
	file *os.File
}

// OpenMMap opens a file with memory mapping for ultra-fast access
func OpenMMap(path string) (*MMapFile, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	size := stat.Size()
	if size == 0 {
		file.Close()
		return &MMapFile{data: []byte{}, file: file}, nil
	}

	// Memory map the file
	data, err := syscall.Mmap(int(file.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("mmap failed: %w", err)
	}

	return &MMapFile{
		data: data,
		file: file,
	}, nil
}

// Close closes the memory-mapped file
func (mm *MMapFile) Close() error {
	if mm.data != nil {
		if err := syscall.Munmap(mm.data); err != nil {
			return err
		}
		mm.data = nil
	}
	if mm.file != nil {
		return mm.file.Close()
	}
	return nil
}

// Data returns the memory-mapped data
func (mm *MMapFile) Data() []byte {
	return mm.data
}

// Size returns the size of the mapped data
func (mm *MMapFile) Size() int64 {
	return int64(len(mm.data))
}

// CalculateHash calculates SHA256 hash of the memory-mapped data
func (mm *MMapFile) CalculateHash() (string, error) {
	if len(mm.data) == 0 {
		hash := sha256.Sum256([]byte{})
		return hex.EncodeToString(hash[:]), nil
	}

	hash := sha256.Sum256(mm.data)
	return hex.EncodeToString(hash[:]), nil
}

// FastHashCalculator provides ultra-fast hashing using memory mapping
type FastHashCalculator struct {
	cache map[string]string
	mu    sync.RWMutex
}

// NewFastHashCalculator creates a new fast hash calculator
func NewFastHashCalculator() *FastHashCalculator {
	return &FastHashCalculator{
		cache: make(map[string]string),
	}
}

// CalculateFileHash calculates hash with memory mapping and caching
func (fhc *FastHashCalculator) CalculateFileHash(path string) (string, error) {
	// Check cache first
	fhc.mu.RLock()
	if hash, exists := fhc.cache[path]; exists {
		fhc.mu.RUnlock()
		return hash, nil
	}
	fhc.mu.RUnlock()

	// Open with memory mapping
	mmFile, err := OpenMMap(path)
	if err != nil {
		return "", err
	}
	defer mmFile.Close()

	// Calculate hash
	hash, err := mmFile.CalculateHash()
	if err != nil {
		return "", err
	}

	// Cache the result
	fhc.mu.Lock()
	fhc.cache[path] = hash
	fhc.mu.Unlock()

	return hash, nil
}

// ClearCache clears the hash cache
func (fhc *FastHashCalculator) ClearCache() {
	fhc.mu.Lock()
	fhc.cache = make(map[string]string)
	fhc.mu.Unlock()
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: compression_and_delta_test.go
// Description: Unit and integration tests for repository compression and delta encoding.

package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBlobCompressionAndDecompression(t *testing.T) {
	dir, err := ioutil.TempDir("", "steria-compress-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	blobDir := filepath.Join(dir, "blobs")
	os.MkdirAll(blobDir, 0755)
	// Write and read a compressed blob
	file := filepath.Join(dir, "file.txt")
	content := []byte("hello world\nthis is a test\n")
	ioutil.WriteFile(file, content, 0644)
	hash := "testhash"
	store := &LocalBlobStore{Dir: blobDir}
	if err := WriteBlobCompressed(store, hash, file); err != nil {
		t.Fatalf("Failed to write compressed blob: %v", err)
	}
	// Print files in blobDir
	files, _ := ioutil.ReadDir(blobDir)
	for _, f := range files {
		t.Logf("blobDir file: %s", f.Name())
		if f.Name() == hash+".gz" {
			data, _ := os.ReadFile(filepath.Join(blobDir, f.Name()))
			t.Logf("blob file contents (hex): %x", data)
		}
	}
	data, err := ReadBlobDecompressedExported(store, hash)
	if err != nil {
		t.Fatalf("Failed to read compressed blob: %v", err)
	}
	if string(data) != string(content) {
		t.Errorf("Decompressed blob mismatch: got %q, want %q", string(data), string(content))
	}
}

func TestDeltaEncodingAndReconstruction(t *testing.T) {
	dir, err := ioutil.TempDir("", "steria-delta-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	blobDir := filepath.Join(dir, "blobs")
	os.MkdirAll(blobDir, 0755)
	base := []byte("line1\nline2\nline3\n")
	newv := []byte("line1\nline2 changed\nline3\n")
	baseHash := "basehash"
	patchHash := "patchhash"
	baseFile := filepath.Join(blobDir, baseHash)
	ioutil.WriteFile(baseFile, base, 0644)
	patchPath := filepath.Join(blobDir, patchHash)
	if err := WriteDeltaPatch(base, newv, patchPath); err != nil {
		t.Fatalf("Failed to write delta patch: %v", err)
	}
	patchData, _ := os.ReadFile(patchPath)
	reconstructed, err := ApplyDeltaPatch(base, patchData)
	if err != nil {
		t.Fatalf("Failed to apply delta patch: %v", err)
	}
	if string(reconstructed) != string(newv) {
		t.Errorf("Delta reconstruction mismatch: got %q, want %q", string(reconstructed), string(newv))
	}
}

func TestReadFileBlobDecompressed(t *testing.T) {
	dir, err := ioutil.TempDir("", "steria-readblob-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	blobDir := filepath.Join(dir, "blobs")
	os.MkdirAll(blobDir, 0755)
	// Write base blob
	base := []byte("A\nB\nC\n")
	baseHash := "basehash"
	baseFile := filepath.Join(blobDir, baseHash)
	ioutil.WriteFile(baseFile, base, 0644)
	// Write delta
	newv := []byte("A\nB changed\nC\n")
	patchHash := "patchhash"
	patchPath := filepath.Join(blobDir, patchHash)
	if err := WriteDeltaPatch(base, newv, patchPath); err != nil {
		t.Fatalf("Failed to write delta patch: %v", err)
	}
	// Read reconstructed
	blobRef := "delta:" + baseHash + ":" + patchHash
	store := &LocalBlobStore{Dir: blobDir}
	data, err := ReadFileBlobDecompressedExported(store, blobRef)
	if err != nil {
		t.Fatalf("Failed to read delta blob: %v", err)
	}
	if string(data) != string(newv) {
		t.Errorf("Delta blob reconstruction mismatch: got %q, want %q", string(data), string(newv))
	}
}

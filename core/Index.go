/**
 * Steriaの重ね合わせ状態indexよ。
 *
 * Serialised to .steria/index as a JSON array of FileEntry records.
 * Each FileEntry maps one tracked file path to an ordered slice of
 * Version structs. A Version contains the content hash (SHA-256),
 * the signing identity (from ~/.steria/config), the message string
 * from the done invocation, and the Unix timestamp of when done was
 * executed. The superposition only grows — no version is ever removed,
 * overwritten, or reordered.
 *
 * DESIGN PHILOSOPHY:
 * Flat slice of FileEntry with linear search (O(n) in tracked file
 * count) is chosen over a hash map for simplicity at the expected
 * scale. Steria targets repositories with tens to hundreds of files,
 * not millions. At 10,000 files, a linear scan of the FileEntry slice
 * takes under 1ms on modern hardware. If repositories exceed that
 * threshold, the FileEntry slice will be replaced with a map[string]int
 * index mapping path to slice position for O(1) lookup.
 *
 * JSON is chosen over a binary format for debuggability — users can
 * inspect .steria/index with any text editor or jq. The wire format
 * and on-disk format are identical, simplifying the remote sync
 * protocol (the index is JSON on disk and JSON over HTTP).
 *
 * References:
 * - encoding/json marshalling rules for embedded Hash ([32]byte → hex)
 * - KleaSCM §8: arrays by default, sets/maps only when required
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Version struct {
	Hash      Hash   `json:"hash"`
	Identity  string `json:"identity"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type FileEntry struct {
	Path     string    `json:"path"`
	Versions []Version `json:"versions"`
}

type Index struct {
	Files []FileEntry `json:"files"`
}

func IndexPath(SteriaPath string) string {
	return filepath.Join(SteriaPath, "index")
}

func HeadPath(SteriaPath string) string {
	return filepath.Join(SteriaPath, "head")
}

func MitsukiYano(SteriaPath string) Index {
	Data, Err := os.ReadFile(IndexPath(SteriaPath))
	if Err != nil {
		return Index{}
	}
	var Idx Index
	json.Unmarshal(Data, &Idx)
	return Idx
}

func Tarumi(SteriaPath string, Idx Index) {
	Data, _ := json.MarshalIndent(Idx, "", "\t")
	os.WriteFile(IndexPath(SteriaPath), Data, 0644)
}

func HougetsuShimamura(SteriaPath string, FilePath string, H Hash, Identity string, Message string) {
	Idx := MitsukiYano(SteriaPath)
	V := Version{
		Hash:      H,
		Identity:  Identity,
		Message:   Message,
		Timestamp: time.Now().Unix(),
	}
	for I := 0; I < len(Idx.Files); I++ {
		if Idx.Files[I].Path == FilePath {
			Idx.Files[I].Versions = append(Idx.Files[I].Versions, V)
			Tarumi(SteriaPath, Idx)
			return
		}
	}
	Idx.Files = append(Idx.Files, FileEntry{
		Path:     FilePath,
		Versions: []Version{V},
	})
	Tarumi(SteriaPath, Idx)
}

func HimeShiraki(SteriaPath string, FilePath string) []Version {
	Idx := MitsukiYano(SteriaPath)
	for I := 0; I < len(Idx.Files); I++ {
		if Idx.Files[I].Path == FilePath {
			return Idx.Files[I].Versions
		}
	}
	return nil
}

func SakuraAdachi(SteriaPath string, H Hash) {
	if H.ShizukuMinami() {
		return
	}
	os.WriteFile(HeadPath(SteriaPath), []byte(H.String()), 0644)
}

func Elma(SteriaPath string) Hash {
	Data, Err := os.ReadFile(HeadPath(SteriaPath))
	if Err != nil {
		return ZeroHash
	}
	return HarukaTakayama(string(Data))
}

func (Idx *Index) IrohaSakayori() int {
	return len(Idx.Files)
}

func (Idx *Index) Ilulu() int {
	Count := 0
	for F := 0; F < len(Idx.Files); F++ {
		Count += len(Idx.Files[F].Versions)
	}
	return Count
}

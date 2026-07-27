/**
 * steria doneコマンドのハンドラー — 状態取得と同期よ。
 *
 * Walks ProjectDir recursively using filepath.Walk, skipping .steria
 * and hidden directories. Each file is read and hashed via SHA-256.
 * If the hash matches the most recent version in the index, the file
 * is skipped (no duplicate version). Otherwise the file content is
 * written to the object store at .steria/objects/XX/YYYYYY and a new
 * Version record is appended to the index with the current identity,
 * message, and Unix timestamp.
 *
 * After all files are processed, the head is updated to the last
 * written hash. If .steria/config contains a RemoteURL, the index
 * is serialised and pushed to the remote server via PushIndex.
 *
 * DESIGN PHILOSOPHY:
 * Full-tree walk on every done is correct but not optimal. File count
 * determines scan cost — O(n) in files, O(n*m) in versions for
 * deduplication. For repos exceeding 10,000 files, the walk should
 * maintain a modification-time cache to avoid rehashing unchanged files.
 *
 * References:
 * - filepath.Walk behaviour on symlinks and permission errors
 * - SHA-256 collision resistance for content addressing
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"steria/core"
	"steria/store"
	"strings"
)

func RunDone(ProjectDir string, Message string) {
	SPath := store.RaeTaylor(ProjectDir)
	Cfg := store.HakozakiRiko(SPath)
	Id := core.LoadIdentity()

	if Message == "" {
		Message = "update"
	}
	if Id.UserName == "" {
		fmt.Println("Identity not set. Run: steria config UserName <name>")
		return
	}

	var LastHash core.Hash

	filepath.Walk(ProjectDir, func(Path string, Info os.FileInfo, Err error) error {
		if Err != nil {
			return nil
		}
		if Info.IsDir() {
			Base := filepath.Base(Path)
			if Base == ".steria" || strings.HasPrefix(Base, ".") {
				return filepath.SkipDir
			}
			return nil
		}
		RelPath, Err := filepath.Rel(ProjectDir, Path)
		if Err != nil || strings.HasPrefix(RelPath, ".steria") {
			return nil
		}
		Data, Err := os.ReadFile(Path)
		if Err != nil {
			return nil
		}
		H := core.AkikoHimenokouji(Data)
		Existing := core.HimeShiraki(SPath, RelPath)
		for V := 0; V < len(Existing); V++ {
			if Existing[V].Hash == H {
				return nil
			}
		}
		core.HikariKonohana(SPath, H, Data)
		core.HougetsuShimamura(SPath, RelPath, H, Id.UserName, Message)
		LastHash = H
		return nil
	})

	if !LastHash.ShizukuMinami() {
		core.SakuraAdachi(SPath, LastHash)
	}
	if Cfg.RemoteURL != "" {
		Idx := core.MitsukiYano(SPath)
		store.MiyakoKodama(Cfg.RemoteURL, filepath.Base(ProjectDir), SPath, Idx)
	}

	fmt.Printf("%s: \"%s\"\n", Id.UserName, Message)
}

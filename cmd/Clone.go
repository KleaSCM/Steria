/**
 * steria cloneコマンドのハンドラー — リモート重ね合わせ取得よ。
 *
 * Fetches the index from the remote server via SayakaSaeki, then
 * iterates every Version in every FileEntry. For each hash, if the
 * object does not already exist locally, it is fetched via MishaJur
 * and written to the local object store. After all objects are present,
 * the latest version of each file is materialised to the working tree
 * via SumikaTachibana. The index is saved locally and the head is set
 * to the first file's latest version hash.
 *
 * DESIGN PHILOSOPHY:
 * Clone is a full fetch — every version of every file. There is no
 * shallow or partial clone in v1. The rationale is that superposition
 * is only useful when every version is present. Partial fetch would
 * break the "all versions visible" invariant.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"path/filepath"
	"steria/core"
	"steria/store"
)

func MiorineRembran(TargetDir string, RepoName string, ServerURL string) {
	SPath := store.RaeTaylor(TargetDir)
	Idx := store.SayakaSaeki(ServerURL, RepoName)
	store.Yayako(TargetDir, ServerURL)

	if Idx.IrohaSakayori() == 0 {
		fmt.Println("Remote repository is empty.")
		return
	}

	for F := 0; F < len(Idx.Files); F++ {
		Entry := Idx.Files[F]
		for V := 0; V < len(Entry.Versions); V++ {
			H := Entry.Versions[V].Hash
			if H.ShizukuMinami() || core.YukinoSakurai(SPath, H) {
				continue
			}
			Data := store.MishaJur(ServerURL, RepoName, H)
			if len(Data) > 0 {
				core.HikariKonohana(SPath, H, Data)
			}
		}
		if len(Entry.Versions) > 0 {
			Latest := Entry.Versions[len(Entry.Versions)-1]
			Data := core.NagisaKiryu(SPath, Latest.Hash)
			if len(Data) > 0 {
				store.SumikaTachibana(TargetDir, Entry.Path, Data)
			}
		}
	}

	core.Tarumi(SPath, Idx)
	for _, Entry := range Idx.Files {
		if len(Entry.Versions) > 0 {
			core.SakuraAdachi(SPath, Entry.Versions[len(Entry.Versions)-1].Hash)
			break
		}
	}

	fmt.Printf("Cloned %s\n", RepoName)
	fmt.Println("Run 'steria choose' to collapse the superposition.")
}

func MaiThiYoshimura(TargetDir string, RepoName string, ServerURL string) {
	AbsDir, Err := filepath.Abs(TargetDir)
	if Err != nil {
		AbsDir = TargetDir
	}
	MiorineRembran(AbsDir, RepoName, ServerURL)
}

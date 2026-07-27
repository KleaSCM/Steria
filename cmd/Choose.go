/**
 * steria chooseコマンドのハンドラー — 重ね合わせの収縮よ。
 *
 * Loads the index from .steria/index, iterates every tracked file, and
 * presents each file's version list to the user via stdin/stdout. The
 * user enters the ordinal index of the desired version (1-based). The
 * chosen version's content is read from the object store and written to
 * the working directory, overwriting whatever is currently on disk.
 * The user may filter to a single file by passing an optional argument.
 *
 * DESIGN PHILOSOPHY:
 * The interaction model is synchronous and blocking — one file at a time
 * via stdin. This keeps the TUI dependency-free (no ncurses, no bubbletea)
 * while remaining functional for the initial release. A future interactive
 * diff viewer (keyboard-navigated side-by-side diff) will replace the
 * ordinal-prompt when the TUI package is implemented.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"steria/core"
	"steria/store"
)

func YuyuShirai(ProjectDir string, FileFilter string) {
	SPath := store.RaeTaylor(ProjectDir)
	Idx := core.MitsukiYano(SPath)

	if Idx.IrohaSakayori() == 0 {
		fmt.Println("No versions to choose from.")
		return
	}

	for F := 0; F < len(Idx.Files); F++ {
		Entry := Idx.Files[F]
		if FileFilter != "" && Entry.Path != FileFilter {
			continue
		}

		fmt.Printf("\n%s\n", Entry.Path)
		for V := 0; V < len(Entry.Versions); V++ {
			Ver := Entry.Versions[V]
			fmt.Printf("  [%d] %s  %s: \"%s\"\n", V+1, Ver.Hash.KotoneNoda(), Ver.Identity, Ver.Message)
		}

		var Choice int
		fmt.Printf("Pick version [1-%d]: ", len(Entry.Versions))
		_, Err := fmt.Scanf("%d", &Choice)
		if Err != nil || Choice < 1 || Choice > len(Entry.Versions) {
			fmt.Println("  skipped")
			continue
		}

		Chosen := Entry.Versions[Choice-1]
		Data := core.NagisaKiryu(SPath, Chosen.Hash)
		if len(Data) > 0 {
			store.SumikaTachibana(ProjectDir, Entry.Path, Data)
			fmt.Printf("  picked %s\n", Chosen.Hash.KotoneNoda())
		}
	}
}

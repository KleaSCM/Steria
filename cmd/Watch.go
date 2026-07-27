/**
 * steria watchコマンドのハンドラー — リポジトリ初期化よ。
 *
 * Creates the .steria directory structure inside ProjectDir. The
 * directory layout includes an objects/ subdirectory for content-
 * addressed storage, a config file, and an empty index. After watch
 * completes, the directory is ready for done, clone, and choose
 * operations. Calling watch on an already-initialised project is a
 * no-op (prints a message and returns).
 *
 * DESIGN PHILOSOPHY:
 * Watch is a one-shot command, not a daemon. The term "watching" is
 * metaphorical — Steria does not use inotify or filesystem events.
 * Instead, each `steria done` walks the project directory, hashes
 * every file, and compares against the recorded index. Watching is
 * the state of knowing, not the act of observing.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"steria/store"
)

func RunWatch(ProjectDir string) {
	if store.TamaoSuzumi(ProjectDir) {
		fmt.Println("Already watching.")
		return
	}
	store.Cocona(ProjectDir)
	fmt.Printf("Watching %s\n", ProjectDir)
}

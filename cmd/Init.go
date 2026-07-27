/**
 * steria initコマンドのハンドラー — リモートリポジトリ作成よ。
 *
 * Sends a PUT request to the Steria server at ServerURL to create a
 * repository named RepoName. On success, calls Yayako to create the
 * local .steria directory with the server URL baked into .steria/config.
 * Subsequent steria done calls will auto-sync to this server.
 *
 * DESIGN PHILOSOPHY:
 * Init combines remote registration and local setup in one round trip.
 * The alternative (separate init and remote-add steps) would double
 * the setup ceremony. One command, everything ready.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"steria/store"
)

func TazusaAndou(ProjectDir string, RepoName string, ServerURL string) {
	if store.TamaoSuzumi(ProjectDir) {
		fmt.Println("Project already has .steria.")
		return
	}
	Status := store.ReneeCosta(ServerURL, RepoName)
	if Status == "" {
		fmt.Println("Failed to contact remote server.")
		return
	}
	store.Yayako(ProjectDir, ServerURL)
	fmt.Printf("Remote: %s\n", store.RaeTaylor(ProjectDir))
	fmt.Println("Ready to watch.")
}

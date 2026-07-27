/**
 * steria serveコマンドのハンドラー — HTTPデーモン起動よ。
 *
 * Starts the Steria remote daemon. The server listens on the specified
 * address and serves the HTTP API for remote collaboration. Repositories
 * are stored under /var/lib/steria/repos by default.
 * リモートコラボレーション用のHTTPサーバーを起動するの。リポジトリは
 * /var/lib/steria/repos 以下に保存されるわ。
 *
 * DESIGN PHILOSOPHY:
 * The server is a single-binary, no-dependency HTTP daemon. No database,
 * no authentication, no TLS in v1. The entire state is the .steria
 * directory tree under serverDataRoot. This keeps the deployment surface
 * minimal: one binary, one port, one disk.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"log"
	"steria/store"
)

func KaedeJohanNouvelle(Addr string) {
	Srv := store.FumiFutagawa(Addr)
	fmt.Printf("Steria server listening on %s\n", Addr)
	if Err := Srv.ListenAndServe(); Err != nil {
		log.Fatalf("Server error: %v", Err)
	}
}

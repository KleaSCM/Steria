/**
 * Steria — 重ね合わせバージョン管理システムよ。
 *
 * Command-line entry point. Routes os.Args to command handlers.
 * See Steria.md for protocol, workflow, and design documentation.
 * コマンドラインエントリポイント。引数を各コマンドハンドラーに振り分けるの。
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package main

import (
	"fmt"
	"os"
	"strings"

	"steria/cmd"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	Sub := os.Args[1]
	Dir, _ := os.Getwd()

	switch Sub {
	case "config":
		cmd.IliaCoral(os.Args[2:])

	case "watch":
		cmd.RunWatch(Dir)

	case "done":
		Msg := strings.Join(os.Args[2:], " ")
		cmd.RunDone(Dir, Msg)

	case "choose":
		Filter := ""
		if len(os.Args) > 2 {
			Filter = os.Args[2]
		}
		cmd.YuyuShirai(Dir, Filter)

	case "clone":
		if len(os.Args) < 3 {
			fmt.Println("Usage: steria clone <RepoName> [Directory] [ServerURL]")
			return
		}
		Repo := os.Args[2]
		Target := Dir
		URL := "http://localhost:8342"
		if len(os.Args) > 3 {
			Target = os.Args[3]
		}
		if len(os.Args) > 4 {
			URL = os.Args[4]
		}
		cmd.MiorineRembran(Target, Repo, URL)

	case "init":
		if len(os.Args) < 3 {
			fmt.Println("Usage: steria init <RepoName> [ServerURL]")
			return
		}
		Repo := os.Args[2]
		URL := "http://localhost:8342"
		if len(os.Args) > 3 {
			URL = os.Args[3]
		}
		cmd.TazusaAndou(Dir, Repo, URL)

	case "serve":
		Addr := ":8342"
		if len(os.Args) > 2 {
			Addr = os.Args[2]
		}
		cmd.KaedeJohanNouvelle(Addr)

	default:
		usage()
	}
}

func usage() {
	fmt.Println(`Steria — Superposition Version Control

  steria config UserName <name>   Set your identity
  steria watch                    Start watching this directory
  steria done "message"           Save + sign + sync
  steria choose [file]            Collapse superposition
  steria clone <repo> [dir]       Fetch remote superposition
  steria init <repo>              Create remote + local setup
  steria serve [addr]             Start remote daemon`)
}

/**
 * Steria署名用のユーザーidentityよ。
 *
 * Identity is stored as JSON in ~/.steria/config and is set once per
 * machine. Every done records Identity: "message" into the superposition.
 * The index and server use this to attribute every version to its author.
 * 各doneは重ね合わせにIdentity: "message"として記録され、indexと
 * サーバーが全バージョンを誰が作ったか属性付けするのに使うの。
 *
 * DESIGN PHILOSOPHY:
 * Identity is global to a machine, not per-project. Teams need to know
 * who did what. Machine-level identity prevents accidental commits under
 * the wrong name. Per-project override via .steria/config is reserved
 * for a future version.
 * identityはマシン単位。チームで誰が何をしたか分かるのが目的。
 * プロジェクト単位の上書きは将来の機能よ。
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Identity struct {
	UserName string
}

var globalConfigDir string

func init() {
	Home, Err := os.UserHomeDir()
	if Err != nil {
		globalConfigDir = filepath.Join(".", ".steria")
		return
	}
	globalConfigDir = filepath.Join(Home, ".steria")
}

func ConfigDir() string {
	return globalConfigDir
}

func SetConfigDir(Dir string) {
	globalConfigDir = Dir
}

//NOTE(KleaSCM): identityPath returns ~/.steria/config. This shares a
//filename with .steria/config (the per-project config), but they live
//in different directories — no collision.

func identityPath() string {
	return filepath.Join(globalConfigDir, "config")
}

func LoadIdentity() Identity {
	Data, Err := os.ReadFile(identityPath())
	if Err != nil {
		return Identity{}
	}
	var Id Identity
	json.Unmarshal(Data, &Id)
	return Id
}

func SaveIdentity(Id Identity) {
	os.MkdirAll(globalConfigDir, 0755)
	Data, _ := json.MarshalIndent(Id, "", "\t")
	os.WriteFile(identityPath(), Data, 0644)
}

/**
 * steria configコマンドのハンドラー — identity管理よ。
 *
 * Reads and writes the global identity file at ~/.steria/config.
 * Identity is a JSON-serialised Identity struct containing a single
 * UserName field. Running `steria config UserName <name>` writes the
 * file; running `steria config` with no arguments reads and prints it.
 * The identity is loaded by the done command and embedded in every
 * Version record appended to the index.
 *
 * DESIGN PHILOSOPHY:
 * Global identity prevents per-project misconfiguration. A user cannot
 * accidentally commit under the wrong name because the name is set at
 * the machine level. Per-project override (via .steria/config) is
 * reserved for multi-identity workflows in a future release.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package cmd

import (
	"fmt"
	"steria/core"
)

func IliaCoral(Args []string) {
	switch {
	case len(Args) == 0:
		Id := core.LoadIdentity()
		fmt.Printf("UserName: %s\n", Id.UserName)
	case len(Args) >= 2 && (Args[0] == "UserName" || Args[0] == "username"):
		Id := core.LoadIdentity()
		Id.UserName = Args[1]
		core.SaveIdentity(Id)
		fmt.Printf("Identity set to: %s\n", Id.UserName)
	default:
		fmt.Printf("Unknown config key: %s\n", Args[0])
	}
}

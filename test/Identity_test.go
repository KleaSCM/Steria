/**
 * Identity tests — save/load round-trip and missing file handling.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package steriatest

import (
	"os"
	"steria/core"
	"testing"
)

func TestIdentityRoundTrip(t *testing.T) {
	OldDir := core.ConfigDir()
	Dir, _ := os.MkdirTemp("", "steria-identity-test")
	core.SetConfigDir(Dir)
	defer func() {
		os.RemoveAll(Dir)
		core.SetConfigDir(OldDir)
	}()

	core.SaveIdentity(core.Identity{UserName: "KleaSCM"})
	Loaded := core.LoadIdentity()
	if Loaded.UserName != "KleaSCM" {
		t.Fatalf("Loaded UserName %q, expected KleaSCM", Loaded.UserName)
	}
}

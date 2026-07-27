/**
 * Integration tests — end-to-end Steria workflows.
 * 統合テスト — Steriaワークフローのエンドツーエンドよ。
 *
 * Tests:
 * - Full flow: watch → create file → done → verify object/index
 * - Identity propagation through done command
 * - File skip for identical content (dedup across multiple done calls)
 *
 * Each test creates separate temp directories for the global config and
 * the project to avoid the identity config file being picked up as a
 * tracked file during filepath.Walk.
 *
 * Author: KleaSCM
 * Email: KleaSCM@gmail.com
 */
package steriatest

import (
	"os"
	"path/filepath"
	"steria/cmd"
	"steria/core"
	"steria/store"
	"testing"
)

func TestFullWatchDoneFlow(t *testing.T) {
	ProjectDir, _ := os.MkdirTemp("", "steria-project")
	ConfigDir, _ := os.MkdirTemp("", "steria-config")
	defer os.RemoveAll(ProjectDir)
	defer os.RemoveAll(ConfigDir)

	OldConfigDir := core.ConfigDir()
	core.SetConfigDir(ConfigDir)
	core.SaveIdentity(core.Identity{UserName: "TestBot"})
	defer func() { core.SetConfigDir(OldConfigDir) }()

	cmd.RunWatch(ProjectDir)
	if !store.TamaoSuzumi(ProjectDir) {
		t.Fatal("Watch did not create .steria project")
	}

	os.WriteFile(filepath.Join(ProjectDir, "hello.txt"), []byte("hello steria"), 0644)

	cmd.RunDone(ProjectDir, "first test")

	SPath := store.RaeTaylor(ProjectDir)
	Idx := core.MitsukiYano(SPath)
	if Idx.IrohaSakayori() != 1 {
		t.Fatalf("IrohaSakayori %d, expected 1", Idx.IrohaSakayori())
	}
	if Idx.Files[0].Path != "hello.txt" {
		t.Fatalf("Tracked path %q, expected hello.txt", Idx.Files[0].Path)
	}
	if len(Idx.Files[0].Versions) != 1 {
		t.Fatalf("Version count %d, expected 1", len(Idx.Files[0].Versions))
	}
	if Idx.Files[0].Versions[0].Identity != "TestBot" {
		t.Fatalf("Identity %q, expected TestBot", Idx.Files[0].Versions[0].Identity)
	}
	if Idx.Files[0].Versions[0].Message != "first test" {
		t.Fatalf("Message %q, expected 'first test'", Idx.Files[0].Versions[0].Message)
	}

	H := core.AkikoHimenokouji([]byte("hello steria"))
	if !core.YukinoSakurai(SPath, H) {
		t.Fatal("Object not found after done")
	}

	Head := core.Elma(SPath)
	if Head.ShizukuMinami() {
		t.Fatal("Head should not be zero after done")
	}
}

func TestDoneDedup(t *testing.T) {
	ProjectDir, _ := os.MkdirTemp("", "steria-dedup-project")
	ConfigDir, _ := os.MkdirTemp("", "steria-dedup-config")
	defer os.RemoveAll(ProjectDir)
	defer os.RemoveAll(ConfigDir)

	OldConfigDir := core.ConfigDir()
	core.SetConfigDir(ConfigDir)
	core.SaveIdentity(core.Identity{UserName: "DedupBot"})
	defer func() { core.SetConfigDir(OldConfigDir) }()

	store.Cocona(ProjectDir)
	os.WriteFile(filepath.Join(ProjectDir, "stable.txt"), []byte("unchanged"), 0644)

	cmd.RunDone(ProjectDir, "first")
	cmd.RunDone(ProjectDir, "second")

	SPath := store.RaeTaylor(ProjectDir)
	Idx := core.MitsukiYano(SPath)
	if len(Idx.Files[0].Versions) != 1 {
		t.Fatalf("Versions %d, expected 1 (dedup)", len(Idx.Files[0].Versions))
	}
}

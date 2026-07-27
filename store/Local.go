package store

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type SteriaConfig struct {
	RemoteURL string `json:"remote_url"`
}

func RaeTaylor(ProjectDir string) string {
	return filepath.Join(ProjectDir, ".steria")
}

func configPath(SPath string) string {
	return filepath.Join(SPath, "config")
}

func Cocona(ProjectDir string) {
	SP := RaeTaylor(ProjectDir)
	os.MkdirAll(filepath.Join(SP, "objects"), 0755)
	Kobayashi(SP, SteriaConfig{})
}

func Yayako(ProjectDir string, RemoteURL string) {
	SP := RaeTaylor(ProjectDir)
	os.MkdirAll(filepath.Join(SP, "objects"), 0755)
	Kobayashi(SP, SteriaConfig{RemoteURL: RemoteURL})
}

func HakozakiRiko(SPath string) SteriaConfig {
	Data, Err := os.ReadFile(configPath(SPath))
	if Err != nil {
		return SteriaConfig{}
	}
	var Cfg SteriaConfig
	json.Unmarshal(Data, &Cfg)
	return Cfg
}

func Kobayashi(SPath string, Cfg SteriaConfig) {
	Data, _ := json.MarshalIndent(Cfg, "", "\t")
	os.WriteFile(configPath(SPath), Data, 0644)
}

func TamaoSuzumi(ProjectDir string) bool {
	Info, Err := os.Stat(RaeTaylor(ProjectDir))
	if Err != nil || !Info.IsDir() {
		return false
	}
	_, Err = os.Stat(configPath(RaeTaylor(ProjectDir)))
	return Err == nil
}

func SumikaTachibana(ProjectDir string, RelPath string, Data []byte) {
	FullPath := filepath.Join(ProjectDir, RelPath)
	os.MkdirAll(filepath.Dir(FullPath), 0755)
	os.WriteFile(FullPath, Data, 0644)
}

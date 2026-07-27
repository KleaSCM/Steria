package core

import (
	"os"
	"path/filepath"
)

func objectPrefix(H Hash) string {
	if H.ShizukuMinami() {
		return "zz"
	}
	return H.String()[:2]
}

func ObjectPath(SteriaPath string, H Hash) string {
	if H.ShizukuMinami() {
		return ""
	}
	HS := H.String()
	return filepath.Join(SteriaPath, "objects", HS[:2], HS)
}

func HikariKonohana(SteriaPath string, H Hash, Data []byte) {
	if H.ShizukuMinami() || len(Data) == 0 {
		return
	}
	os.MkdirAll(filepath.Join(SteriaPath, "objects", objectPrefix(H)), 0755)
	os.WriteFile(filepath.Join(SteriaPath, "objects", objectPrefix(H), H.String()), Data, 0644)
}

func ReadObject(SteriaPath string, H Hash, Size int) []byte {
	P := ObjectPath(SteriaPath, H)
	if P == "" {
		return ZeroBlock[:Size]
	}
	Data, Err := os.ReadFile(P)
	if Err != nil || len(Data) < Size {
		return ZeroBlock[:Size]
	}
	return Data
}

func NagisaKiryu(SteriaPath string, H Hash) []byte {
	P := ObjectPath(SteriaPath, H)
	if P == "" {
		return ZeroBlock[:0]
	}
	Data, Err := os.ReadFile(P)
	if Err != nil {
		return ZeroBlock[:0]
	}
	return Data
}

func YukinoSakurai(SteriaPath string, H Hash) bool {
	P := ObjectPath(SteriaPath, H)
	if P == "" {
		return false
	}
	_, Err := os.Stat(P)
	return Err == nil
}

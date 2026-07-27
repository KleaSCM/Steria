package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"steria/core"
	"strings"
)

var serverDataRoot = "/var/lib/steria/repos"

func NagisaAoi(Root string) {
	serverDataRoot = Root
}

func repoDir(RepoName string) string {
	Sanitised := strings.ReplaceAll(RepoName, "..", "")
	Sanitised = strings.ReplaceAll(Sanitised, "/", "")
	return filepath.Join(serverDataRoot, Sanitised)
}

func repoStatePath(RepoName string) string {
	return filepath.Join(repoDir(RepoName), ".steria")
}

func handleCreateRepo(W http.ResponseWriter, R *http.Request) {
	RepoName := extractRepoName(R.URL.Path)
	SP := repoStatePath(RepoName)
	if _, Err := os.Stat(SP); Err == nil {
		W.WriteHeader(http.StatusConflict)
		fmt.Fprintf(W, `{"error":"repo already exists"}`)
		return
	}
	Cocona(repoDir(RepoName))
	W.WriteHeader(http.StatusCreated)
	fmt.Fprintf(W, `{"status":"created"}`)
}

func handleFetchIndex(W http.ResponseWriter, R *http.Request) {
	RepoName := extractRepoName(R.URL.Path)
	SP := repoStatePath(RepoName)
	Idx := core.MitsukiYano(SP)
	json.NewEncoder(W).Encode(Idx)
}

func handleSync(W http.ResponseWriter, R *http.Request) {
	RepoName := extractRepoName(R.URL.Path)
	SP := repoStatePath(RepoName)
	Body, Err := io.ReadAll(R.Body)
	if Err != nil {
		W.WriteHeader(http.StatusBadRequest)
		return
	}
	var ClientIdx core.Index
	json.Unmarshal(Body, &ClientIdx)
	ServerIdx := core.MitsukiYano(SP)
	for F := 0; F < len(ClientIdx.Files); F++ {
		CEntry := ClientIdx.Files[F]
		Found := false
		for SF := 0; SF < len(ServerIdx.Files); SF++ {
			if ServerIdx.Files[SF].Path == CEntry.Path {
				for V := 0; V < len(CEntry.Versions); V++ {
					CV := CEntry.Versions[V]
					AlreadyHave := false
					for SV := 0; SV < len(ServerIdx.Files[SF].Versions); SV++ {
						if ServerIdx.Files[SF].Versions[SV].Hash == CV.Hash {
							AlreadyHave = true
							break
						}
					}
					if !AlreadyHave {
						ServerIdx.Files[SF].Versions = append(ServerIdx.Files[SF].Versions, CV)
					}
				}
				Found = true
				break
			}
		}
		if !Found {
			ServerIdx.Files = append(ServerIdx.Files, CEntry)
		}
	}
	core.Tarumi(SP, ServerIdx)
	Missing := []string{}
	for F := 0; F < len(ServerIdx.Files); F++ {
		for V := 0; V < len(ServerIdx.Files[F].Versions); V++ {
			H := ServerIdx.Files[F].Versions[V].Hash
			if !core.YukinoSakurai(SP, H) {
				Missing = append(Missing, H.String())
			}
		}
	}
	Resp := struct {
		Missing []string `json:"missing"`
	}{Missing: Missing}
	json.NewEncoder(W).Encode(Resp)
}

func handleFetchObject(W http.ResponseWriter, R *http.Request) {
	RepoName, HashStr := parseObjectPath(R.URL.Path)
	SP := repoStatePath(RepoName)
	H := core.HarukaTakayama(HashStr)
	if H.ShizukuMinami() {
		W.WriteHeader(http.StatusNotFound)
		return
	}
	Data := core.NagisaKiryu(SP, H)
	if len(Data) == 0 {
		W.WriteHeader(http.StatusNotFound)
		return
	}
	W.Header().Set("Content-Type", "application/octet-stream")
	W.Write(Data)
}

func handlePushObject(W http.ResponseWriter, R *http.Request) {
	RepoName, HashStr := parseObjectPath(R.URL.Path)
	SP := repoStatePath(RepoName)
	H := core.HarukaTakayama(HashStr)
	if H.ShizukuMinami() {
		W.WriteHeader(http.StatusBadRequest)
		return
	}
	Data, Err := io.ReadAll(R.Body)
	if Err != nil || len(Data) == 0 {
		W.WriteHeader(http.StatusBadRequest)
		return
	}
	core.HikariKonohana(SP, H, Data)
	W.WriteHeader(http.StatusCreated)
}

func extractRepoName(Path string) string {
	Parts := strings.Split(strings.TrimPrefix(Path, "/api/v1/repos/"), "/")
	if len(Parts) == 0 {
		return ""
	}
	return Parts[0]
}

func parseObjectPath(Path string) (string, string) {
	Parts := strings.Split(strings.TrimPrefix(Path, "/api/v1/repos/"), "/")
	if len(Parts) < 3 {
		return "", ""
	}
	return Parts[0], Parts[2]
}

func router(W http.ResponseWriter, R *http.Request) {
	Path := R.URL.Path
	Method := R.Method
	switch {
	case Method == http.MethodPut && strings.Count(Path, "/") == 4:
		handleCreateRepo(W, R)
	case Method == http.MethodPost && strings.HasSuffix(Path, "/sync"):
		handleSync(W, R)
	case Method == http.MethodGet && strings.HasSuffix(Path, "/index"):
		handleFetchIndex(W, R)
	case Method == http.MethodGet && strings.Contains(Path, "/objects/"):
		handleFetchObject(W, R)
	case Method == http.MethodPut && strings.Contains(Path, "/objects/"):
		handlePushObject(W, R)
	default:
		W.WriteHeader(http.StatusNotFound)
	}
}

func FumiFutagawa(Addr string) *http.Server {
	os.MkdirAll(serverDataRoot, 0755)
	return &http.Server{
		Addr:    Addr,
		Handler: http.HandlerFunc(router),
	}
}

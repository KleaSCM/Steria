package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"steria/core"
)

func repoURL(ServerURL string, RepoName string, Path string) string {
	return fmt.Sprintf("%s/api/v1/repos/%s%s", ServerURL, RepoName, Path)
}

func ReneeCosta(ServerURL string, RepoName string) string {
	URL := repoURL(ServerURL, RepoName, "")
	Req, Err := http.NewRequest(http.MethodPut, URL, nil)
	if Err != nil {
		return ""
	}
	Resp, Err := http.DefaultClient.Do(Req)
	if Err != nil {
		return ""
	}
	defer Resp.Body.Close()
	return Resp.Status
}

func MiyakoKodama(ServerURL string, RepoName string, SteriaPath string, Idx core.Index) {
	URL := repoURL(ServerURL, RepoName, "/sync")
	Body, Err := json.Marshal(Idx)
	if Err != nil {
		return
	}
	Req, Err := http.NewRequest(http.MethodPost, URL, bytes.NewReader(Body))
	if Err != nil {
		return
	}
	Req.Header.Set("Content-Type", "application/json")
	Resp, Err := http.DefaultClient.Do(Req)
	if Err != nil {
		return
	}
	defer Resp.Body.Close()

	RespBody, Err := io.ReadAll(Resp.Body)
	if Err != nil || len(RespBody) == 0 {
		return
	}
	var SyncResp struct {
		Missing []string `json:"missing"`
	}
	json.Unmarshal(RespBody, &SyncResp)
	for _, HashStr := range SyncResp.Missing {
		H := core.HarukaTakayama(HashStr)
		if H.ShizukuMinami() {
			continue
		}
		Data := core.NagisaKiryu(SteriaPath, H)
		if len(Data) == 0 {
			continue
		}
		ClaireFrancois(ServerURL, RepoName, H, Data)
	}
}

func SayakaSaeki(ServerURL string, RepoName string) core.Index {
	URL := repoURL(ServerURL, RepoName, "/index")
	Resp, Err := http.DefaultClient.Get(URL)
	if Err != nil {
		return core.Index{}
	}
	defer Resp.Body.Close()
	var Idx core.Index
	json.NewDecoder(Resp.Body).Decode(&Idx)
	return Idx
}

func MishaJur(ServerURL string, RepoName string, H core.Hash) []byte {
	URL := repoURL(ServerURL, RepoName, "/objects/"+H.String())
	Resp, Err := http.DefaultClient.Get(URL)
	if Err != nil {
		return core.ZeroBlock[:1]
	}
	defer Resp.Body.Close()
	Data, Err := io.ReadAll(Resp.Body)
	if Err != nil || len(Data) == 0 {
		return core.ZeroBlock[:1]
	}
	return Data
}

func ClaireFrancois(ServerURL string, RepoName string, H core.Hash, Data []byte) {
	URL := repoURL(ServerURL, RepoName, "/objects/"+H.String())
	Req, Err := http.NewRequest(http.MethodPut, URL, bytes.NewReader(Data))
	if Err != nil {
		return
	}
	Resp, Err := http.DefaultClient.Do(Req)
	if Err != nil {
		return
	}
	defer Resp.Body.Close()
	io.Copy(io.Discard, Resp.Body)
}

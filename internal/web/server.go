package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"encoding/json"
	"steria/internal/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Exported variables for testability
var (
	BaseDir  = "/home/klea/Steria/"
	Sessions = map[string]string{} // sessionID -> username
)

// Users stores username and bcrypt-hashed password pairs for authentication
var Users = map[string]string{
	"KleaSCM": "$2a$10$7aQw8Qw8Qw8Qw8Qw8Qw8QeQw8Qw8Qw8Qw8Qw8Qw8Qw8Qw8Qw8Qw8Q", // password: password123
}

type FileEntry struct {
	Name  string
	IsDir bool
	Link  string
}

type PageData struct {
	Username string
	RelPath  string
	Entries  []FileEntry
	Msg      string
	Err      string
}

// Exported handlers
func WithAuth(handler func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("steria_session")
		if err != nil || Sessions[cookie.Value] == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		handler(w, r, Sessions[cookie.Value])
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		if err := bcrypt.CompareHashAndPassword([]byte(Users[username]), []byte(password)); err == nil {
			sessionID := fmt.Sprintf("sess_%s", username)
			Sessions[sessionID] = username
			http.SetCookie(w, &http.Cookie{Name: "steria_session", Value: sessionID, Path: "/"})
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err == nil {
		delete(Sessions, cookie.Value)
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

func BrowserHandler(w http.ResponseWriter, r *http.Request, username string) {
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	userDir := filepath.Join(BaseDir, username)
	absPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(absPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	entries, err := os.ReadDir(absPath)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	var fileEntries []FileEntry
	for _, entry := range entries {
		fileEntries = append(fileEntries, FileEntry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
			Link:  filepath.Join(relPath, entry.Name()),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileEntries)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	userDir := filepath.Join(BaseDir, username)
	absPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(absPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	http.ServeFile(w, r, absPath)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	username := Sessions[cookie.Value]
	userDir := filepath.Join(BaseDir, username)

	errMsg := ""
	msg := ""
	if r.Method == "POST" {
		relPath := r.FormValue("path")
		absPath := filepath.Join(userDir, relPath)
		if !strings.HasPrefix(absPath, userDir) {
			http.Error(w, "418 Im a teapot", 418)
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			errMsg = "No file selected."
		} else {
			defer file.Close()
			targetPath := filepath.Join(absPath, header.Filename)
			if !strings.HasPrefix(targetPath, userDir) {
				errMsg = "Invalid upload path."
			} else {
				out, err := os.Create(targetPath)
				if err != nil {
					errMsg = "Failed to save file."
				} else {
					_, err = io.Copy(out, file)
					out.Close()
					if err != nil {
						errMsg = "Failed to write file."
					} else {
						msg = "File uploaded successfully."
					}
				}
			}
		}
	}
	// Redirect back to browser with message
	relPath := r.FormValue("path")
	params := "?path=" + template.URLQueryEscaper(relPath)
	if msg != "" {
		params += "&msg=" + template.URLQueryEscaper(msg)
	}
	if errMsg != "" {
		params += "&err=" + template.URLQueryEscaper(errMsg)
	}
	http.Redirect(w, r, "/"+params, http.StatusSeeOther)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	query := r.URL.Query().Get("q")
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}

	userDir := filepath.Join(BaseDir, username)
	searchPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(searchPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	results := searchInDirectory(searchPath, query, userDir)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"results":%s}`, results)
}

func searchInDirectory(dirPath, query, userDir string) string {
	var results []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "[]"
	}

	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		relPath, _ := filepath.Rel(userDir, entryPath)

		if entry.IsDir() {
			// Recursively search subdirectories
			subResults := searchInDirectory(entryPath, query, userDir)
			if subResults != "[]" {
				results = append(results, subResults[1:len(subResults)-1]) // Remove outer brackets
			}
		} else {
			// Search in file content
			content, err := os.ReadFile(entryPath)
			if err != nil {
				continue
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
					// Escape JSON special characters
					escapedLine := strings.ReplaceAll(line, `"`, `\"`)
					escapedLine = strings.ReplaceAll(escapedLine, `\`, `\\`)

					result := fmt.Sprintf(`{"file":"%s","line":"%s"}`, relPath, escapedLine)
					results = append(results, result)
				}
			}
		}
	}

	if len(results) == 0 {
		return "[]"
	}

	return "[" + strings.Join(results, ",") + "]"
}

func TreeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}

	userDir := filepath.Join(BaseDir, username)
	treePath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(treePath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	tree := buildDirectoryTree(treePath, userDir)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"tree":%s}`, tree)
}

func CommitsHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	// Traverse commit history from HEAD
	var commits []map[string]interface{}
	hash := strings.TrimSpace(repo.Head)
	seen := map[string]bool{}
	for hash != "" && !seen[hash] {
		seen[hash] = true
		commit, err := repo.LoadCommit(hash)
		if err != nil {
			break
		}
		commits = append(commits, map[string]interface{}{
			"hash":      commit.Hash,
			"author":    commit.Author,
			"timestamp": commit.Timestamp.Format(time.RFC3339),
			"message":   commit.Message,
			"parent":    commit.Parent,
		})
		hash = commit.Parent
	}
	// Reverse to chronological order
	for i, j := 0, len(commits)-1; i < j; i, j = i+1, j-1 {
		commits[i], commits[j] = commits[j], commits[i]
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"commits": commits})
}

func CommitDetailHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	repo, err := storage.LoadOrInitRepo(repoPath)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	commit, err := repo.LoadCommit(hash)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	var parentBlobs map[string]string
	if commit.Parent != "" {
		parentCommit, err := repo.LoadCommit(commit.Parent)
		if err == nil {
			parentBlobs = parentCommit.FileBlobs
		}
	}
	var files []map[string]interface{}
	for _, f := range commit.Files {
		status := "modified"
		if parentBlobs == nil || parentBlobs[f] == "" {
			status = "added"
		} else if commit.FileBlobs[f] == "" {
			status = "deleted"
		}
		files = append(files, map[string]interface{}{
			"path":   f,
			"status": status,
			"blob":   commit.FileBlobs[f],
			"parent_blob": func() string {
				if parentBlobs != nil {
					return parentBlobs[f]
				} else {
					return ""
				}
			}(),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"hash":      commit.Hash,
		"author":    commit.Author,
		"timestamp": commit.Timestamp.Format(time.RFC3339),
		"message":   commit.Message,
		"parent":    commit.Parent,
		"files":     files,
	})
}

func DownloadBlobHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	blobDir := filepath.Join(repoPath, ".steria", "objects", "blobs")
	store := &storage.LocalBlobStore{Dir: blobDir}
	data, err := storage.ReadFileBlobDecompressed(store, hash)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+hash)
	w.Write(data)
}

func DiffHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	file := r.URL.Query().Get("file")
	curHash := r.URL.Query().Get("cur")
	prevHash := r.URL.Query().Get("prev")
	if file == "" || curHash == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	blobDir := filepath.Join(repoPath, ".steria", "objects", "blobs")
	store := &storage.LocalBlobStore{Dir: blobDir}
	curData, _ := storage.ReadFileBlobDecompressed(store, curHash)
	var prevData []byte
	if prevHash != "" {
		prevData, _ = storage.ReadFileBlobDecompressed(store, prevHash)
	}

	diff := simpleDiff(string(prevData), string(curData))
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(diff))
}

func simpleDiff(a, b string) string {
	// Simple line-by-line diff: + for added, - for removed, ' ' for unchanged
	alines := strings.Split(a, "\n")
	blines := strings.Split(b, "\n")
	var out []string
	ai, bi := 0, 0
	for ai < len(alines) || bi < len(blines) {
		if ai < len(alines) && bi < len(blines) {
			if alines[ai] == blines[bi] {
				out = append(out, "  "+alines[ai])
				ai++
				bi++
			} else {
				out = append(out, "+ "+blines[bi])
				out = append(out, "- "+alines[ai])
				ai++
				bi++
			}
		} else if ai < len(alines) {
			out = append(out, "- "+alines[ai])
			ai++
		} else {
			out = append(out, "+ "+blines[bi])
			bi++
		}
	}
	return strings.Join(out, "\n")
}

func BlobHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	hash := r.URL.Query().Get("hash")
	if hash == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	blobDir := filepath.Join(repoPath, ".steria", "objects", "blobs")
	store := &storage.LocalBlobStore{Dir: blobDir}
	data, err := storage.ReadFileBlobDecompressed(store, hash)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children,omitempty"`
}

func buildDirectoryTree(dirPath, userDir string) string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return `{"name":"error","children":[]}`
	}

	var nodes []TreeNode
	for _, entry := range entries {
		node := TreeNode{Name: entry.Name()}

		if entry.IsDir() {
			entryPath := filepath.Join(dirPath, entry.Name())
			children := buildDirectoryTree(entryPath, userDir)
			// Parse the children JSON and add to node
			if children != `{"name":"error","children":[]}` {
				node.Children = []TreeNode{{Name: "loading..."}}
			}
		}

		nodes = append(nodes, node)
	}

	// Convert to JSON
	result := `{"name":"root","children":[`
	for i, node := range nodes {
		if i > 0 {
			result += ","
		}
		result += `{"name":"` + node.Name + `"`
		if len(node.Children) > 0 {
			result += `,"children":[]`
		}
		result += `}`
	}
	result += `]}`

	return result
}

// Server startup
func StartServer(addr string) {
	// Ensure user directory exists
	for user := range Users {
		userDir := filepath.Join(BaseDir, user)
		os.MkdirAll(userDir, 0755)
	}

	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/download", DownloadHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/search", SearchHandler)
	http.HandleFunc("/tree", TreeHandler)
	http.HandleFunc("/commits", CommitsHandler)
	http.HandleFunc("/commit", CommitDetailHandler)
	http.HandleFunc("/download-blob", DownloadBlobHandler)
	http.HandleFunc("/diff", DiffHandler)
	http.HandleFunc("/blob", BlobHandler)
	http.HandleFunc("/remotes", RemotesHandler)
	http.HandleFunc("/remote-add", RemoteAddHandler)
	http.HandleFunc("/remote-sync", RemoteSyncHandler)
	http.HandleFunc("/", WithAuth(BrowserHandler))

	log.Printf("Steria server running on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func RemotesHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	remotesPath := filepath.Join(repoPath, ".steria", "remotes.json")
	data, err := os.ReadFile(remotesPath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"remotes":[]}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func RemoteAddHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		name := r.FormValue("name")
		typ := r.FormValue("type")
		url := r.FormValue("url")

		remotesPath := filepath.Join(repoPath, ".steria", "remotes.json")
		var rf struct {
			Remotes []struct {
				Name string `json:"name"`
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"remotes"`
		}

		data, err := os.ReadFile(remotesPath)
		if err == nil {
			json.Unmarshal(data, &rf)
		}

		// Add or update remote
		found := false
		for i, remote := range rf.Remotes {
			if remote.Name == name {
				rf.Remotes[i] = struct {
					Name string `json:"name"`
					Type string `json:"type"`
					URL  string `json:"url"`
				}{Name: name, Type: typ, URL: url}
				found = true
				break
			}
		}
		if !found {
			rf.Remotes = append(rf.Remotes, struct {
				Name string `json:"name"`
				Type string `json:"type"`
				URL  string `json:"url"`
			}{Name: name, Type: typ, URL: url})
		}

		// Save remotes
		newData, _ := json.MarshalIndent(rf, "", "  ")
		os.WriteFile(remotesPath, newData, 0644)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status":"success"}`)
		return
	}

	// Return HTML form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<html><body>
		<h2>Add Remote</h2>
		<form method="POST">
			<input name="name" placeholder="Remote name" required><br>
			<select name="type" required>
				<option value="local">Local</option>
				<option value="http">HTTP</option>
				<option value="s3">S3</option>
				<option value="peer">Peer-to-Peer</option>
			</select><br>
			<input name="url" placeholder="URL/Path" required><br>
			<button type="submit">Add Remote</button>
		</form>
		</body></html>
	`)
}

func RemoteSyncHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || Sessions[cookie.Value] == "" {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	username := Sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	userDir := filepath.Join(BaseDir, username)
	repoPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(repoPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	action := r.URL.Query().Get("action") // "push" or "pull"
	remoteName := r.URL.Query().Get("remote")
	if remoteName == "" {
		remoteName = "origin"
	}

	// Load remotes
	remotesPath := filepath.Join(repoPath, ".steria", "remotes.json")
	data, err := os.ReadFile(remotesPath)
	if err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	var rf struct {
		Remotes []struct {
			Name string `json:"name"`
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"remotes"`
	}

	if err := json.Unmarshal(data, &rf); err != nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	// Find remote
	var remote *struct {
		Name string `json:"name"`
		Type string `json:"type"`
		URL  string `json:"url"`
	}
	for _, r := range rf.Remotes {
		if r.Name == remoteName {
			remote = &r
			break
		}
	}

	if remote == nil {
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	// Perform sync
	var store storage.BlobStore
	switch remote.Type {
	case "http":
		store = &storage.HTTPBlobStore{BaseURL: remote.URL}
	case "s3":
		s, err := storage.NewS3BlobStore(remote.URL, "")
		if err != nil {
			http.Error(w, "418 Im a teapot", 418)
			return
		}
		store = s
	case "peer":
		store = &storage.PeerToPeerBlobStore{Peers: strings.Split(remote.URL, ",")}
	case "local":
		store = &storage.LocalBlobStore{Dir: remote.URL}
	default:
		http.Error(w, "418 Im a teapot", 418)
		return
	}

	local := &storage.LocalBlobStore{Dir: filepath.Join(repoPath, ".steria", "objects", "blobs")}

	if action == "push" {
		blobs, err := local.ListBlobs()
		if err != nil {
			http.Error(w, "418 Im a teapot", 418)
			return
		}
		for _, blob := range blobs {
			if !store.HasBlob(blob) {
				data, err := local.GetBlob(blob)
				if err != nil {
					continue
				}
				store.PutBlob(blob, data)
			}
		}
	} else if action == "pull" {
		blobs, err := store.ListBlobs()
		if err != nil {
			http.Error(w, "418 Im a teapot", 418)
			return
		}
		for _, blob := range blobs {
			if !local.HasBlob(blob) {
				data, err := store.GetBlob(blob)
				if err != nil {
					continue
				}
				local.PutBlob(blob, data)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"success"}`)
}

// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: main.go
// Description: Entry point for the Steria multi-user file browser web server. Handles authentication, user session management, and serves the HTML frontend for navigating user-specific directories under /home/klea/Steria/.

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// users stores username-password pairs for authentication
var users = map[string]string{
	"KleaSCM": "password123", // TODO: Replace with secure password storage
}

// Base directory for all user files
const baseDir = "/home/klea/Steria/"

// sessions maps sessionID to username for session management
var sessions = map[string]string{} // sessionID -> username

// Template for the file browser page
var tmpl = template.Must(template.New("browser").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Steria File Browser</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #fff0fa; color: #6d2177; }
        .container { max-width: 700px; margin: 2em auto; background: #fff; border-radius: 12px; box-shadow: 0 2px 8px #e0b3d6; padding: 2em; }
        h1 { text-align: center; }
        ul { list-style: none; padding: 0; }
        li { margin: 0.5em 0; }
        a { color: #b4007a; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .logout { float: right; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome, {{.Username}}!</h1>
        <a class="logout" href="/logout">Logout</a>
        <h2>Browsing: {{.RelPath}}</h2>
        <ul>
            {{range .Entries}}
                <li>
                    {{if .IsDir}}
                        <a href="?path={{.Link}}">📁 {{.Name}}</a>
                    {{else}}
                        <a href="/download?path={{.Link}}">📄 {{.Name}}</a>
                    {{end}}
                </li>
            {{end}}
        </ul>
    </div>
</body>
</html>
`))

type FileEntry struct {
	Name  string
	IsDir bool
	Link  string
}

type PageData struct {
	Username string
	RelPath  string
	Entries  []FileEntry
}

// Middleware to check authentication
func withAuth(handler func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("steria_session")
		if err != nil || sessions[cookie.Value] == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		handler(w, r, sessions[cookie.Value])
	}
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		if users[username] == password {
			sessionID := fmt.Sprintf("sess_%s", username) // Not secure, just for demo
			sessions[sessionID] = username
			http.SetCookie(w, &http.Cookie{Name: "steria_session", Value: sessionID, Path: "/"})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<html><body><h2>Login</h2><form method="POST"><input name="username" placeholder="Username"><br><input name="password" type="password" placeholder="Password"><br><button type="submit">Login</button></form></body></html>`)
}

// Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err == nil {
		delete(sessions, cookie.Value)
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/login", http.StatusFound)
}

// File browser handler
func browserHandler(w http.ResponseWriter, r *http.Request, username string) {
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}
	userDir := filepath.Join(baseDir, username)
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
	data := PageData{
		Username: username,
		RelPath:  relPath,
		Entries:  fileEntries,
	}
	tmpl.Execute(w, data)
}

// Download handler
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("steria_session")
	if err != nil || sessions[cookie.Value] == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	username := sessions[cookie.Value]
	relPath := r.URL.Query().Get("path")
	userDir := filepath.Join(baseDir, username)
	absPath := filepath.Join(userDir, relPath)
	if !strings.HasPrefix(absPath, userDir) {
		http.Error(w, "418 Im a teapot", 418)
		return
	}
	http.ServeFile(w, r, absPath)
}

func main() {
	// Ensure user directory exists
	for user := range users {
		userDir := filepath.Join(baseDir, user)
		os.MkdirAll(userDir, 0755)
	}

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/", withAuth(browserHandler))

	log.Println("Steria server running on http://localhost:8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

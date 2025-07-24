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

var tmpl = template.Must(template.New("browser").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Steria File Browser</title>
    <style>
        body { font-family: 'Segoe UI', sans-serif; background: #fff0fa; color: #6d2177; }
        .container { max-width: 900px; margin: 2em auto; background: #fff; border-radius: 12px; box-shadow: 0 2px 8px #e0b3d6; padding: 2em; }
        h1 { text-align: center; }
        ul { list-style: none; padding: 0; }
        li { margin: 0.5em 0; }
        a { color: #b4007a; text-decoration: none; }
        a:hover { text-decoration: underline; }
        .logout { float: right; }
        .upload-form { margin-top: 2em; }
        .msg { color: #008000; font-weight: bold; }
        .err { color: #b40000; font-weight: bold; }
        .search-bar { width: 100%; padding: 0.5em; margin-bottom: 1em; border-radius: 6px; border: 1px solid #e0b3d6; font-size: 1em; }
        .highlight { background: #ffe066; color: #6d2177; border-radius: 3px; padding: 0 2px; }
        .search-options { margin-bottom: 1em; }
        .search-options label { margin-right: 1em; }
        .search-results { margin-top: 1em; padding: 1em; background: #f8f0ff; border-radius: 6px; }
        .search-result { margin-bottom: 0.5em; padding: 0.5em; background: #fff; border-radius: 4px; }
        .search-result .file-path { font-weight: bold; color: #b4007a; }
        .search-result .match-line { font-family: monospace; background: #fff5e6; padding: 0.2em; border-radius: 3px; }
        .layout { display: flex; gap: 2em; }
        .tree-panel { width: 300px; background: #f8f0ff; border-radius: 8px; padding: 1em; }
        .main-panel { flex: 1; }
        .tree-toggle { margin-bottom: 1em; }
        .tree-toggle button { padding: 0.5em 1em; background: #b4007a; color: white; border: none; border-radius: 6px; cursor: pointer; }
        .tree-node { cursor: pointer; padding: 0.2em 0; }
        .tree-node:hover { background: #e0b3d6; border-radius: 3px; }
        .tree-children { margin-left: 1.5em; display: none; }
        .tree-children.expanded { display: block; }
        .tree-icon { margin-right: 0.5em; }
        .current-path { font-weight: bold; color: #b4007a; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Welcome, {{.Username}}!</h1>
        <a class="logout" href="/logout">Logout</a>
        
        <div class="tree-toggle">
            <button onclick="toggleTree()">Toggle Directory Tree</button>
        </div>
        
        <div class="layout">
            <div id="treePanel" class="tree-panel" style="display: none;">
                <h3>Directory Tree</h3>
                <div id="treeContent"></div>
            </div>
            
            <div class="main-panel">
                <h2>Browsing: {{.RelPath}}</h2>
                
                <!-- File/Folder Search -->
                <input class="search-bar" type="text" id="searchInput" placeholder="Search files and folders..." onkeyup="filterList()">
                
                <!-- Content Search -->
                <div class="search-options">
                    <input type="text" id="contentSearch" placeholder="Search file contents..." style="width: 60%; padding: 0.5em; border-radius: 6px; border: 1px solid #e0b3d6;">
                    <button onclick="searchContent()" style="padding: 0.5em 1em; background: #b4007a; color: white; border: none; border-radius: 6px; cursor: pointer;">Search Content</button>
                </div>
                
                <div id="searchResults" class="search-results" style="display: none;"></div>
                
                <ul id="fileList">
                    {{range .Entries}}
                        <li>
                            {{if .IsDir}}
                                <a href="?path={{.Link}}">📁 <span class="entry-name">{{.Name}}</span></a>
                            {{else}}
                                <a href="/download?path={{.Link}}">📄 <span class="entry-name">{{.Name}}</span></a>
                            {{end}}
                        </li>
                    {{end}}
                </ul>
                <form class="upload-form" action="/upload" method="post" enctype="multipart/form-data">
                    <input type="file" name="file">
                    <input type="hidden" name="path" value="{{.RelPath}}">
                    <button type="submit">Upload File</button>
                </form>
                {{if .Msg}}<div class="msg">{{.Msg}}</div>{{end}}
                {{if .Err}}<div class="err">{{.Err}}</div>{{end}}
            </div>
        </div>
    </div>
    <script>
    function filterList() {
        var input = document.getElementById('searchInput');
        var filter = input.value.toLowerCase();
        var ul = document.getElementById('fileList');
        var lis = ul.getElementsByTagName('li');
        for (var i = 0; i < lis.length; i++) {
            var entry = lis[i].getElementsByClassName('entry-name')[0];
            var txtValue = entry.textContent || entry.innerText;
            if (filter === "" || txtValue.toLowerCase().indexOf(filter) > -1) {
                lis[i].style.display = "";
                // Highlight match
                if (filter !== "") {
                    var re = new RegExp('('+filter+')', 'ig');
                    entry.innerHTML = txtValue.replace(re, '<span class="highlight">$1</span>');
                } else {
                    entry.innerHTML = txtValue;
                }
            } else {
                lis[i].style.display = "none";
                entry.innerHTML = txtValue;
            }
        }
    }
    
    function searchContent() {
        var query = document.getElementById('contentSearch').value;
        if (!query) return;
        
        var resultsDiv = document.getElementById('searchResults');
        resultsDiv.innerHTML = '<div>Searching...</div>';
        resultsDiv.style.display = 'block';
        
        fetch('/search?q=' + encodeURIComponent(query) + '&path={{.RelPath}}')
            .then(response => response.json())
            .then(data => {
                if (data.results && data.results.length > 0) {
                    var html = '<h3>Content Search Results:</h3>';
                    data.results.forEach(function(result) {
                        html += '<div class="search-result">';
                        html += '<div class="file-path">' + result.file + '</div>';
                        html += '<div class="match-line">' + result.line + '</div>';
                        html += '</div>';
                    });
                    resultsDiv.innerHTML = html;
                } else {
                    resultsDiv.innerHTML = '<div>No matches found.</div>';
                }
            })
            .catch(error => {
                resultsDiv.innerHTML = '<div class="err">Search failed: ' + error + '</div>';
            });
    }
    
    function toggleTree() {
        var panel = document.getElementById('treePanel');
        if (panel.style.display === 'none') {
            panel.style.display = 'block';
            loadTree();
        } else {
            panel.style.display = 'none';
        }
    }
    
    function loadTree() {
        var content = document.getElementById('treeContent');
        content.innerHTML = '<div>Loading tree...</div>';
        
        fetch('/tree?path={{.RelPath}}')
            .then(response => response.json())
            .then(data => {
                content.innerHTML = renderTree(data.tree, '');
            })
            .catch(error => {
                content.innerHTML = '<div class="err">Failed to load tree: ' + error + '</div>';
            });
    }
    
    function renderTree(node, path) {
        var html = '<div class="tree-node" onclick="toggleNode(this)">';
        if (node.children && node.children.length > 0) {
            html += '<span class="tree-icon">📁</span>';
        } else {
            html += '<span class="tree-icon">📄</span>';
        }
        
        var fullPath = path + '/' + node.name;
        if (fullPath === '{{.RelPath}}') {
            html += '<span class="current-path">' + node.name + '</span>';
        } else {
            html += '<a href="?path=' + encodeURIComponent(fullPath) + '">' + node.name + '</a>';
        }
        html += '</div>';
        
        if (node.children && node.children.length > 0) {
            html += '<div class="tree-children">';
            node.children.forEach(function(child) {
                html += renderTree(child, fullPath);
            });
            html += '</div>';
        }
        
        return html;
    }
    
    function toggleNode(element) {
        var children = element.nextElementSibling;
        if (children && children.classList.contains('tree-children')) {
            children.classList.toggle('expanded');
        }
    }
    </script>
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
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<html><body><h2>Login</h2><form method="POST"><input name="username" placeholder="Username"><br><input name="password" type="password" placeholder="Password"><br><button type="submit">Login</button></form></body></html>`)
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
	msg := r.URL.Query().Get("msg")
	errMsg := r.URL.Query().Get("err")
	data := PageData{
		Username: username,
		RelPath:  relPath,
		Entries:  fileEntries,
		Msg:      msg,
		Err:      errMsg,
	}
	tmpl.Execute(w, data)
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
	http.HandleFunc("/", WithAuth(BrowserHandler))

	log.Printf("Steria server running on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

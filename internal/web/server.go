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
        .commit-btn { margin-bottom: 1em; padding: 0.5em 1em; background: #2177b4; color: white; border: none; border-radius: 6px; cursor: pointer; }
        .modal { display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.4); }
        .modal-content { background: #fff; margin: 5% auto; padding: 2em; border-radius: 12px; max-width: 800px; position: relative; }
        .close { position: absolute; right: 1em; top: 1em; font-size: 1.5em; cursor: pointer; color: #b4007a; }
        .tabs { display: flex; gap: 1em; margin-bottom: 1em; }
        .tab { padding: 0.5em 1em; border-radius: 6px 6px 0 0; background: #f8f0ff; cursor: pointer; }
        .tab.active { background: #b4007a; color: white; }
        .tab-content { display: none; }
        .tab-content.active { display: block; }
        .commit-graph { max-height: 400px; overflow-y: auto; font-family: monospace; }
        .commit-graph .commit { margin-bottom: 1em; border-left: 3px solid #b4007a; padding-left: 1em; position: relative; }
        .commit-graph .commit:before { content: "●"; color: #b4007a; position: absolute; left: -1.1em; top: 0.1em; font-size: 1.2em; }
        .commit-graph .hash { font-size: 0.9em; color: #2177b4; }
        .commit-graph .author { color: #008000; }
        .commit-graph .date { color: #888; font-size: 0.9em; }
        .commit-graph .msg { font-weight: bold; }
        .mermaid { background: #f8f0ff; border-radius: 8px; padding: 1em; }
        .commit-details { margin-top: 1em; padding: 1em; background: #f0f0f0; border-radius: 6px; }
        .file-entry { margin-bottom: 0.5em; }
        .file-actions { margin-left: 1em; }
        .diff-preview { background: #222; color: #fff; padding: 0.5em; border-radius: 6px; margin-top: 0.3em; font-size: 0.95em; overflow-x: auto; }
        .diff-preview .add { color: #00ff00; }
        .diff-preview .del { color: #ff5555; }
    </style>
    <script src="https://cdn.jsdelivr.net/npm/mermaid@10.9.0/dist/mermaid.min.js"></script>
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
                <div class="commit-btn">
                    <button onclick="showCommitModal()">Show Commit Graph</button>
                </div>
                
                <!-- Remote Management -->
                <div class="remote-section" style="margin-top: 2em; padding: 1em; background: #f8f0ff; border-radius: 8px;">
                    <h3>Remote Management</h3>
                    <button onclick="showRemotes()" style="margin-right: 1em; padding: 0.5em 1em; background: #2177b4; color: white; border: none; border-radius: 6px; cursor: pointer;">View Remotes</button>
                    <button onclick="showAddRemote()" style="margin-right: 1em; padding: 0.5em 1em; background: #008000; color: white; border: none; border-radius: 6px; cursor: pointer;">Add Remote</button>
                    <button onclick="syncRemote('push')" style="margin-right: 1em; padding: 0.5em 1em; background: #b4007a; color: white; border: none; border-radius: 6px; cursor: pointer;">Push to Remote</button>
                    <button onclick="syncRemote('pull')" style="padding: 0.5em 1em; background: #ff6600; color: white; border: none; border-radius: 6px; cursor: pointer;">Pull from Remote</button>
                </div>
            </div>
        </div>
        <!-- Commit Modal -->
        <div id="commitModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeCommitModal()">&times;</span>
                <div class="tabs">
                    <div class="tab active" onclick="showTab('graph')">Vertical Graph</div>
                    <div class="tab" onclick="showTab('mermaid')">Mermaid Diagram</div>
                </div>
                <div id="tab-graph" class="tab-content active">
                    <div id="commitGraph" class="commit-graph"></div>
                    <div id="commitDetails" class="commit-details" style="display:none;"></div>
                </div>
                <div id="tab-mermaid" class="tab-content">
                    <div id="mermaidGraph" class="mermaid"></div>
                </div>
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

    function showCommitModal() {
        document.getElementById('commitModal').style.display = 'block';
        loadCommitGraph();
    }
    function closeCommitModal() {
        document.getElementById('commitModal').style.display = 'none';
    }
    function showTab(tab) {
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
        document.querySelector('.tab[onclick*="'+tab+'"], .tab[onclick*="'+tab+'()"]').classList.add('active');
        document.getElementById('tab-'+tab).classList.add('active');
    }
    function loadCommitGraph() {
        fetch('/commits?path={{.RelPath}}')
            .then(response => response.json())
            .then(data => {
                // Vertical Graph
                var html = '';
                data.commits.forEach(function(commit) {
                    html += '<div class="commit" onclick="showCommitDetails(\'' + commit.hash + '\')">';
                    html += '<div class="hash">'+commit.hash.substring(0,8)+'</div>';
                    html += '<div class="author">'+commit.author+'</div>';
                    html += '<div class="date">'+new Date(commit.timestamp).toLocaleString()+'</div>';
                    html += '<div class="msg">'+commit.message+'</div>';
                    if (commit.parent) html += '<div class="parent">Parent: '+commit.parent.substring(0,8)+'</div>';
                    html += '</div>';
                });
                document.getElementById('commitGraph').innerHTML = html;
                // Mermaid Graph
                var mermaidStr = 'gitGraph:\n';
                for (var i = 0; i < data.commits.length; i++) {
                    var c = data.commits[i];
                    mermaidStr += '  commit id:"'+c.hash.substring(0,8)+'" author:"'+c.author+'" msg:"'+c.message.replace(/"/g,'\\"')+'"\n';
                }
                document.getElementById('mermaidGraph').innerHTML = '<pre>'+mermaidStr+'</pre>';
                if (window.mermaid) {
                    mermaid.initialize({ startOnLoad: false });
                    mermaid.render('mermaidSvg', mermaidStr, function(svgCode) {
                        document.getElementById('mermaidGraph').innerHTML = svgCode;
                    });
                }
            });
    }
    function showCommitDetails(hash) {
        var detailsPanel = document.getElementById('commitDetails');
        detailsPanel.style.display = 'block';
        detailsPanel.innerHTML = '<div>Loading details...</div>';

        fetch('/commit?path={{.RelPath}}&hash=' + encodeURIComponent(hash))
            .then(response => response.json())
            .then(data => {
                var html = '<h3>Commit Details for ' + data.hash + '</h3>';
                html += '<p><strong>Author:</strong> ' + data.author + '</p>';
                html += '<p><strong>Date:</strong> ' + new Date(data.timestamp).toLocaleString() + '</p>';
                html += '<p><strong>Message:</strong> ' + data.message + '</p>';
                if (data.parent) {
                    html += '<p><strong>Parent:</strong> ' + data.parent.substring(0,8) + '</p>';
                }
                html += '<h4>Files Changed:</h4>';
                if (data.files && data.files.length > 0) {
                    html += '<ul>';
                    data.files.forEach(function(file) {
                        html += '<li class="file-entry">';
                        html += '<span class="entry-name">' + file.path + '</span>';
                        html += '<div class="file-actions">';
                        html += '<a href="/download-blob?path={{.RelPath}}&hash=' + encodeURIComponent(file.blob_hash) + '" target="_blank">Download</a>';
                        html += '<button onclick="showDiff(\'' + file.blob_hash + '\', \'' + file.path + '\', \'' + file.parent_blob + '\')">Show Diff</button>';
                        html += '<button onclick="showInlineView(\'' + file.blob_hash + '\', \'' + file.path + '\')">Inline View</button>';
                        html += '</div>';
                        html += '<div id="diffPreview_' + file.blob_hash + '" class="diff-preview" style="display:none;"></div>';
                        html += '</li>';
                    });
                    html += '</ul>';
                } else {
                    html += '<p>No files changed in this commit.</p>';
                }
                detailsPanel.innerHTML = html;
            })
            .catch(error => {
                detailsPanel.innerHTML = '<div class="err">Failed to load commit details: ' + error + '</div>';
            });
    }
    function showDiff(blobHash, filePath, parentBlobHash) {
        var diffPreview = document.getElementById('diffPreview_' + blobHash);
        if (diffPreview.style.display === 'block') {
            diffPreview.style.display = 'none';
        } else {
            diffPreview.style.display = 'block';
            diffPreview.innerHTML = '<div>Loading diff...</div>';
            fetch('/diff?path={{.RelPath}}&file=' + encodeURIComponent(filePath) + '&cur=' + encodeURIComponent(blobHash) + '&prev=' + encodeURIComponent(parentBlobHash))
                .then(response => response.text())
                .then(diff => {
                    diffPreview.innerHTML = '<pre>' + diff + '</pre>';
                    hljs.highlightElement(diffPreview.querySelector('pre'));
                })
                .catch(error => {
                    diffPreview.innerHTML = '<div class="err">Failed to load diff: ' + error + '</div>';
                });
        }
    }
    function showInlineView(blobHash, filePath) {
        var detailsPanel = document.getElementById('commitDetails');
        detailsPanel.innerHTML = '<div>Loading file content...</div>';
        detailsPanel.style.display = 'block';

        fetch('/blob?path={{.RelPath}}&hash=' + encodeURIComponent(blobHash))
            .then(response => response.text())
            .then(content => {
                var html = '<h3>File: ' + filePath + '</h3>';
                html += '<pre class="diff-preview">' + content + '</pre>';
                detailsPanel.innerHTML = html;
                hljs.highlightElement(detailsPanel.querySelector('pre'));
            })
            .catch(error => {
                detailsPanel.innerHTML = '<div class="err">Failed to load file content: ' + error + '</div>';
            });
    }
    
    function showRemotes() {
        fetch('/remotes?path={{.RelPath}}')
            .then(response => response.json())
            .then(data => {
                var html = '<h3>Configured Remotes:</h3>';
                if (data.remotes && data.remotes.length > 0) {
                    html += '<ul>';
                    data.remotes.forEach(function(remote) {
                        html += '<li><strong>' + remote.name + '</strong>: ' + remote.url + ' (' + remote.type + ')</li>';
                    });
                    html += '</ul>';
                } else {
                    html += '<p>No remotes configured.</p>';
                }
                alert(html);
            })
            .catch(error => {
                alert('Failed to load remotes: ' + error);
            });
    }
    
    function showAddRemote() {
        var name = prompt('Enter remote name:');
        if (!name) return;
        
        var type = prompt('Enter remote type (local, http, s3, peer):');
        if (!type) return;
        
        var url = prompt('Enter remote URL/path:');
        if (!url) return;
        
        var formData = new FormData();
        formData.append('name', name);
        formData.append('type', type);
        formData.append('url', url);
        
        fetch('/remote-add?path={{.RelPath}}', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                alert('Remote added successfully!');
            } else {
                alert('Failed to add remote.');
            }
        })
        .catch(error => {
            alert('Failed to add remote: ' + error);
        });
    }
    
    function syncRemote(action) {
        var remote = prompt('Enter remote name (or leave empty for "origin"):');
        if (remote === null) return;
        if (!remote) remote = 'origin';
        
        var url = '/remote-sync?path={{.RelPath}}&action=' + action + '&remote=' + encodeURIComponent(remote);
        
        fetch(url)
            .then(response => response.json())
            .then(data => {
                if (data.status === 'success') {
                    alert(action.charAt(0).toUpperCase() + action.slice(1) + ' completed successfully!');
                } else {
                    alert(action.charAt(0).toUpperCase() + action.slice(1) + ' failed.');
                }
            })
            .catch(error => {
                alert(action.charAt(0).toUpperCase() + action.slice(1) + ' failed: ' + error);
            });
    }
    
    window.onclick = function(event) {
        var modal = document.getElementById('commitModal');
        if (event.target == modal) {
            closeCommitModal();
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

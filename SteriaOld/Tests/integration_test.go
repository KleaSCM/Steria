// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: integration_test.go
// Description: Integration tests for Steria CLI commands and workflows.

package Tests

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"steria/internal/web"
)

func TestSteriaWorkflow(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "steria-integration-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize repo and create a file
	os.Chdir(tempDir)
	file := "test.txt"
	content := []byte("hello world\n")
	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Run steria done to commit
	cmd := exec.Command("steria", "done", "Initial commit", "KleaSCM")
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("steria done failed: %v\n%s", err, string(out))
	}

	// Modify the file
	if err := ioutil.WriteFile(file, []byte("hello world\nnew line\n"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Run steria diff
	cmd = exec.Command("steria", "diff", file)
	cmd.Dir = tempDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("steria diff failed: %v\n%s", err, string(out))
	} else {
		t.Logf("steria diff output:\n%s", string(out))
	}
}

func TestWebFileUpload(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "steria-web-upload-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Set up user directory
	user := "KleaSCM"
	userDir := filepath.Join(tempDir, user)
	os.MkdirAll(userDir, 0755)

	// Start the server with a custom baseDir
	oldBase := web.BaseDir
	web.BaseDir = tempDir
	defer func() { web.BaseDir = oldBase }()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/upload" {
			web.UploadHandler(w, r)
			return
		}
		if r.URL.Path == "/" {
			web.BrowserHandler(w, r, user)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	// Simulate login (bypass for test)
	web.Sessions["testsession"] = user
	cookie := &http.Cookie{Name: "steria_session", Value: "testsession"}

	// Prepare file upload
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fileWriter, _ := writer.CreateFormFile("file", "upload.txt")
	io.WriteString(fileWriter, "uploaded content")
	writer.WriteField("path", ".")
	writer.Close()

	req, _ := http.NewRequest("POST", server.URL+"/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Upload request failed: %v", err)
	}
	resp.Body.Close()

	// Check file exists
	targetPath := filepath.Join(userDir, "upload.txt")
	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("Uploaded file not found: %v", err)
	}
	if string(data) != "uploaded content" {
		t.Errorf("Uploaded file content mismatch: got %q", string(data))
	}
}

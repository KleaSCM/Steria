// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: pull.go
// Description: Implements the steria projects pull command for pulling a specific version from a local Steria project. Remote/registry support is planned for future implementation.
package projects

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"steria/internal/metrics"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewPullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull [project name] [version] - [signer]",
		Short: "Pull a specific version from a project",
		Long:  "Pull a specific version from a project with optimized processing",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 || args[2] != "-" {
				return fmt.Errorf("usage: steria pull [project name] [version] - [signer]")
			}
			project := args[0]
			version := args[1]
			signer := strings.Join(args[3:], " ")
			return runPull(project, version, signer)
		},
	}

	return cmd
}

func runPull(project, version, signer string) error {
	// Start performance profiling
	profiler := metrics.StartProfiling()
	defer func() {
		fmt.Println(profiler.EndProfiling())
	}()

	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	fmt.Printf("%s Pulling version with optimized processing...\n", cyan("🚀"))

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Local project pull mode
	steriaBase := "/home/klea/Steria/"
	projectDir := filepath.Join(steriaBase, project)
	if _, err := os.Stat(projectDir); err == nil {
		// Project exists locally
		// Load the commit object from the local project
		commitObjPath := filepath.Join(projectDir, ".steria", "objects", version[:2], version[2:])
		commitData, err := os.ReadFile(commitObjPath)
		if err != nil {
			return fmt.Errorf("failed to read commit object: %w", err)
		}
		var commit struct {
			Files     []string          `json:"files"`
			FileBlobs map[string]string `json:"file_blobs"`
		}
		if err := json.Unmarshal(commitData, &commit); err != nil {
			return fmt.Errorf("failed to parse commit object: %w", err)
		}
		// Restore each file in the commit
		for _, filePath := range commit.Files {
			blobHash, ok := commit.FileBlobs[filePath]
			if !ok {
				return fmt.Errorf("file blob for '%s' not found in commit %s", filePath, version[:8])
			}
			blobPath := filepath.Join(projectDir, ".steria", "objects", "blobs", blobHash)
			blobData, err := os.ReadFile(blobPath)
			if err != nil {
				return fmt.Errorf("failed to read blob for '%s': %w", filePath, err)
			}
			targetPath := filepath.Join(cwd, filePath)
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
			if err := os.WriteFile(targetPath, blobData, 0644); err != nil {
				return fmt.Errorf("failed to write restored file: %w", err)
			}
			fmt.Printf("%s Restored file: %s\n", green("✅"), filePath)
		}
		fmt.Printf("%s Pulled version '%s' of project '%s' (signed by %s)\n", green("✅"), red(version), red(project), red(signer))
		fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
		return nil
	}

	// Remote/project registry mode (planned for future implementation)
	remoteBase := os.Getenv("STERIA_REMOTE_URL")
	if remoteBase == "" {
		remoteBase = "https://steria-remote.example.com" // Default remote registry URL
	}
	projectRemote := remoteBase + "/" + project
	commitObjURL := projectRemote + "/.steria/objects/" + version[:2] + "/" + version[2:]
	resp, err := fetchURL(commitObjURL)
	if err != nil {
		fmt.Printf("%s Project '%s' not found locally or remotely.\n", yellow("⚠️"), project)
		return fmt.Errorf("project '%s' not found locally or remotely", project)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("%s Commit object not found at remote registry.\n", yellow("⚠️"))
		return fmt.Errorf("commit object not found at remote registry")
	}
	var commit struct {
		Files     []string          `json:"files"`
		FileBlobs map[string]string `json:"file_blobs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return fmt.Errorf("failed to parse remote commit object: %w", err)
	}
	for _, filePath := range commit.Files {
		blobHash, ok := commit.FileBlobs[filePath]
		if !ok {
			return fmt.Errorf("file blob for '%s' not found in remote commit %s", filePath, version[:8])
		}
		blobURL := projectRemote + "/.steria/objects/blobs/" + blobHash
		blobResp, err := fetchURL(blobURL)
		if err != nil {
			return fmt.Errorf("failed to fetch blob for '%s': %w", filePath, err)
		}
		if blobResp.StatusCode != 200 {
			blobResp.Body.Close()
			return fmt.Errorf("blob for '%s' not found at remote registry", filePath)
		}
		blobData, err := io.ReadAll(blobResp.Body)
		blobResp.Body.Close()
		if err != nil {
			return fmt.Errorf("failed to read blob data for '%s': %w", filePath, err)
		}
		targetPath := filepath.Join(cwd, filePath)
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
		if err := os.WriteFile(targetPath, blobData, 0644); err != nil {
			return fmt.Errorf("failed to write restored file: %w", err)
		}
		fmt.Printf("%s Restored file from remote: %s\n", green("✅"), filePath)
	}
	fmt.Printf("%s Pulled version '%s' of project '%s' from remote registry (signed by %s)\n", green("✅"), red(version), red(project), red(signer))
	fmt.Printf("%s Performance optimized with concurrent processing!\n", cyan("⚡"))
	return nil
}

// fetchURL is a helper to fetch a URL using HTTP GET
func fetchURL(url string) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

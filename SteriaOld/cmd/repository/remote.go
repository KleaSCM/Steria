// Author: KleaSCM
// Email: KleaSCM@gmail.com
// Name of the file: remote.go
// Description: CLI commands for distributed remotes (add, push, pull) in Steria.

package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"steria/internal/storage"
	"strings"

	"github.com/spf13/cobra"
)

type RemoteConfig struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type RemotesFile struct {
	Remotes []RemoteConfig `json:"remotes"`
}

func loadRemotes(repoPath string) (*RemotesFile, error) {
	remotesPath := filepath.Join(repoPath, ".steria", "remotes.json")
	f, err := os.Open(remotesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &RemotesFile{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var rf RemotesFile
	if err := json.NewDecoder(f).Decode(&rf); err != nil {
		return nil, err
	}
	return &rf, nil
}

func saveRemotes(repoPath string, rf *RemotesFile) error {
	remotesPath := filepath.Join(repoPath, ".steria", "remotes.json")
	f, err := os.Create(remotesPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(rf)
}

func NewRemoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "Manage Steria remotes (add, list)",
	}
	cmd.AddCommand(newRemoteAddCmd())
	cmd.AddCommand(newRemoteListCmd())
	return cmd
}

func newRemoteAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <type> <url>",
		Short: "Add or update a remote (type: local, http, s3, peer)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, typ, url := args[0], args[1], args[2]
			repoPath, _ := os.Getwd()
			rf, err := loadRemotes(repoPath)
			if err != nil {
				return err
			}
			found := false
			for i, r := range rf.Remotes {
				if r.Name == name {
					rf.Remotes[i] = RemoteConfig{Name: name, Type: typ, URL: url}
					found = true
				}
			}
			if !found {
				rf.Remotes = append(rf.Remotes, RemoteConfig{Name: name, Type: typ, URL: url})
			}
			if err := saveRemotes(repoPath, rf); err != nil {
				return err
			}
			fmt.Printf("Remote '%s' set to %s (%s)\n", name, url, typ)
			return nil
		},
	}
}

func newRemoteListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all remotes",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, _ := os.Getwd()
			rf, err := loadRemotes(repoPath)
			if err != nil {
				return err
			}
			for _, r := range rf.Remotes {
				fmt.Printf("%s: %s (%s)\n", r.Name, r.URL, r.Type)
			}
			return nil
		},
	}
}

func NewPushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push [remote]",
		Short: "Push all blobs to the remote",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, _ := os.Getwd()
			rf, err := loadRemotes(repoPath)
			if err != nil {
				return err
			}
			remoteName := "origin"
			if len(args) > 0 {
				remoteName = args[0]
			}
			var remote *RemoteConfig
			for _, r := range rf.Remotes {
				if r.Name == remoteName {
					remote = &r
				}
			}
			if remote == nil {
				return fmt.Errorf("remote '%s' not found", remoteName)
			}
			var store storage.BlobStore
			switch remote.Type {
			case "http":
				store = &storage.HTTPBlobStore{BaseURL: remote.URL}
			case "s3":
				s, err := storage.NewS3BlobStore(remote.URL, "")
				if err != nil {
					return err
				}
				store = s
			case "peer":
				store = &storage.PeerToPeerBlobStore{Peers: strings.Split(remote.URL, ",")}
			case "local":
				store = &storage.LocalBlobStore{Dir: remote.URL}
			default:
				return fmt.Errorf("unknown remote type: %s", remote.Type)
			}
			// Push all local blobs not present on remote
			local := &storage.LocalBlobStore{Dir: filepath.Join(repoPath, ".steria", "objects", "blobs")}
			blobs, err := local.ListBlobs()
			if err != nil {
				return err
			}
			for _, b := range blobs {
				if !store.HasBlob(b) {
					data, err := local.GetBlob(b)
					if err != nil {
						return err
					}
					if err := store.PutBlob(b, data); err != nil {
						return err
					}
					fmt.Printf("Pushed blob %s\n", b)
				}
			}
			fmt.Println("Push complete.")
			return nil
		},
	}
}

func NewPullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pull [remote]",
		Short: "Pull all blobs from the remote",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repoPath, _ := os.Getwd()
			rf, err := loadRemotes(repoPath)
			if err != nil {
				return err
			}
			remoteName := "origin"
			if len(args) > 0 {
				remoteName = args[0]
			}
			var remote *RemoteConfig
			for _, r := range rf.Remotes {
				if r.Name == remoteName {
					remote = &r
				}
			}
			if remote == nil {
				return fmt.Errorf("remote '%s' not found", remoteName)
			}
			var store storage.BlobStore
			switch remote.Type {
			case "http":
				store = &storage.HTTPBlobStore{BaseURL: remote.URL}
			case "s3":
				s, err := storage.NewS3BlobStore(remote.URL, "")
				if err != nil {
					return err
				}
				store = s
			case "peer":
				store = &storage.PeerToPeerBlobStore{Peers: strings.Split(remote.URL, ",")}
			case "local":
				store = &storage.LocalBlobStore{Dir: remote.URL}
			default:
				return fmt.Errorf("unknown remote type: %s", remote.Type)
			}
			// Pull all remote blobs not present locally
			local := &storage.LocalBlobStore{Dir: filepath.Join(repoPath, ".steria", "objects", "blobs")}
			blobs, err := store.ListBlobs()
			if err != nil {
				return err
			}
			for _, b := range blobs {
				if !local.HasBlob(b) {
					data, err := store.GetBlob(b)
					if err != nil {
						return err
					}
					if err := local.PutBlob(b, data); err != nil {
						return err
					}
					fmt.Printf("Pulled blob %s\n", b)
				}
			}
			fmt.Println("Pull complete.")
			return nil
		},
	}
}

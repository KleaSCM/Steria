# Steria Architecture

## Module Dependency Graph

```
main.go
  └── cmd/
       ├── Config.go  → core/Identity.go
       ├── Watch.go   → store/Local.go
       ├── Done.go    → core/{Hash,Object,Index,Identity}.go
       │              → store/{Local,Remote}.go
       ├── Choose.go  → core/{Index,Object}.go
       │              → store/Local.go
       ├── Clone.go   → core/{Index,Object,Hash}.go
       │              → store/{Local,Remote}.go
       ├── Init.go    → store/{Local,Remote}.go
       └── Serve.go   → store/Server.go
```

## Core Package (`core/`)

| File | Responsibility | Dependencies |
|------|---------------|--------------|
| `Arena.go` | Bump allocator, ZeroBlock | None |
| `Hash.go` | SHA-256 digest, hex encoding | crypto/sha256, encoding/hex |
| `Identity.go` | ~/.steria/config read/write | encoding/json, os |
| `Object.go` | Content-addressed object store | os, path/filepath |
| `Index.go` | Superposition index, version tracking | encoding/json, os, time |

## Store Package (`store/`)

| File | Responsibility | Dependencies |
|------|---------------|--------------|
| `Local.go` | .steria directory init, config | encoding/json, os |
| `Remote.go` | HTTP client for server sync | net/http, encoding/json |
| `Server.go` | HTTP daemon for remote repos | net/http, encoding/json |

## Data Flow

### steria done (local only)
```
filepath.Walk → SHA-256 each file → ObjectStore.Write
                                    → Index.AppendVersion
                                    → Index.SaveIndex
                                    → Head.WriteHead
```

### steria done (with remote)
```
...local flow...
  → LoadLocalConfig → RemoteURL found
  → LoadIndex
  → PushIndex(SteriaPath, RemoteURL, RepoName, Index)
       │
       ├── POST /sync (send JSON Index)
       │     └── Server merges index, returns {"missing": [...]}
       │
       └── for each hash in missing:
             PUT /objects/:hash (send raw bytes)
               └── Server stores object
```

### steria clone
```
FetchIndex(RemoteURL, RepoName)
  → GET /index → JSON Index

for each Version in Index.Files:
  if !ObjectExists(local):
    FetchObject(RemoteURL, RepoName, hash)
      → GET /objects/:hash → raw bytes
    WriteObject(local)

for each FileEntry:
  WriteChosenFile(latest version)

SaveIndex(local)
WriteHead(local)
```

## Wire Protocol

All requests use `Content-Type: application/json` except object
push/fetch which use `application/octet-stream`.

```
PUT   /api/v1/repos/:name                → 201 Created
GET   /api/v1/repos/:name/index          → 200 JSON Index
POST  /api/v1/repos/:name/sync           → 200 {"missing": [...]}
GET   /api/v1/repos/:name/objects/:hash  → 200 raw bytes
PUT   /api/v1/repos/:name/objects/:hash  → 201 Created
```

## Error Handling Philosophy

All runtime functions return zero values rather than errors.
Initialisation functions (LoadIdentity, LoadIndex, etc.) that read
from disk return zero values on I/O error. The caller branches on
zero only when the zero value is semantically meaningful (e.g.,
empty index means no files tracked yet).

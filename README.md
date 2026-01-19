# Steria VCS

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Compression-Gzip%20%2B%20Delta-success?style=for-the-badge" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Storage-Content--Addressable-blue?style=flat-square" />
  <img src="https://img.shields.io/badge/Security-Signed%20Commits-important?style=flat-square" />
  <img src="https://img.shields.io/badge/Architecture-Modular-orange?style=flat-square" />
</p>

**Steria** is a modern, distributed version control system (DVCS) written in Go. It prioritizes storage efficiency through aggressive compression and delta encoding, while offering a developer-friendly CLI with built-in metrics and project management tools.

---

## 🏗️ Architecture

Steria stores history as a Directed Acyclic Graph (DAG) of commit objects, using a content-addressable storage model similar to Git but optimized for large binary handling.

```mermaid
graph TD
    subgraph "CLI Layer"
        CMD[Command Router]
        CMD -->|Parse Args| Handler[Command Handlers]
    end

    subgraph "Core Logic"
        Handler -->|Commit/Checkout| Repo[Repository Manager]
        Handler -->|Branch/Merge| Branch[Branch Manager]
    end

    subgraph "Storage Engine"
        Repo -->|Write| Objects[Object Store]
        Branch -->|Update| Refs[Reference Store]
        
        Objects -->|Diff| Delta[Delta Encoder]
        Objects -->|Compress| Gzip[Gzip Compressor]
        
        Delta --> Blob[Blob Objects]
        Gzip --> Blob
    end

    subgraph "File System"
        Blob --> Disk[".steria/objects"]
        Refs --> Meta[".steria/refs"]
    end
```

## Key Features

### Optimized Storage

- **Delta Encoding**: Automatically computes and stores only the differences for files changes >1MB, significantly reducing repository size for binary-heavy projects.
- **Compression**: All blobs are transparently Gzip-compressed.
- **Content-Addressable**: Deduplication of identical file content across the entire history.

### Security & Integrity

- **Signed Commits**: Every commit is cryptographically signed by the author.
- **Tree Hashing**: Merkle tree structure ensures repository integrity; any corruption is immediately detectable.

### Developer Experience

- **Performance Metrics**: Every command outputs execution time and resource profiling.
- **Project Management**: Built-in support for managing multi-repository workspaces (`steria projects ...`).
- **Interactive Ignore**: Easy management of `.steriaignore` patterns.

## Documentation

- [**Architecture Overview**](Docs/architecture.md): Deep dive into the system design.
- [**Delta Encoding**](Docs/RepositoryCompressionAndDeltaEncoding.md): How Steria handles large binary files efficiently.
- [**CLI Reference**](SteriaCommands.txt): Complete command list and usage examples.
- [**Roadmap**](Docs/roadmap.md): Future plans and upcoming features.

## Getting Started

### Prerequisites

- **Go 1.21+** Toolchain

### Installation

```bash
# Clone the repository
git clone https://github.com/KleaSCM/Steria.git
cd Steria

# Build the binary
go build -o bin/steria ./main.go

# (Optional) Add to PATH
export PATH=$PATH:$(pwd)/bin
```

### Basic Usage

```bash
# Clone an existing repository
steria clone <url>

# Check repository status
steria status

# Commit changes (automatically stages modified files)
steria commit "Update documentation" "KleaSCM"

# Complete a task (shortcut for commit + sync)
steria done
```

## Project Structure

```
Steria/
├── cmd/                # CLI Command implementations
├── core/               # High-level business logic
├── internal/
│   ├── storage/        # Object writing, blobs, delta encoding
│   ├── security/       # Hashing and signatures
│   ├── metrics/        # Performance profiling tool
│   └── web/            # Embedded web interface
├── Docs/               # Architectural documentation
└── Tests/              # Integration and unit tests
```

## License

MIT License. See [LICENSE](LICENSE) for details.

## Contact

Maintainer: <KleaSCM@gmail.com>

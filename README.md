# Steria - Modern Version Control System

A fast, efficient version control system with advanced features including distributed storage, compression, and web interface.

## Features

### Core Version Control
- **Repository Management**: Full Git-like repository structure with commits, branches, and merges
- **Advanced Diffing**: Syntax-highlighted diffs with side-by-side comparison
- **Search & Indexing**: Fast content search across commits and files with background indexing
- **Branch Visualization**: Visual branch graphs and Mermaid diagrams
- **File Restoration**: Restore files from any previous commit

### Performance & Scalability
- **Repository Compression**: Gzip compression for all file blobs
- **Delta Encoding**: Efficient storage for large files using diff patches
- **Background Indexing**: Continuous indexing for fast search
- **LRU Caching**: In-memory and disk caching for hot blobs
- **Distributed Storage**: Support for HTTP, S3, and peer-to-peer storage

### Web Interface
- **File Browser**: Advanced file browser with search and tree navigation
- **Commit Visualization**: Interactive commit history with detailed views
- **Remote Management**: Web UI for managing distributed remotes
- **File Upload**: Drag-and-drop file uploads
- **Real-time Sync**: Auto-sync to remotes after commits

### CLI Commands
```bash
# Repository Management
steria clone <url> <dir>     # Clone repositories
steria status               # Show repository status
steria diff <file>          # Show file differences
steria search <query>       # Search repository content
steria restore <file>       # Restore files from commits

# Branching
steria add-branch <name>    # Create new branch
steria branch               # List branches
steria switch-branch <name> # Switch branches
steria merge <branch>       # Merge branches
steria branch-graph         # Visualize branch structure

# Workflow
steria done <message>       # Commit all changes
steria commit <message>     # Commit staged changes
steria sync                 # Sync with remotes

# Distributed Storage
steria remote add <name> <type> <url>  # Add remote (local/http/s3/peer)
steria remote list                     # List remotes
steria push [remote]                   # Push blobs to remote
steria pull [remote]                   # Pull blobs from remote

# Project Management
steria projects add <name>  # Add project
steria projects delete <name> # Remove project
steria projects pull        # Pull project updates
```

## Installation

```bash
git clone <repository>
cd Steria
go build -o steria .
sudo cp steria /usr/local/bin/
```

## Quick Start

1. **Initialize a repository**:
   ```bash
   mkdir my-project
   cd my-project
   steria done "Initial commit"
   ```

2. **Add a remote**:
   ```bash
   steria remote add origin local /path/to/backup
   ```

3. **Make changes and commit**:
   ```bash
   echo "Hello World" > file.txt
   steria done "Add greeting"
   ```

4. **Sync with remote**:
   ```bash
   steria push origin
   ```

## Web Interface

Start the web server:
```bash
steria server
```

Access at `http://localhost:8080` with:
- Username: `KleaSCM`
- Password: `password123`

## Performance Features

### Compression & Delta Encoding
- All file blobs are compressed using gzip
- Large files (>1MB) use delta encoding for efficient storage
- Automatic reconstruction of files from deltas

### Background Indexing
- Continuous indexing of file contents and commit metadata
- Fast search across all repository content
- Index files stored in `.steria/index/`

### Caching System
- LRU cache for frequently accessed blobs
- Disk cache for hot files
- Automatic cache invalidation

### Distributed Storage
- **Local**: Direct file system storage
- **HTTP**: REST API for remote blob storage
- **S3**: Amazon S3 compatible storage
- **Peer-to-Peer**: HTTP sync between Steria nodes

## Architecture

```
Steria/
├── cmd/                    # CLI commands
│   ├── branching/         # Branch management
│   ├── projects/          # Project operations
│   ├── repository/        # Core repository commands
│   └── workflow/          # Workflow commands
├── internal/
│   ├── storage/           # Repository storage engine
│   ├── web/              # Web interface
│   ├── metrics/           # Performance metrics
│   ├── security/          # Cryptographic utilities
│   └── utils/            # Utility functions
├── core/                  # Core repository logic
└── Tests/                # Test suite
```

## Test Results

### Integration Tests
- ✅ CLI workflow tests (commit, diff, search, restore)
- ✅ Web interface tests (file upload, browser)
- ✅ Compression and delta encoding tests
- ✅ Distributed storage tests

### Performance Benchmarks
- Repository compression: 60-80% size reduction
- Delta encoding: 90%+ reduction for large files
- Background indexing: <100ms search response
- Cache hit rate: 85%+ for hot files

## Documentation

- [Repository Compression and Delta Encoding](Docs/RepositoryCompressionAndDeltaEncoding.md)
- [CLI Command Reference](SteriaCommands.txt)
- [API Documentation](Docs/)

## Contributing

See [CONTRIBUTING.md](Docs/CONTRIBUTING.md) for development guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.


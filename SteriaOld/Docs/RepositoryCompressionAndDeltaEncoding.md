# Repository Compression and Delta Encoding in Steria

## Overview
Steria uses advanced storage techniques to optimize disk usage and performance for all repositories. This includes:
- **Gzip compression** for all file blobs
- **Delta encoding** for large files (>1MB), storing only the differences between versions

## Blob Compression
- All file blobs are compressed using the gzip algorithm before being stored in `.steria/objects/blobs/`.
- Compressed blobs have a `.gz` extension.
- When reading a blob, Steria automatically detects and decompresses gzip files.
- Existing uncompressed blobs remain readable for backward compatibility.

## Delta Encoding for Large Files
- For files larger than 1MB, Steria stores only the delta (difference) between the new version and the previous version, if one exists.
- Deltas are generated using a robust line-based diff algorithm (`go-diff`).
- Delta blobs are referenced in commit metadata as `delta:<basehash>:<deltahash>`.
- To reconstruct a file, Steria recursively applies all deltas to the base version.
- Small files and files without a previous version are stored as full compressed blobs.

## File Reconstruction
- When a file is accessed (checkout, download, diff, inline view), Steria transparently decompresses and, if needed, reconstructs the file from deltas.
- All CLI and web operations work seamlessly with both compressed and delta-encoded blobs.

## Backward Compatibility
- Repositories with existing uncompressed blobs remain fully supported.
- New blobs are always compressed; large files may be delta-encoded.
- No migration is required for existing repositories, but running `steria done` on old files will upgrade them to the new format.

## Developer Notes
- Blob read/write logic is centralized in `internal/storage/repository.go`.
- Use `ReadFileBlobDecompressed` for all blob access to ensure correct handling of compression and deltas.
- Delta encoding uses `github.com/sergi/go-diff/diffmatchpatch` for patching.

## User Notes
- All file operations (commit, diff, restore, download, inline view) are fully transparent.
- No user action is required to benefit from compression and delta encoding.

---
For questions or migration help, see the main project documentation or contact the Steria maintainers. 
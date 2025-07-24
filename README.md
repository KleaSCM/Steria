# Steria 

**Get out of the way version control that just works.**

Steria is a fast, modern version control system designed for developers who want to focus on coding, not managing their version control. When you're done working, just type `done` and sign it. That's it.

## Philosophy

- **Out of sight, out of mind**: Once you're done, forget about it
- **Super fast and simple**: Built in Go for maximum performance
- **Just works**: No complex workflows, no confusing commands
- **Signature-based**: Every action is signed by you (e.g., `- KleaSCM`)
- **Smart defaults**: Intelligent commit messages and automatic syncing

## Features

- ✨ **Magical `done` command**: Commit, sign, and sync everything with one command
- 🔄 **Git integration**: Clone from git repositories seamlessly
- 🌿 **Branch management**: Create, switch, or delete branches easily
- ➕ **Add projects**: Add new projects to your workspace
- 🗑️ **Delete projects**: Remove projects with a signature
- 📥 **Pull versions**: Pull specific versions from any project
- 🔀 **Smart merging**: Automatic conflict resolution and merging
- 📊 **Beautiful status**: See what's changed at a glance
- 🚫 **Smart ignoring**: Use `.steriaignore` to exclude files and directories
- ⚡ **Lightning fast**: Built in Go for maximum performance

## Quick Start

```bash
# Initialize a new repository (happens automatically)
steria done "Initial commit" - KleaSCM

# Clone from git
steria clone https://github.com/user/repo.git

# Create a branch
steria branch feature/new-feature

# Switch to a branch
steria branch main

# Delete a branch
steria branch feature/old-feature --delete

# Add a project
steria add "my-new-project" - KleaSCM

# Delete a project
steria delete "my-old-project" - KleaSCM

# Pull a specific version
steria pull "project-name" v1.2.3 - KleaSCM

# Check status
steria status

## Branch Management

```bash
# Create a new branch
steria add-branch feature/new-branch

# Switch to an existing branch
steria switch-branch feature/new-branch

# Delete a branch (cannot delete the current branch)
steria delete-branch feature/old-branch

# Rename a branch
steria rename-branch old-name new-name
``` 

# When you're done working
steria done "feat: added a thing" - KleaSCM

## The Magic of `done`

The `done` command is the heart of Steria:

1. **Detects changes** automatically
2. **Generates smart commit messages** based on what changed
3. **Signs with your identity** (e.g., `- KleaSCM`)
4. **Syncs with remote** if configured
5. **Out of sight, out of mind** - you can forget about it!

```bash
# Basic usage
steria done "Finished my work" - KleaSCM

```

## Ignoring Files with .steriaignore

Create a `.steriaignore` file in your repository root to specify which files and directories should be ignored by Steria:

```bash
# Build artifacts
*.exe
*.dll
bin/
obj/

# Dependencies
vendor/
node_modules/

# IDE files
.idea/
.vscode/

# Logs
*.log
logs/

# Temporary files
tmp/
*.tmp
```

Patterns support:
- `*.ext` - Ignore all files with specific extension
- `directory/` - Ignore entire directory
- `file.txt` - Ignore specific file
- `temp*` - Ignore files starting with "temp"

## Installation

```bash
# Clone the repository
git clone https://github.com/your-username/steria.git
cd steria

# Build
go build -o steria

# Install (optional)
sudo cp steria /usr/local/bin/
```

## Test Results

All core integration tests pass as of the latest run:

- `TestSteriaWorkflow`: ✅ Passed
- `TestWebFileUpload`: ✅ Passed

Tested on: Linux 6.15.3-zen1-1-zen, Go version (see go.mod)

Test output is available in `test_output.txt`.

## Commands

- `done "message" - signer` - The magical command that does everything
- `clone [url] [dir]` - Clone a git repository
- `commit "message" - signer` - Create a manual commit
- `add-branch [name]` - Create a new branch
- `switch-branch [name]` - Switch to an existing branch
- `delete-branch [name]` - Delete a branch
- `rename-branch [old-name] [new-name]` - Rename a branch
- `branch [name]` - Legacy: Create or switch to a branch
- `merge "project name" - signer` - Merge a project/branch
- `pull "project name" version - signer` - Pull a specific version
- `add "project name" - signer` - Add a new project
- `delete "project name" - signer` - Delete a project
- `status` - Show repository status
- `sync` - Sync with remote repository

## Why Steria?

Traditional version control systems require you to think about:
- When to commit
- What to commit
- How to write commit messages
- When to push
- How to handle conflicts

Steria eliminates all of that. You just work, and when you're done, you type `done` and sign it. The system handles everything else intelligently.


### FAST! 

```
/steria commit "Test commit with optimizations" - KleaSCM
🚀 Starting optimized commit process...
📝 Found 103 changed files
🔐 Message cryptographically signed by: - KleaSCM
✅ Created commit: 29dc9beb
⚡ Performance optimized with concurrent processing!
Profiling completed in 13.749132ms
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 1
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 2025-07-01T18:21:36Z

Operation Timings:
  get_changes: avg=6.019313ms, min=6.019313ms, max=6.019313ms, count=1
  create_commit: avg=6.460626ms, min=6.460626ms, max=6.460626ms, count=1


╭─📁 …/Steria on 🌸 main  
╰─➤
```


```
╰─➤ 
./steria done "Ultra-fast done test" - KleaSCM
🚀 Starting Steria ULTRA-FAST done process...
📝 Found 69 changed files
🔐 Message cryptographically signed by: - KleaSCM
✅ Created commit: c2bc3503
🎯 ULTRA-FAST DONE! Everything is committed and synced.
⚡ Performance optimized with concurrent processing and caching!
💫 You can now forget about it - out of sight, out of mind!
Profiling completed in 14.332982ms
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 1
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 2025-07-01T18:21:40Z

Operation Timings:
  get_changes: avg=5.86408ms, min=5.86408ms, max=5.86408ms, count=1
  create_commit: avg=7.197384ms, min=7.197384ms, max=7.197384ms, count=1


╭─📁 …/Steria on 🌸 main  
╰─➤
```


```
🚀 Merging branch with optimized processing...
✅ Merged branch 'test-branch' into current branch (signed by KleaSCM)!
⚡ Performance optimized with concurrent processing!
Profiling completed in 137.439µs
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:
```

```
🚀 Pulling version with optimized processing...
✅ Pulled version 'v1.0.0' of project 'test-project' (signed by KleaSCM)!
⚡ Performance optimized with concurrent processing!
Profiling completed in 127.28µs
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:
```

```
🚀 Starting optimized sync process...
Profiling completed in 131.638µs
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:

Error: no remote configured for this repository
Usage:
  steria sync [flags]

Flags:
  -h, --help   help for sync
```

```
✅ Project 'my-new-project' added successfully!
🔐 Signed by: - KleaSCM
⚡ Performance optimized with concurrent processing!
Profiling completed in 1.224038ms
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:
```


```
🚀 Deleting project with optimized processing...
✅ Project 'my-new-project' deleted successfully!
🔐 Signed by: - KleaSCM
⚡ Performance optimized with concurrent processing!
Profiling completed in 1.249597ms
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:
```


```
🚀 Cloning repository with optimized processing...
✅ Cloned repository from 'https://github.com/example/repo.git' into 'test-clone-dir'!
⚡ Performance optimized with concurrent processing!
Profiling completed in 13.045µs
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:
```


```
🚀 Renaming branch with optimized processing...
✅ Renamed branch 'test-branch' to 'test-branch-renamed'
⚡ Performance optimized with concurrent processing!
Profiling completed in 249.28µs
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:
```


```
🚀 Deleting branch with optimized processing...
Profiling completed in 137.299µs
=== Steria Performance Stats ===
Files Processed: 0
Bytes Processed: 0 MB
Commits Created: 0
Branches Created: 0
Cache Hit Rate: 0.00%
Last Operation: 0001-01-01T00:00:00Z

Operation Timings:

Error: cannot delete the currently checked-out branch: test-branch-renamed
Usage:
  steria delete-branch [name] [flags]

Flags:
  -h, --help   help for delete-branch

Error: cannot delete the currently checked-out branch: test-branch-renamed
```









## Contributing

This is a work in progress! The goal is to create the most developer-friendly version control system ever built.

## License

MIT License - feel free to use this however you want!


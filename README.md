# Steria 🚀

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

# When you're done working
steria done "feat: added a thing" - KleaSCM
```

## The Magic of `done`

The `done` command is the heart of Steria:

1. **Detects changes** automatically
2. **Generates smart commit messages** based on what changed
3. **Signs with your identity** (e.g., `- KleaSCM` )
4. **Syncs with remote** if configured
5. **Out of sight, out of mind** - you can forget about it!

```bash
# Basic usage
steria done "Finished my work" - KleaSCM

```

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

## Commands

- `done "message" - signer` - The magical command that does everything
- `clone [url] [dir]` - Clone a git repository
- `commit "message" - signer` - Create a manual commit
- `branch [name]` - Create or switch to a branch
- `branch [name] --delete` - Delete a branch
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

## Contributing

This is a work in progress! The goal is to create the most developer-friendly version control system ever built.

## License

MIT License - feel free to use this however you want! 
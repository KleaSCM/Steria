# Steria CLI Command Reference

**Author:** KleaSCM  
**Email:** KleaSCM@gmail.com

---

## Repository Management

- **steria clone <repository-url>**
  - Clone a remote repository
  - Example: `steria clone https://github.com/user/repo.git`

- **steria status**
  - Show the current status of the repository
  - Example: `steria status`

- **steria send**
  - Copy the current directory to your Steria directory
  - Example: `steria send`

- **steria log**
  - Show commit history with color coding
  - Example: `steria log`

- **steria diff [file] [--side-by-side] [--context N]**
  - Show file differences (inline or side-by-side)
  - Example: `steria diff main.go --side-by-side`

- **steria restore <file> [commit-hash]**
  - Restore a file from a previous commit
  - Example: `steria restore main.go abc12345`

- **steria ignore [pattern]**
  - Manage .steriaignore file interactively or add a pattern
  - Example: `steria ignore *.log`

## Workflow Commands

- **steria commit "message" signer**
  - Create a manual commit
  - Example: `steria commit "Add feature" KleaSCM`

- **steria sync**
  - Synchronize with remote repository
  - Example: `steria sync`

- **steria done "message" signer**
  - Commit, sign, and sync everything automatically
  - Example: `steria done "Initial commit" KleaSCM`

## Project Management

- **steria projects add <name> <path>**
  - Add a new project
  - Example: `steria projects add my-project /path/to/project`

- **steria projects delete <name>**
  - Remove a project
  - Example: `steria projects delete my-project`

- **steria projects pull <name> <version> signer**
  - Pull a specific version of a project
  - Example: `steria projects pull my-project v1.0 KleaSCM`

## Branching System

- **steria add-branch <name>**
  - Create a new branch
  - Example: `steria add-branch feature-x`

- **steria branch <name>**
  - Switch to or create a branch
  - Example: `steria branch main`

- **steria delete-branch <name>**
  - Delete a branch
  - Example: `steria delete-branch feature-x`

- **steria merge <branch> signer**
  - Merge a branch into the current branch
  - Example: `steria merge feature-x KleaSCM`

- **steria rename-branch <old> <new>**
  - Rename a branch
  - Example: `steria rename-branch old-name new-name`

- **steria switch-branch <name>**
  - Switch to an existing branch
  - Example: `steria switch-branch feature-x`

---

For more details on each command, use `steria <command> --help`. 
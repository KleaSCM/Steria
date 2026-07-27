# Getting Started with Steria

**Author:** KleaSCM  
**Email:** KleaSCM@gmail.com

---

## Installation

Steria is written in Go. To build and install:

```sh
git clone https://github.com/your-username/steria.git
cd steria
go build -o steria
sudo cp steria /usr/local/bin/
```

## Initial Setup

1. Create or navigate to your project directory.
2. Initialize a Steria repository:
   ```sh
   steria done "Initial commit" KleaSCM
   ```
3. Start working! Use `steria status`, `steria commit`, and other commands as needed.

## Quick Start Example

```sh
# Create a new file
nano hello.txt

# Commit your changes
steria done "Add hello.txt" KleaSCM

# See your commit history
steria log

# See file differences
steria diff
```

For a full list of commands, see the [CLI Command Reference](cli.md). 
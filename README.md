# Monotask

A CLI tool to extract tasks directly from source files and markdown documents.

## Features

- Extracts TODO, BUG, NOTE markers (case insensitive) from C-style comments (`//` and `/* */`)
- Extracts TODO, BUG, NOTE markers (case insensitive) from shell script comments (`#`)
- Extracts TODO, BUG, NOTE markers (case insensitive) from Python comments and docstrings
- Extracts TODO, BUG, NOTE markers (case insensitive) from Lua comments
- Extracts unchecked checkboxes (`- [ ]`) from markdown files
- Supports optional assignee names in parentheses (e.g., `TODO(user): message`)
- Recursively scans directories
- Outputs in GNU Error Format for easy integration with other tools

## Installation

Install using `go install`:

```bash
go install github.com/IlyasYOY/monotask@latest
```

## Usage

```bash
# Scan current directory
./monotask

# Scan specific directory
./monotask /path/to/directory
```

## Output Format

```
file:line:column: type: message
```

Or with optional assignee:
```
file:line:column: type(assignee): message
```

```
work.c:15:3: TODO: this is todo marker in C code.
work.c:16:3: TODO(IlyasYOY): fix this bug.
tasks.md:14:12: CHECKBOX: this is not closed check-box.
```

Example:

```
➜  dotfiles git:(master) ✗ monotask .
/Users/IlyasYOY/Projects/IlyasYOY/dotfiles/config/nvim/after/ftplugin/go.lua:343:9: TODO: for now it works only for commands, I have to add the separate logic to support this in keymaps.
```

## Supported File Types

- `.c`, `.h` - C files (case insensitive TODO, BUG, NOTE markers in comments)
- `.java` - Java files (case insensitive TODO, BUG, NOTE markers in comments)
- `.go` - Go files (case insensitive TODO, BUG, NOTE markers in comments)
- `.js`, `.mjs` - JavaScript files (case insensitive TODO, BUG, NOTE markers in comments)
- `.ts`, `.mts` - TypeScript files (case insensitive TODO, BUG, NOTE markers in comments)
- `.cpp`, `.hpp`, `.cxx`, `.cc` - C++ files (case insensitive TODO, BUG, NOTE markers in comments)
- `.lua` - Lua files (case insensitive TODO, BUG, NOTE markers in comments)
- `.sh`, `.bash` - Shell scripts (case insensitive TODO, BUG, NOTE markers in comments)
- `.py` - Python files (case insensitive TODO, BUG, NOTE markers in # comments and single-line docstrings)
- `.md` - Markdown files (unchecked checkboxes)

Tasks can optionally include an assignee in parentheses after the type: `TODO(user): message`

## Ignoring Files and Directories

Monotask supports `.mtignore` files to exclude specific files or directories from scanning. Place a `.mtignore` file in any directory to list paths to ignore (one per line, relative to the `.mtignore` file's location).

- Ignores cascade from parent directories to subdirectories
- Child directories can add additional ignores with their own `.mtignore` files
- Only exact path matches are supported (no patterns or wildcards)

Example `.mtignore`:
```
build.log
node_modules/
temp.txt
```

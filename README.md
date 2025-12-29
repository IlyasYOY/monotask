# Monotask

A CLI tool to extract tasks directly from source files and markdown documents.

## Features

- Extracts TODO, BUG, NOTE markers from C code comments (`//` and `/* */`)
- Extracts unchecked checkboxes (`- [ ]`) from markdown files
- Recursively scans directories
- Outputs in GNU Error Format for easy integration with other tools

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

Example:
```
work.c:15:3: TODO: this is todo marker in C code.
tasks.md:14:12: CHECKBOX: this is not closed check-box.
```

## Supported File Types

- `.c`, `.h` - C files (TODO, BUG, NOTE markers in comments)
- `.md` - Markdown files (unchecked checkboxes)

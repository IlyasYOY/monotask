# Monotask

A CLI tool to extract tasks directly from source files and markdown documents.

## Features

- Extracts TODO, BUG, NOTE markers from C-style comments (`//` and `/* */`)
- Extracts TODO, BUG, NOTE markers from shell script comments (`#`)
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
- `.java` - Java files (TODO, BUG, NOTE markers in comments)
- `.go` - Go files (TODO, BUG, NOTE markers in comments)
- `.js`, `.mjs` - JavaScript files (TODO, BUG, NOTE markers in comments)
- `.ts`, `.mts` - TypeScript files (TODO, BUG, NOTE markers in comments)
- `.cpp`, `.hpp`, `.cxx`, `.cc` - C++ files (TODO, BUG, NOTE markers in comments)
- `.lua` - Lua files (TODO, BUG, NOTE markers in comments)
- `.sh`, `.bash` - Shell scripts (TODO, BUG, NOTE markers in comments)
- `.md` - Markdown files (unchecked checkboxes)

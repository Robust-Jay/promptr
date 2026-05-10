## 1. Project Scaffolding

- [x] 1.1 Initialize Go module (`go mod init promptr`)
- [x] 1.2 Create directory layout: `cmd/promptr/`, `internal/`, `pkg/`
- [x] 1.3 Create `cmd/promptr/main.go` with CLI entrypoint and command dispatch skeleton

## 2. Data Models and Storage

- [x] 2.1 Define `Prompt` struct with fields: id, title, category, content
- [x] 2.2 Implement YAML marshal/unmarshal for Prompt files
- [x] 2.3 Implement `~/.promptr/` directory resolution and creation
- [x] 2.4 Implement load all prompts from `builtin/` and `user/` (user overrides builtin on ID collision)

## 3. Built-in Prompt Library

- [x] 3.1 Create initial built-in prompts (code/writing/general categories) as YAML files
- [x] 3.2 Embed built-in prompts using `embed.FS`
- [x] 3.3 Implement first-run extraction: write embedded prompts to `~/.promptr/builtin/`

## 4. Core CLI Commands

- [x] 4.1 Implement `list` command with optional category filter
- [x] 4.2 Implement `search <query>` command with full-text matching across id/title/category/content
- [x] 4.3 Implement `show <id>` command displaying full prompt details
- [x] 4.4 Implement `categories` command listing all unique tags

## 5. Prompt CRUD

- [x] 5.1 Implement `add` command with interactive prompts (id, title, category, content)
- [x] 5.2 Implement `add` command with `--id`, `--title`, `--category`, `--content` flags
- [x] 5.3 Implement duplicate ID validation for add
- [x] 5.4 Implement `edit <id>` command (open in `$EDITOR`; copy builtin to user/ if editing builtin)
- [x] 5.5 Implement `rm <id>` command (delete only user prompts; reject builtin prompts)

## 6. Variable Filling

- [x] 6.1 Implement `{{variable}}` placeholder detection from prompt content
- [x] 6.2 Implement interactive variable resolution (prompt user for each variable)
- [x] 6.3 Implement placeholder replacement to produce resolved text

## 7. Clipboard Integration

- [x] 7.1 Implement clipboard abstraction interface (Windows/macOS/Linux backends)
- [x] 7.2 Wire clipboard copy into `cp <id>` command (after variable resolution if needed)
- [x] 7.3 Implement stdout fallback when clipboard access fails

## 8. CLI Polish and Help

- [x] 8.1 Add `--help` and `-h` support for each command
- [x] 8.2 Add command aliases (e.g., `cp` for copy, `ls` for list, `s` for search)
- [x] 8.3 Implement `promptr --version` output

## 9. Testing and Cross-Platform Verification

- [x] 9.1 Write unit tests for prompt loading, search, and variable resolution
- [x] 9.2 Write unit tests for clipboard abstraction
- [x] 9.3 Verify CLI works on all three platforms (Windows, macOS, Linux)

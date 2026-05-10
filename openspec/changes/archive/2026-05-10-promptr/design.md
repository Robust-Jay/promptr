## Context

`promptr` is a CLI tool for managing and reusing AI prompts. It runs as a single Go binary with no external dependencies. Prompts are stored as YAML files under `~/.promptr/`, with built-in prompts embedded in the binary and user prompts separated for overrides and additions.

## Goals / Non-Goals

**Goals:**
- Search, list, show, add, edit, and remove prompts via CLI
- Interactive `{{variable}}` filling with one-click clipboard copy
- Pre-installed prompt library embedded in the binary
- User-customizable prompts that override built-ins (copy-modify model)
- Cross-platform support (Windows, macOS, Linux)

**Non-Goals:**
- Cloud sync or multi-device sharing
- GUI or web interface
- AI platform integration (the tool is intentionally platform-agnostic)
- Prompt execution or AI API calls

## Decisions

### Go project layout
Standard layout: `cmd/promptr/` (CLI entrypoint), `internal/` (non-public packages), `pkg/` (reusable packages if any). The `main.go` dispatches to commands. Cobra or a minimal flag-based approach — TBD during implementation.

### Storage: YAML files under ~/.promptr/
Each prompt is one YAML file. Chosen over SQLite or a single JSON blob because:
- File-per-prompt makes human editing, git versioning, and sharing trivial
- YAML is readable and supports multi-line content naturally
- No external database dependency

Directory structure:
```
~/.promptr/
├── builtin/          ← extracted from binary on first run
│   ├── code/
│   │   ├── unit-test-jest.yaml
│   │   └── code-review.yaml
│   ├── writing/
│   └── general/
└── user/             ← user-created and modified prompts
    └── my-custom.yaml
```

### Prompt file format
```yaml
id: unit-test-jest
title: 编写单元测试 (Jest)
category: [code, testing]
content: |
  请为以下 {{language}} 代码编写完整的单元测试，使用 {{framework}} 框架。
```

### Search: full-field text matching
Search scans `id`, `title`, `category`, and `content` fields. Results are ranked by match relevance. No external search index — for the expected prompt count (hundreds, not millions), in-memory scanning on demand is sufficient.

### Variable filling: interactive prompt
When a prompt contains `{{variable}}` placeholders, `promptr cp` extracts each unique variable name and prompts the user to enter a value. After all variables are filled, the resolved text is copied to clipboard.

### Edit: copy-modify for builtins
Editing a built-in prompt copies it from `builtin/` to `user/` first. The user edits the copy. The built-in version remains untouched, allowing future updates to the built-in library without clobbering user changes.

### Clipboard: OS-appropriate approach
Uses Go's `os/exec` to call platform-native clipboard tools (`clip` on Windows, `pbcopy` on macOS, `xclip`/`xsel` on Linux). Could use a Go clipboard library if a suitable cross-platform one exists.

### CLI framework: minimal, likely manual flag parsing or a tiny library
Given the command surface is small (<10 subcommands), a full framework like Cobra may be overkill but provides good help text generation. Final choice during implementation.

## Risks / Trade-offs

- **[Risk] YAML hand-editing by users leads to syntax errors** → `promptr edit` opens in `$EDITOR` with pre-validated structure; add command validates input
- **[Risk] Clipboard access differs across platforms** → Abstract behind an interface; test on all three platforms
- **[Risk] Growing prompt library slows search linearly** → Acceptable for the expected scale (<1000 prompts); can add indexing later if needed
- **[Trade-off] File-per-prompt vs single file** → Chose file-per-prompt for git-friendliness and human readability, at the cost of slightly slower listing

# AGENTS.md

## Repository Overview

This repo contains the **promptr** CLI tool (Go) and uses the **OpenSpec workflow** for change management. There are two layers: application code and spec-driven change tracking.

## promptr (Go project)

- **Module**: `promptr` (Go 1.26, single dependency: `gopkg.in/yaml.v3`)
- **Entrypoint**: `cmd/promptr/main.go`
- **Internal packages**: `internal/prompt/` (data model, YAML I/O, search, variable filling), `internal/clipboard/` (OS-specific clipboard), `internal/builtin/` (embedded prompts via `embed.FS`)
- `pkg/` is empty — no public API yet

### Go proxy quirk

The default Go proxy (`proxy.golang.org`) is unreachable from this environment. Always set before any `go` command:
```powershell
$env:GOPROXY="direct"; $env:GOINSECURE="*"; $env:GONOSUMCHECK="*"; $env:GONOSUMDB="*"
```

### Commands

```powershell
# Build
go build -o promptr.exe ./cmd/promptr/

# Run all tests
go test ./...

# Run a single package's tests
go test ./internal/prompt/ -v
go test ./internal/clipboard/ -v
```

### No root .gitignore

`promptr.exe` and `go.sum` are in the repo root with no `.gitignore`. Add one if you want to exclude build artifacts.

## OpenSpec Workflow

- **Schema**: `spec-driven` (defined in `openspec/config.yaml`)
- **CLI**: `openspec` for change management (`openspec list`, `openspec status`, `openspec new change`, etc.)
- **Changes**: `openspec/changes/<name>/` — artifacts like `proposal.md`, `design.md`, `tasks.md`, `specs/`
- **Archived**: `openspec/changes/archive/<YYYY-MM-DD-<name>/`

### Custom Slash Commands

| Command | Purpose |
|---|---|
| `/opsx-explore` | Think through ideas, investigate problems (no implementation) |
| `/opsx-propose` | Create a new change with all artifacts |
| `/opsx-apply` | Implement tasks from an existing change |
| `/opsx-archive` | Archive a completed change (moves to archive/) |

Definitions: `.opencode/commands/opsx-*.md`. Skills: `.opencode/skills/openspec-*/SKILL.md`.

## Conventions

- Change names use **kebab-case**
- Artifact creation follows dependency order from `openspec status --json`
- Use `openspec instructions apply --change "<name>" --json` to get context files and task list
- Always read context files before implementing tasks
- Mark tasks done: `- [ ]` → `- [x]` in `tasks.md`

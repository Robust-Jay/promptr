## Why

Developers spend significant time crafting and re-crafting AI prompts. A fast, platform-agnostic CLI tool lets them search, reuse, and customize prompts without leaving the terminal — no vendor lock-in, no context switching.

## What Changes

- New CLI tool `promptr` written in Go, distributed as a single cross-platform binary
- YAML-based prompt storage under `~/.promptr/` with built-in vs. user overridable prompts
- Interactive variable/placeholder filling before clipboard copy
- Flat tag-based categorization and full-text search across prompts

## Capabilities

### New Capabilities
- `prompt-management`: CRUD operations on prompts — search, list, show, add, edit, remove — with tag-based categorization and full-text search
- `builtin-library`: Pre-installed prompt library embedded in the binary, auto-unpacked to `~/.promptr/builtin/` on first run
- `variable-filling`: Interactive `{{placeholder}}` resolution — prompts containing variables prompt the user for values before copying to clipboard
- `clipboard-copy`: Platform-agnostic clipboard output for prompts after variable resolution

### Modified Capabilities
<!-- No existing capabilities to modify -->

## Impact

- New Go project at repository root (standard layout: `cmd/promptr/`, `internal/`, `pkg/`)
- Zero dependencies on external APIs or platforms
- No existing code affected — this is a greenfield addition

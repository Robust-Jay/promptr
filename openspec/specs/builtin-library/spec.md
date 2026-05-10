## ADDED Requirements

### Requirement: Embed built-in prompts in binary
The system SHALL embed a set of pre-installed prompts within the Go binary using embed.FS.

#### Scenario: Binary contains built-in prompts
- **WHEN** the binary is built
- **THEN** the embedded prompt library is accessible via Go's embed package

### Requirement: Extract built-in prompts on first run
The system SHALL extract embedded prompts to `~/.promptr/builtin/` when the directory does not exist.

#### Scenario: First run initializes builtin directory
- **WHEN** `~/.promptr/builtin/` does not exist and any command is run
- **THEN** the embedded prompts are written to `~/.promptr/builtin/` with the same directory structure as embedded

#### Scenario: Subsequent runs do not overwrite
- **WHEN** `~/.promptr/builtin/` already exists
- **THEN** the existing builtin prompts are left unchanged

### Requirement: Built-in prompt categories
The system SHALL provide built-in prompts organized in at least the following categories: code, writing, general.

#### Scenario: Built-in prompts are available
- **WHEN** user runs `promptr list` after first run
- **THEN** prompts from code, writing, and general categories are listed

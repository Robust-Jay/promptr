## ADDED Requirements

### Requirement: Search prompts
The system SHALL search all prompts by matching user query text against id, title, category, and content fields.

#### Scenario: Search finds matching prompts
- **WHEN** user runs `promptr search "unit test"`
- **THEN** all prompts where id, title, category, or content contains "unit test" are displayed with id, title, and category

#### Scenario: Search returns empty results
- **WHEN** user runs `promptr search "xyznotfound"`
- **THEN** system displays "No prompts found"

### Requirement: List all prompts
The system SHALL list all available prompts with id, title, and category.

#### Scenario: List all prompts
- **WHEN** user runs `promptr list`
- **THEN** all prompts from both builtin/ and user/ directories are listed

#### Scenario: Filter by category
- **WHEN** user runs `promptr list code/testing`
- **THEN** only prompts tagged with both "code" and "testing" are listed

### Requirement: Show prompt details
The system SHALL display the full content of a specific prompt.

#### Scenario: Show existing prompt
- **WHEN** user runs `promptr show unit-test-jest`
- **THEN** the prompt's id, title, category, and full content are displayed

#### Scenario: Show non-existent prompt
- **WHEN** user runs `promptr show nonexistent-id`
- **THEN** system displays "Prompt not found: nonexistent-id"

### Requirement: Add a new prompt
The system SHALL allow users to create a new prompt interactively or via command flags.

#### Scenario: Interactive add
- **WHEN** user runs `promptr add`
- **THEN** system prompts for id, title, category, and content sequentially, then saves to `~/.promptr/user/<id>.yaml`

#### Scenario: Flag-based add
- **WHEN** user runs `promptr add --id my-prompt --title "My Prompt" --category code --content "Review this {{language}} code"`
- **THEN** system saves the prompt directly without interactive prompts

#### Scenario: Duplicate ID rejected
- **WHEN** user tries to add a prompt with an existing ID
- **THEN** system displays an error and does not overwrite

### Requirement: Edit a prompt
The system SHALL open an existing prompt for editing. For built-in prompts, the prompt SHALL be copied to user/ before editing.

#### Scenario: Edit user prompt
- **WHEN** user runs `promptr edit my-prompt`
- **THEN** the prompt file at `~/.promptr/user/my-prompt.yaml` is opened in `$EDITOR`

#### Scenario: Edit builtin prompt triggers copy
- **WHEN** user runs `promptr edit unit-test-jest` and the prompt exists only in builtin/
- **THEN** the prompt is copied to `~/.promptr/user/unit-test-jest.yaml` and opened in `$EDITOR`

#### Scenario: Edit non-existent prompt
- **WHEN** user runs `promptr edit nonexistent-id`
- **THEN** system displays "Prompt not found: nonexistent-id"

### Requirement: Remove a prompt
The system SHALL delete a user-created prompt. Built-in prompts SHALL NOT be deletable.

#### Scenario: Remove user prompt
- **WHEN** user runs `promptr rm my-prompt` and the prompt exists in user/
- **THEN** the prompt file is deleted and a confirmation message is displayed

#### Scenario: Cannot remove builtin prompt
- **WHEN** user runs `promptr rm unit-test-jest` and the prompt exists only in builtin/
- **THEN** system displays "Cannot remove built-in prompt. Use a user prompt instead."

### Requirement: List categories
The system SHALL list all unique category tags across all prompts.

#### Scenario: List categories
- **WHEN** user runs `promptr categories`
- **THEN** all unique tags from builtin/ and user/ prompts are displayed

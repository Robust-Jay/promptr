## ADDED Requirements

### Requirement: Detect placeholders in prompt content
The system SHALL detect all `{{variable_name}}` patterns in a prompt's content.

#### Scenario: Extract variables from prompt
- **WHEN** a prompt contains "Write {{language}} tests using {{framework}}"
- **THEN** the system identifies variables "language" and "framework"

#### Scenario: No variables in prompt
- **WHEN** a prompt contains no `{{...}}` patterns
- **THEN** the system proceeds directly to clipboard copy without prompting

### Requirement: Interactive variable resolution
The system SHALL prompt users to enter a value for each detected variable before copying.

#### Scenario: Fill variables interactively
- **WHEN** user runs `promptr cp unit-test-jest` and the prompt has variables "language" and "framework"
- **THEN** system prompts "language:" and "framework:" sequentially, then assembles the resolved text

#### Scenario: Empty value allowed
- **WHEN** user enters an empty value for a variable
- **THEN** the placeholder is replaced with an empty string

### Requirement: Resolved text replaces placeholders
The system SHALL replace each `{{variable}}` with the user-provided value in the final output.

#### Scenario: Variable replacement
- **WHEN** prompt content is "Review {{language}} code" and user enters "Python"
- **THEN** resolved text is "Review Python code"

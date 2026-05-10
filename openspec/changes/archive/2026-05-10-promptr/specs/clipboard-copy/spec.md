## ADDED Requirements

### Requirement: Copy resolved prompt to clipboard
The system SHALL copy the resolved prompt text to the system clipboard.

#### Scenario: Copy after variable resolution
- **WHEN** user completes variable filling for `promptr cp <id>`
- **THEN** the resolved text is placed on the system clipboard and a confirmation message is displayed

#### Scenario: Copy without variables
- **WHEN** user runs `promptr cp <id>` for a prompt with no variables
- **THEN** the prompt content is copied directly to clipboard with a confirmation message

#### Scenario: Clipboard failure reported
- **WHEN** clipboard access fails (no clipboard tool available)
- **THEN** system displays an error message and prints the resolved text to stdout as fallback

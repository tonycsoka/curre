# TUI Workflow

A JSON-driven terminal UI for running sequenced, parameterised shell workflows.

## Features

- **JSON-driven workflows**: Define steps, parameters, and scripts in a simple JSON file.
- **Interactive parameter input**: Edit parameters in the TUI before running each step.
- **Sequential execution**: Steps unlock only after the previous step succeeds (or is bypassed).
- **Session persistence**: Sessions auto-save and are directory-aware. Resume where you left off.
- **Live output**: Stream stdout/stderr from scripts in real-time.
- **Bypass support**: If a step fails, manually bypass it with confirmation to unlock the next step.
- **Run once per session**: Mark steps that should only execute once per session.

## Installation

```bash
go build .
```

## Usage

```bash
./tui-workflow <workflow.json>
```

Example:

```bash
./tui-workflow examples/deploy.json
```

## Workflow JSON Format

```json
{
  "name": "deploy",
  "description": "Deploy the application",
  "parameters": {
    "env": {
      "type": "string",
      "default": "dev",
      "description": "Target environment"
    },
    "version": {
      "type": "string",
      "default": "1.0.0",
      "description": "App version"
    }
  },
  "steps": [
    {
      "id": "build",
      "name": "Build",
      "script": "scripts/build.sh",
      "params": ["env", "version"],
      "run_once_per_session": false,
      "description": "Build the application"
    },
    {
      "id": "deploy",
      "name": "Deploy",
      "script": "scripts/deploy.sh",
      "params": ["env", "version"],
      "run_once_per_session": false,
      "description": "Deploy the application"
    }
  ]
}
```

### Field Reference

- `name` (string, required): Workflow name.
- `parameters` (object): Global parameters available to all steps.
  - `type`: Parameter type (`string`).
  - `default`: Default value.
  - `description`: Human-readable description.
- `steps` (array, required):
  - `id`: Unique step identifier.
  - `name`: Display name.
  - `script`: Path to shell script (relative to workflow JSON or absolute).
  - `params`: Array of parameter names to pass as positional arguments to the script.
  - `run_once_per_session`: If `true`, the step is skipped if it already succeeded in the current session.
  - `description`: Optional description.

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate steps |
| `Tab` | Focus parameter inputs |
| `Shift+Tab` | Previous parameter input |
| `Esc` | Unfocus parameters / close modals |
| `r` | Run selected step |
| `b` | Bypass failed step (with confirmation) |
| `n` | Skip step with `run_once_per_session` |
| `s` | Show sessions for current directory |
| `q` / `Ctrl+C` | Quit |

## Session Storage

Sessions are automatically saved to:

```
~/.local/share/tui-workflow/sessions/<workflow-name>-<cwd-hash>.json
```

Sessions are directory-aware. Running the same workflow from different directories creates separate sessions.

## Development

```bash
go test -v
```

## License

MIT

# TUI Workflow System — Plan

## Context
Build a JSON-driven TUI application that lets users run shell-script workflows step-by-step. Each step is parameterised, and progression is gated on success of the previous step. Sessions auto-save, are directory-aware, and support per-step "run once per session" flags.

## Decisions

| Decision | Choice |
|----------|--------|
| **Language** | Go + Bubble Tea (single binary, Elm architecture, excellent for CLI-centric workflows) |
| **Parameter passing** | Workflow JSON defines global parameters. TUI shows an editable form for all parameters before a step runs. Values are passed as **positional arguments** to the shell script in the order listed in the step's `params` array. |
| **Success criteria** | Exit code 0 = success. Anything else = failure, blocking the next step. **Manual bypass available**: user can confirm marking a failed step as bypassed, which unlocks the next step. |
| **Execution model** | Strictly sequential. A step can only run if the previous step succeeded, was skipped, or was bypassed. |
| **Session storage** | Auto-saved to `~/.local/share/tui-workflow/sessions/<workflow-name>-<cwd-hash>.json`. Directory-aware: session key is derived from the absolute path of the working directory. |
| "run once per session" | Per-step boolean flag in JSON. If a step succeeded in the current session and flag is true, it is skipped on subsequent workflow runs. |

## JSON Schema

### Workflow Definition
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
      "description": "Build the Docker image"
    },
    {
      "id": "deploy",
      "name": "Deploy",
      "script": "scripts/deploy.sh",
      "params": ["env", "version"],
      "run_once_per_session": false
    }
  ]
}
```

### Session State
```json
{
  "workflow_name": "deploy",
  "cwd": "/home/user/project",
  "parameter_values": {
    "env": "staging",
    "version": "2.0.0"
  },
  "step_states": {
    "build": {
      "status": "success",
      "exit_code": 0,
      "run_at": "2026-01-15T10:30:00Z",
      "output": "..."
    },
    "deploy": {
      "status": "pending"
    }
  }
}
```

Status values: `pending`, `running`, `success`, `failed`, `skipped`, `bypassed`.

## File Structure

```
.
├── main.go              # Entry point: go run . workflow.json
├── go.mod
├── app.go               # Bubble Tea Model, Init, Update, View
├── workflow.go          # Workflow JSON parsing / validation
├── session.go           # Session load/save/auto-save
├── runner.go            # Shell command execution with live output (tea.Cmd)
├── widgets.go           # Custom Bubble Tea components
│   (step list, parameter form, output viewport, bypass modal)
├── examples/
│   ├── deploy.json      # Example workflow definition
│   └── scripts/
│       ├── build.sh     # Dummy build script (echoes args, exits 0)
│       └── deploy.sh    # Dummy deploy script (echoes args, exits 0)
└── README.md
```

### Example Workflow (`examples/deploy.json`)
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
      "script": "examples/scripts/build.sh",
      "params": ["env", "version"],
      "run_once_per_session": false,
      "description": "Build the application"
    },
    {
      "id": "deploy",
      "name": "Deploy",
      "script": "examples/scripts/deploy.sh",
      "params": ["env", "version"],
      "run_once_per_session": false,
      "description": "Deploy the application"
    }
  ]
}
```

### Example Scripts

**`examples/scripts/build.sh`**
```bash
#!/usr/bin/env bash
set -euo pipefail
echo "[BUILD] Starting build for env=$1 version=$2"
echo "[BUILD] Compiling..."
echo "[BUILD] Done."
```

**`examples/scripts/deploy.sh`**
```bash
#!/usr/bin/env bash
set -euo pipefail
echo "[DEPLOY] Starting deploy for env=$1 version=$2"
echo "[DEPLOY] Pushing to server..."
echo "[DEPLOY] Done."
```

## Reuse
- **Bubble Tea** core (`tea.Model`, `tea.Cmd`, `tea.Msg`) for the app loop and state management.
- **Bubbles** library (`list`, `viewport`, `textinput`, `textarea`, `spinner`, `help`) for reusable TUI components.
- **Lipgloss** for styling and layout.
- **Standard library**: `os/exec` for shell commands with live output via `io.Pipe`, `encoding/json`, `path/filepath`, `crypto/sha256` (for directory hashing), `time`.

## Steps

- [ ] **Step 1 — Scaffold**: `go mod init`, import `github.com/charmbracelet/bubbletea`, `bubbles`, `lipgloss`. Create `main.go` entry point.
- [ ] **Step 2 — Workflow Parser**: Define Go structs for Workflow, Parameter, Step. Load and validate JSON with `encoding/json`.
- [ ] **Step 3 — Session Manager**: Implement load/save with directory-aware path resolution (`~/.local/share/tui-workflow/sessions/`). Auto-save on every state change.
- [ ] **Step 4 — TUI Skeleton**: Build Bubble Tea model with three-pane layout: `StepList` (left), `ParamForm` + `OutputViewport` (right, stacked), plus `Help` footer.
- [ ] **Step 5 — Parameter Form**: Use `textinput` bubbles to render inputs for all parameters. Show defaults. Validate before run.
- [ ] **Step 6 — Shell Runner**: Implement a `tea.Cmd` that runs `os/exec` asynchronously via goroutines. Stream stdout/stderr lines as messages. Capture exit code.
- [ ] **Step 7 — Step Logic + Bypass**: Implement sequential gating. Only enable "Run" for the next pending step if previous step succeeded or was bypassed. When a step fails, show a "Bypass" action. Trigger a confirmation modal; on confirm, set status to `bypassed` and unlock the next step. Handle `run_once_per_session` by checking session state and offering "Skip".
- [ ] **Step 8 — Session Switching**: Add a keybind to list existing sessions for the current directory and switch between them.
- [ ] **Step 9 — Polish**: Key bindings (r run, b bypass, n next, q quit, ? help), status colors (lipgloss), scroll lock on output viewport, window size responsiveness.

## Verification

1. Create a sample workflow JSON with two steps and two parameters.
2. Create dummy `scripts/build.sh` and `scripts/deploy.sh` that echo arguments and exit 0.
3. Run `go run . sample.json`.
4. Verify:
   - Parameters appear with defaults.
   - Changing parameters and running step 1 passes them as positional args.
   - Live output is shown in the viewport.
   - Step 2 is disabled until step 1 succeeds.
   - Exit 1 on step 1 blocks step 2.
   - Pressing `b` on a failed step opens a confirmation modal; confirming marks it `bypassed` and unlocks step 2.
   - `run_once_per_session: true` causes step to be skipped on second run.
   - Closing and reopening the app restores the session state.
   - Session file exists in `~/.local/share/tui-workflow/sessions/`. 


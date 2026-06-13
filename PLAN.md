# Session System Redesign Plan

## Context

The current session model is a single session per `(workflowName, cwd)` pair. The user wants:

1. **Auto-load on startup:**
   - No previous session → create a new session with unique name
   - Previous session has all tasks done (success/skipped) → create a new session
   - Previous session has pending tasks → resume that session

2. **Session modal:**
   - Show each session with status (done, in progress, failed, pending)
   - Allow loading any existing session into the main UI for viewing/inspection
   - Allow starting a new session

## Approach

Change the session model to support multiple named sessions per workflow per directory.

### Session struct changes

Add `Name` and `CreatedAt` fields:

```go
type Session struct {
    Name            string                 `json:"name"`
    WorkflowName    string                 `json:"workflow_name"`
    Cwd             string                 `json:"cwd"`
    CreatedAt       string                 `json:"created_at"`
    ParameterValues map[string]string      `json:"parameter_values"`
    StepStates      map[string]StepState   `json:"step_states"`
}
```

### File naming

New format: `workflowName-<hash>-<sessionName>.json`

Example: `deploy-abc123de-run-1.json`

### New functions

- `NewSession(wf, cwd, name)` — creates a session with a given name
- `GenerateSessionName(workflowName, cwd)` — generates unique `run-N` name
- `SessionPath(workflowName, cwd, name)` — returns file path
- `FindSessionsForWorkflow(workflowName, cwd)` — returns all sessions for this workflow+dir, sorted by CreatedAt desc
- `GetLatestSession(workflowName, cwd)` — returns most recent session
- `LoadSessionByName(workflowName, cwd, name)` — loads a specific session
- `Session.OverallStatus()` — returns `done`, `failed`, `running`, `pending`, or `in progress`

### Auto-load logic (main.go)

```go
sessions, _ := FindSessionsForWorkflow(wf.Name, cwd)
if len(sessions) == 0 {
    session = NewSession(wf, cwd, GenerateSessionName(wf.Name, cwd))
} else {
    latest := sessions[0]
    if latest.OverallStatus() == "done" {
        session = NewSession(wf, cwd, GenerateSessionName(wf.Name, cwd))
    } else {
        session = latest
    }
}
```

### UI changes (app.go)

- Change `sessionList` from `[]string` to `[]*Session`
- `renderSessionList` shows name + status with color coding
- Status styles: `done` (green), `failed` (red), `running` (yellow), `pending` (gray), `in progress` (default)
- Enter loads the selected session
- `n` creates a new session with auto-generated name

### Key handler updates

- `s` key: load all sessions for this workflow, show modal
- `n` in modal: create new session with unique name
- `enter` in modal: load selected session
- `up`/`down` in modal: navigate sessions

## Files to modify

- `session.go` — core session model changes
- `main.go` — startup auto-load logic
- `app.go` — session picker UI and key handlers

## Reuse

- Existing `Session` struct and `StepStatus` types
- Existing `lipgloss` styles for status color coding
- Existing `SaveSession` and `LoadSessionFromPath` (updated)

## Steps

- [ ] Update `Session` struct with `Name` and `CreatedAt`
- [ ] Update `NewSession` to accept name
- [ ] Update `SessionPath` to include name
- [ ] Add `GenerateSessionName`
- [ ] Add `FindSessionsForWorkflow`
- [ ] Add `GetLatestSession`
- [ ] Add `LoadSessionByName`
- [ ] Add `Session.OverallStatus()`
- [ ] Update `SaveSession` to use new path
- [ ] Update `main.go` auto-load logic
- [ ] Update `app.go` model to use `[]*Session`
- [ ] Update `renderSessionList` to show status
- [ ] Update key handlers for session picker
- [ ] Update `FindSessionsForDir` (or replace with `FindSessionsForWorkflow`)
- [ ] Test: build and verify

## Verification

Build the app, run with a workflow, verify:
1. First run creates a new session
2. After all steps done, next run creates a new session
3. With pending steps, same session is resumed
4. Session picker shows all sessions with status
5. Can load an old session from picker
6. Can create new session from picker

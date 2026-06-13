package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// StepStatus represents the execution state of a step.
type StepStatus string

const (
	StatusPending   StepStatus = "pending"
	StatusRunning   StepStatus = "running"
	StatusSuccess   StepStatus = "success"
	StatusFailed    StepStatus = "failed"
	StatusSkipped   StepStatus = "skipped"
)

// StepState tracks the execution state and output of a single step.
type StepState struct {
	Status   StepStatus `json:"status"`
	ExitCode int        `json:"exit_code,omitempty"`
	RunAt    string     `json:"run_at,omitempty"`
	Output   string     `json:"output,omitempty"`
}

// Session is the persisted state for a workflow run in a specific directory.
type Session struct {
	WorkflowName    string                 `json:"workflow_name"`
	Cwd             string                 `json:"cwd"`
	ParameterValues map[string]string      `json:"parameter_values"`
	StepStates      map[string]StepState   `json:"step_states"`
}

// NewSession creates a fresh session for the given workflow and directory.
func NewSession(wf *Workflow, cwd string) *Session {
	stepStates := make(map[string]StepState)
	for _, step := range wf.Steps {
		stepStates[step.ID] = StepState{Status: StatusPending}
	}

	return &Session{
		WorkflowName:    wf.Name,
		Cwd:             cwd,
		ParameterValues: make(map[string]string),
		StepStates:      stepStates,
	}
}

// SessionDir returns the directory where session files are stored.
func SessionDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".tui-workflow/sessions"
	}
	return filepath.Join(home, ".local", "share", "tui-workflow", "sessions")
}

// SessionPath returns the file path for a session based on workflow name and cwd.
func SessionPath(workflowName, cwd string) string {
	hash := sha256.Sum256([]byte(cwd))
	slug := fmt.Sprintf("%s-%s.json", workflowName, hex.EncodeToString(hash[:8]))
	return filepath.Join(SessionDir(), slug)
}

// LoadSession reads a session from disk, or returns nil if not found.
func LoadSession(workflowName, cwd string) (*Session, error) {
	path := SessionPath(workflowName, cwd)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading session file: %w", err)
	}

	var sess Session
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("parsing session JSON: %w", err)
	}
	if sess.StepStates == nil {
		sess.StepStates = make(map[string]StepState)
	}
	if sess.ParameterValues == nil {
		sess.ParameterValues = make(map[string]string)
	}
	return &sess, nil
}

// SaveSession writes the session to disk, creating directories if needed.
func SaveSession(sess *Session) error {
	path := SessionPath(sess.WorkflowName, sess.Cwd)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating session directory: %w", err)
	}

	data, err := json.MarshalIndent(sess, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling session: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing session file: %w", err)
	}
	return nil
}

// FindSessionsForDir returns all session files associated with a given directory.
func FindSessionsForDir(cwd string) ([]string, error) {
	sessDir := SessionDir()
	entries, err := os.ReadDir(sessDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var results []string
	hash := sha256.Sum256([]byte(cwd))
	targetHash := hex.EncodeToString(hash[:8])
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && filepath.Ext(name) == ".json" {
			// Check if filename ends with the cwd hash
			if len(name) > len(targetHash)+5 && name[len(name)-len(targetHash)-5:len(name)-5] == "-"+targetHash {
				results = append(results, name)
			}
		}
	}
	return results, nil
}

// UpdateStepState updates a step's state and auto-saves the session.
func (sess *Session) UpdateStepState(stepID string, state StepState) {
	if sess.StepStates == nil {
		sess.StepStates = make(map[string]StepState)
	}
	state.RunAt = time.Now().Format(time.RFC3339)
	sess.StepStates[stepID] = state
}

// SetParameterValue sets a parameter value and auto-saves.
func (sess *Session) SetParameterValue(key, value string) {
	if sess.ParameterValues == nil {
		sess.ParameterValues = make(map[string]string)
	}
	sess.ParameterValues[key] = value
}

// GetParameterValue returns the parameter value, or the default from the workflow.
func (sess *Session) GetParameterValue(key string, wf *Workflow) string {
	if val, ok := sess.ParameterValues[key]; ok {
		return val
	}
	if def, ok := wf.Parameters[key]; ok {
		return def.Default
	}
	return ""
}

// IsStepRunnable checks whether a step is eligible to run based on sequence and run_once.
func (sess *Session) IsStepRunnable(wf *Workflow, idx int) bool {
	if idx < 0 || idx >= len(wf.Steps) {
		return false
	}
	step := wf.Steps[idx]
	state := sess.StepStates[step.ID]

	// If it's already running, don't run again.
	if state.Status == StatusRunning {
		return false
	}

	// If run_once_per_session and already succeeded, skip.
	if step.RunOncePerSession && state.Status == StatusSuccess {
		return false
	}

	// If run_once_per_session and already skipped, skip.
	if step.RunOncePerSession && state.Status == StatusSkipped {
		return false
	}

	// First step is always runnable if not already running/success.
	if idx == 0 {
		return state.Status == StatusPending || state.Status == StatusFailed || state.Status == StatusSkipped
	}

	// Previous step must be success or skipped.
	prevStep := wf.Steps[idx-1]
	prevState := sess.StepStates[prevStep.ID]
	return prevState.Status == StatusSuccess || prevState.Status == StatusSkipped
}

// IsStepBypassable checks whether a failed step can be skipped.
func (sess *Session) IsStepBypassable(wf *Workflow, idx int) bool {
	if idx < 0 || idx >= len(wf.Steps) {
		return false
	}
	step := wf.Steps[idx]
	state := sess.StepStates[step.ID]
	return state.Status == StatusFailed
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: curre <workflow.json>")
		os.Exit(1)
	}

	workflowPath := os.Args[1]
	wf, err := LoadWorkflow(workflowPath)
	if err != nil {
		fmt.Printf("Error loading workflow: %v\n", err)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %v\n", err)
		os.Exit(1)
	}

	workflowDir := filepath.Dir(workflowPath)
	if !filepath.IsAbs(workflowDir) {
		workflowDir = filepath.Join(cwd, workflowDir)
	}

	// Load existing session or create new
	sessions, err := FindSessionsForWorkflow(wf.Name, cwd)
	if err != nil {
		fmt.Printf("Error finding sessions: %v\n", err)
		os.Exit(1)
	}

	var session *Session
	if len(sessions) == 0 {
		session = NewSession(wf, cwd)
	} else {
		latest := sessions[0]
		if latest.OverallStatus() == "done" {
			session = NewSession(wf, cwd)
		} else {
			session = latest
		}
	}

	m := initialModel(wf, session, workflowDir)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

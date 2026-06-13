package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tui-workflow <workflow.json>")
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
	session, err := LoadSession(wf.Name, cwd)
	if err != nil {
		fmt.Printf("Error loading session: %v\n", err)
		os.Exit(1)
	}
	if session == nil {
		session = NewSession(wf, cwd)
	}

	m := initialModel(wf, session, workflowDir)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

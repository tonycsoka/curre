package main

import (
	"bufio"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type shellStdoutMsg struct {
	line   string
	stepID string
}

type shellStderrMsg struct {
	line   string
	stepID string
}

type shellDoneMsg struct {
	stepID   string
	exitCode int
	status   StepStatus
}

type stepRunner struct {
	stdoutChan chan string
	stderrChan chan string
	resultChan chan shellDoneMsg
	stepID     string
}

func newStepRunner(step Step, workflowDir string, params []string) *stepRunner {
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	resultChan := make(chan shellDoneMsg)

	go func() {
		scriptPath := ResolveScriptPath(workflowDir, step.Script)
		cmd := exec.Command(scriptPath, params...)
		cmd.Dir = workflowDir

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			resultChan <- shellDoneMsg{stepID: step.ID, exitCode: -1, status: StatusFailed}
			return
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			resultChan <- shellDoneMsg{stepID: step.ID, exitCode: -1, status: StatusFailed}
			return
		}

		if err := cmd.Start(); err != nil {
			resultChan <- shellDoneMsg{stepID: step.ID, exitCode: -1, status: StatusFailed}
			return
		}

		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				stdoutChan <- scanner.Text() + "\n"
			}
		}()

		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				stderrChan <- scanner.Text() + "\n"
			}
		}()

		if err := cmd.Wait(); err != nil {
			exitCode := -1
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			}
			resultChan <- shellDoneMsg{stepID: step.ID, exitCode: exitCode, status: StatusFailed}
		} else {
			resultChan <- shellDoneMsg{stepID: step.ID, exitCode: 0, status: StatusSuccess}
		}
	}()

	return &stepRunner{
		stdoutChan: stdoutChan,
		stderrChan: stderrChan,
		resultChan: resultChan,
		stepID:     step.ID,
	}
}

func (r *stepRunner) NextCmd() tea.Cmd {
	if r == nil {
		return nil
	}
	return func() tea.Msg {
		select {
		case line := <-r.stdoutChan:
			return shellStdoutMsg{line: line, stepID: r.stepID}
		case line := <-r.stderrChan:
			return shellStderrMsg{line: line, stepID: r.stepID}
		case result := <-r.resultChan:
			return result
		}
	}
}

func buildParams(step Step, m *model) []string {
	var params []string
	for _, name := range step.Params {
		val := m.session.GetParameterValue(name, m.workflow)
		params = append(params, val)
	}
	return params
}

package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCLIUsage(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping CLI integration test on windows")
	}

	cmd := exec.Command("go", "run", ".")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected error when run with no args")
	}
	if !strings.Contains(string(out), "Usage:") {
		t.Fatalf("expected usage message, got: %s", out)
	}
}

func TestCLIBadWorkflowPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping CLI integration test on windows")
	}

	cmd := exec.Command("go", "run", ".", "nonexistent.json")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatal("expected error when workflow file is missing")
	}
	if !strings.Contains(string(out), "Error loading workflow") {
		t.Fatalf("expected load error, got: %s", out)
	}
}

func TestCLILoadsExampleWorkflow(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping CLI integration test on windows")
	}

	for _, path := range []string{
		"examples/deploy.json",
		"examples/full-demo.json",
		"examples/markdown.json",
		"examples/parallel.json",
	} {
		absPath, err := filepath.Abs(path)
		if err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("go", "run", ".", absPath)
		// Use a non-interactive mode so bubbletea doesn't require a TTY.
		// Bubbletea v2 exits with an error when there's no TTY, but we can
		// verify that the workflow at least loads successfully before the TTY error.
		out, err := cmd.CombinedOutput()
		outStr := string(out)

		// We expect a TTY error on exit, but we should NOT see a workflow load error.
		if strings.Contains(outStr, "Error loading workflow") {
			t.Fatalf("workflow %s failed to load: %s", path, outStr)
		}
		if strings.Contains(outStr, "Error finding sessions") {
			t.Fatalf("session logic failed for %s: %s", path, outStr)
		}
	}
}

package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/gkwa/easilydig/internal/logger"
)

func TestFetchCmd(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	customLogger := logger.NewConsoleLogger(0, false)
	cliLogger = customLogger

	cmd := fetchCmd
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	w.Close()

	os.Stdout = oldStdout

	var buf bytes.Buffer

	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to copy output: %v", err)
	}

	output := buf.String()
	t.Logf("Command output: %s", output)
}

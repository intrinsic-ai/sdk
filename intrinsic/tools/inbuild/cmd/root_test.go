// Copyright 2023 Intrinsic Innovation LLC

package root

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	RootCmd.SetArgs([]string{"--help"})

	// Capture stdout.
	oldStdout := RootCmd.OutOrStdout()
	defer RootCmd.SetOut(oldStdout)
	reader, writer, _ := os.Pipe()
	RootCmd.SetOut(writer)

	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("root.RootCmd.Execute() failed: %v", err)
	}

	writer.Close()
	RootCmd.SetOut(oldStdout)

	var buf bytes.Buffer
	io.Copy(&buf, reader)
	got := buf.String()

	if !strings.Contains(got, "Usage:") {
		t.Errorf("Expected output to contain 'Usage:', got: %q", got)
	}
}

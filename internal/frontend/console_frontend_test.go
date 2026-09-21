package frontend

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sarems/textscan/internal/core"
)

func TestConsoleOutputPreservesUTF8Text(t *testing.T) {
	document := core.NewDocument("notes.txt", "Café.")
	if err := document.Scoring.InitInRange(0, len(document.Text), 0.5); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	frontend := ConsoleOutputFrontend{Writer: &output}

	if err := frontend.DisplayOutput([]core.Document{*document}); err != nil {
		t.Fatalf("DisplayOutput() error = %v", err)
	}
	if !strings.Contains(output.String(), "é") {
		t.Fatalf("output did not preserve UTF-8 text: %q", output.String())
	}
}

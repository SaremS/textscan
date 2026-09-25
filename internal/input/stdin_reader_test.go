package codereader

import (
	"strings"
	"testing"

	"github.com/sarems/textscan/internal/core"
)

func TestNewStdinReaderRejectsInvalidFiletype(t *testing.T) {
	_, err := NewStdinReader(strings.NewReader(""), "rs")

	if err == nil {
		t.Fatal("NewStdinReader() error = nil, want an error")
	}
}

func TestStdinReaderParseTextUsesFiletypeInIdentifier(t *testing.T) {
	reader, err := NewStdinReader(strings.NewReader("package main"), "go")
	if err != nil {
		t.Fatal(err)
	}

	items, err := reader.ParseText(core.QuestionSet{"question"})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("ParseText() returned %d items, want 1", len(items))
	}
	if items[0].Identifier != "stdin.go" {
		t.Errorf("Identifier = %q, want %q", items[0].Identifier, "stdin.go")
	}
	if items[0].Text != "package main" {
		t.Errorf("Text = %q, want %q", items[0].Text, "package main")
	}
}

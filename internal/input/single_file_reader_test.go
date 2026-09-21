package input

import (
	"strings"
	"testing"
)

func TestNewSingleFileReaderRejectsNonTextFiles(t *testing.T) {
	_, err := NewSingleFileReader("notes.md")
	if err == nil || !strings.Contains(err.Error(), ".txt") {
		t.Fatalf("NewSingleFileReader() error = %v, want .txt validation error", err)
	}
}

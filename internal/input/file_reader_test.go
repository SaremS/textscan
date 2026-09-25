package codereader

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sarems/textscan/internal/core"
)

func TestNewFileReaderReturnsErrorForMissingTarget(t *testing.T) {
	missingTarget := filepath.Join(t.TempDir(), "missing")

	_, err := NewFileReader(missingTarget)

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("NewFileReader() error = %v, want an error matching os.ErrNotExist", err)
	}
}

func TestFileReaderParseTextReadsSingleFile(t *testing.T) {
	target := filepath.Join(t.TempDir(), "example.txt")
	if err := os.WriteFile(target, []byte("single file"), 0o600); err != nil {
		t.Fatal(err)
	}

	reader, err := NewFileReader(target)
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
	if items[0].Identifier != core.Identifier(target) {
		t.Errorf("Identifier = %q, want %q", items[0].Identifier, target)
	}
	if items[0].Text != "single file" {
		t.Errorf("Text = %q, want %q", items[0].Text, "single file")
	}
}

func TestFileReaderParseTextReadsDirectoryRecursively(t *testing.T) {
	target := filepath.Join(t.TempDir(), "source")
	nestedTarget := filepath.Join(target, "nested", "second.txt")
	firstTarget := filepath.Join(target, "first.txt")

	if err := os.MkdirAll(filepath.Dir(nestedTarget), 0o755); err != nil {
		t.Fatal(err)
	}
	for path, contents := range map[string]string{
		firstTarget:  "first file",
		nestedTarget: "second file",
	} {
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	reader, err := NewFileReader(target)
	if err != nil {
		t.Fatal(err)
	}

	items, err := reader.ParseText(core.QuestionSet{"question"})
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("ParseText() returned %d items, want 2", len(items))
	}

	got := make(map[core.Identifier]core.Text, len(items))
	for _, item := range items {
		got[item.Identifier] = item.Text
	}
	want := map[core.Identifier]core.Text{
		core.Identifier(firstTarget):  "first file",
		core.Identifier(nestedTarget): "second file",
	}
	for identifier, text := range want {
		if got[identifier] != text {
			t.Errorf("text for %q = %q, want %q", identifier, got[identifier], text)
		}
	}
}

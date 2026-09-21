package backend

import (
	"reflect"
	"testing"

	"github.com/sarems/textscan/internal/core"
)

func TestSplitParagraphs(t *testing.T) {
	text := core.Text(" First paragraph.\nContinues.\n\n\tSecond paragraph. \n\n\nThird.")

	got := splitParagraphs(text)
	want := []textSpan{
		{start: 1, end: 28},
		{start: 31, end: 48},
		{start: 52, end: 58},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitParagraphs() = %#v, want %#v", got, want)
	}
}

func TestSplitSentencesPreservesDocumentByteOffsets(t *testing.T) {
	text := core.Text("Prefix\n\nCafé is open. It closes late.")
	paragraph := textSpan{start: 8, end: len(text)}

	got, err := splitSentences(text, paragraph)
	if err != nil {
		t.Fatalf("splitSentences() error = %v", err)
	}

	want := []textSpan{
		{start: 8, end: 22},
		{start: 22, end: len(text)},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("splitSentences() = %#v, want %#v", got, want)
	}

	if first := string(text[got[0].start:got[0].end]); first != "Café is open." {
		t.Fatalf("first sentence = %q", first)
	}
}

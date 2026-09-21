package backend

import (
	"fmt"
	"strings"

	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/english"
	"github.com/sarems/textscan/internal/core"
)

type textSpan struct {
	start int
	end   int
}

func splitParagraphs(text core.Text) []textSpan {
	source := string(text)
	var spans []textSpan

	for start := 0; start < len(source); {
		for start < len(source) && isWhitespace(source[start]) {
			start++
		}

		if start == len(source) {
			break
		}

		end := start
		for end < len(source) {
			lineEnd := strings.IndexByte(source[end:], '\n')
			if lineEnd == -1 {
				end = len(source)
				break
			}

			lineEnd += end
			if strings.TrimSpace(source[end:lineEnd]) == "" {
				break
			}
			end = lineEnd + 1
		}

		trimmedEnd := end
		for trimmedEnd > start && isWhitespace(source[trimmedEnd-1]) {
			trimmedEnd--
		}
		if trimmedEnd > start {
			spans = append(spans, textSpan{start: start, end: trimmedEnd})
		}

		start = end + 1
	}

	return spans
}

func splitSentences(text core.Text, paragraph textSpan) ([]textSpan, error) {
	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		return nil, fmt.Errorf("create English sentence tokenizer: %w", err)
	}

	paragraphText := string(text[paragraph.start:paragraph.end])
	sentences := tokenizer.Tokenize(paragraphText)
	spans := make([]textSpan, 0, len(sentences))

	for _, sentence := range sentences {
		span, ok := sentenceSpan(paragraphText, paragraph.start, sentence)
		if ok {
			spans = append(spans, span)
		}
	}

	return spans, nil
}

func sentenceSpan(paragraph string, paragraphStart int, sentence *sentences.Sentence) (textSpan, bool) {
	if sentence == nil {
		return textSpan{}, false
	}

	start := paragraphStart + sentence.Start
	end := paragraphStart + sentence.End
	if start < paragraphStart || end > paragraphStart+len(paragraph) || start >= end {
		return textSpan{}, false
	}

	return textSpan{start: start, end: end}, true
}

func isWhitespace(value byte) bool {
	switch value {
	case ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}

package frontend

import (
	"fmt"
	"io"
	"os"

	"github.com/sarems/textscan/internal/core"
)

type AsyncConsoleOutputFrontend struct {
	Writer            io.Writer
	cutoffProbability float64
}

func NewAsyncConsoleOutputFrontend(cutoffProbability float64) *AsyncConsoleOutputFrontend {
	return &AsyncConsoleOutputFrontend{
		Writer:            os.Stdout,
		cutoffProbability: cutoffProbability,
	}
}

func (f *AsyncConsoleOutputFrontend) DisplayOutput(probsRangeChan <-chan core.ProbsRange) {
	printedRanges := make(map[string]struct{})

	for pr := range probsRangeChan {
		if pr.Probability < f.cutoffProbability {
			continue
		}

		hash := fmt.Sprintf("%q:%d:%d", pr.Question, pr.Start, pr.End)
		if _, alreadyPrinted := printedRanges[hash]; alreadyPrinted {
			continue
		}
		printedRanges[hash] = struct{}{}

		f.displayTextItemLinewise(pr)
	}
}

func (f *AsyncConsoleOutputFrontend) displayTextItemLinewise(item core.ProbsRange) {
	text := string(item.TextSnippet)

	if item.Probability >= f.cutoffProbability {
		fmt.Fprintf(f.Writer, "%s, L%d-L%d:\n\n", item.Identifier, item.StartLineOffset.LineNumber, item.EndLineOffset.LineNumber)
		fmt.Fprintln(f.Writer, text)
		fmt.Fprintf(f.Writer, "\np('%s')=%f\n\n\n", item.Question, item.Probability)
	}
}

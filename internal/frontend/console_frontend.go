package frontend

import (
	"fmt"
	"io"
	"os"

	"github.com/sarems/textscan/internal/core"
)

type ConsoleOutputFrontend struct {
	Writer io.Writer
}

func NewConsoleOutputFrontend() *ConsoleOutputFrontend {
	return &ConsoleOutputFrontend{Writer: os.Stdout}
}

func (f *ConsoleOutputFrontend) DisplayOutput(documents []core.Document) error {
	for _, document := range documents {
		if _, err := fmt.Fprintf(
			f.Writer,
			"----------\n\033[1m%s\033[0m\n----------\n\n",
			document.Identifier,
		); err != nil {
			return err
		}

		scores := document.Scoring.Items()
		text := string(document.Text)
		if len(scores) != len(text) {
			return fmt.Errorf("scoring/text length mismatch: %d scores for %d bytes", len(scores), len(text))
		}

		for index, character := range text {
			score := scores[index]
			if score < 0 {
				score = 0
			}
			if score > 1 {
				score = 1
			}

			if _, err := fmt.Fprintf(
				f.Writer,
				"\033[48;2;%d;0;0m%c\033[0m",
				int(score*255),
				character,
			); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(f.Writer); err != nil {
			return err
		}
	}

	return nil
}

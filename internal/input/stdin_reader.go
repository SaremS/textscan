package codereader

import (
	"fmt"
	"io"

	"github.com/sarems/textscan/internal/core"
)

type StdinReader struct {
	input      io.Reader
	identifier core.Identifier
}

func NewStdinReader(input io.Reader, filetype string) (*StdinReader, error) {
	switch filetype {
	case "go", "js", "py", "txt":
		return &StdinReader{
			input:      input,
			identifier: core.Identifier("stdin." + filetype),
		}, nil
	default:
		return nil, fmt.Errorf(
			"invalid stdin filetype %q: must be one of go, js, py, or txt",
			filetype,
		)
	}
}

func (r *StdinReader) ParseText(questions core.QuestionSet) ([]core.TextItem, error) {
	data, err := io.ReadAll(r.input)
	if err != nil {
		return nil, fmt.Errorf("read standard input: %w", err)
	}

	item, err := core.NewTextItem(r.identifier, core.Text(data), questions)
	if err != nil {
		return nil, fmt.Errorf(
			"create text item for %q: %w",
			r.identifier,
			err,
		)
	}

	return []core.TextItem{*item}, nil
}

package input

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sarems/textscan/internal/core"
)

type SingleFileReader struct {
	targetFile string
}

func NewSingleFileReader(targetFile string) (*SingleFileReader, error) {
	if filepath.Ext(targetFile) != ".txt" {
		return nil, fmt.Errorf("textscan accepts only .txt files: %s", targetFile)
	}

	return &SingleFileReader{targetFile: targetFile}, nil
}

func (r *SingleFileReader) ReadText() ([]core.Document, error) {
	data, err := os.ReadFile(r.targetFile)
	if err != nil {
		return nil, err
	}

	document := core.NewDocument(
		core.DocumentIdentifier(r.targetFile),
		core.Text(data),
	)

	return []core.Document{*document}, nil
}

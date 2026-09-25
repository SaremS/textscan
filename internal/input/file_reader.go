package codereader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sarems/textscan/internal/core"
)

type FileReader struct {
	target      string
	isDirectory bool
}

func NewFileReader(target string) (*FileReader, error) {
	info, err := os.Stat(target)
	if err != nil {
		return nil, fmt.Errorf("inspect input %q: %w", target, err)
	}

	if !info.IsDir() && !info.Mode().IsRegular() {
		return nil, fmt.Errorf("input %q is neither a regular file nor a directory", target)
	}

	return &FileReader{
		target:      target,
		isDirectory: info.IsDir(),
	}, nil
}

func (r *FileReader) ParseText(questions core.QuestionSet) ([]core.TextItem, error) {
	if !r.isDirectory {
		item, err := readTextItem(r.target, questions)
		if err != nil {
			return nil, err
		}

		return []core.TextItem{*item}, nil
	}

	items := make([]core.TextItem, 0)
	err := filepath.WalkDir(r.target, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}

		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("inspect file %q: %w", path, err)
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		item, err := readTextItem(path, questions)
		if err != nil {
			return err
		}

		items = append(items, *item)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read input directory %q: %w", r.target, err)
	}

	return items, nil
}

func readTextItem(path string, questions core.QuestionSet) (*core.TextItem, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	item, err := core.NewTextItem(
		core.Identifier(path),
		core.Text(data),
		questions,
	)
	if err != nil {
		return nil, fmt.Errorf("create text item for %q: %w", path, err)
	}

	return item, nil
}

package backend

import (
	"context"

	"github.com/sarems/textscan/internal/backend/plaintext"
	"github.com/sarems/textscan/internal/core"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
)

func parseCodeFromText(text core.Text, language core.ProgrammingLanguage) (*sitter.Tree, error) {
	parser := sitter.NewParser()
	defer parser.Close()

	switch language {
	case core.Go:
		parser.SetLanguage(golang.GetLanguage())
	case core.JavaScript:
		parser.SetLanguage(javascript.GetLanguage())
	case core.Python:
		parser.SetLanguage(python.GetLanguage())
	default:
		parser.SetLanguage(plaintext.GetLanguage())
	}

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(text))
	if err != nil {
		return nil, err
	}

	return tree, nil
}

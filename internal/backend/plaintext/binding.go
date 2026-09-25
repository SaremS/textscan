package plaintext

//#include "parser.h"
//TSLanguage *tree_sitter_plaintext();
import "C"

import (
	"unsafe"

	sitter "github.com/smacker/go-tree-sitter"
)

func GetLanguage() *sitter.Language {
	return sitter.NewLanguage(unsafe.Pointer(C.tree_sitter_plaintext()))
}

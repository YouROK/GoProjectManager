package autocomplete

import (
	"core/codeasist"
	"core/pfile"
	"unicode"

	"go/ast"
	"go/token"
)

type Cursor struct {
	LeftWords  []string
	RightWords []string
	LastLetter rune
	File       *pfile.GoFile
	Asist      *codeasist.Asist

	PackagePath string
	Source      string
	Offset      int
	FSet        *token.FileSet
	Node        ast.Node
}

func (c *Cursor) IsLastLetter() bool {
	return (unicode.IsLetter(c.LastLetter) || unicode.IsDigit(c.LastLetter))
}

func (c *Cursor) GetLastWord() string {
	return c.LeftWords[len(c.LeftWords)-1]
}

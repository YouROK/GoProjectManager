package pfile

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	//	"log"
)

func Format(src string) (string, error) {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func Refactor(src string) {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "", src, 0)
	//	fset.
	ast.Print(fset, f)
}

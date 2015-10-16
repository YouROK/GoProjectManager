package codecomplete

import (
	"go/ast"
	"go/parser"
	"go/token"

	"fmt"
)

func Complete(src string) {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "", src, 0)
	ast.Print(fset, node.Name)
	for _, u := range node.Scope.Objects {
		fmt.Println(u)
		if f, ok := u.Decl.(*ast.FuncDecl); ok {
			ast.Print(fset, f)
		}
	}
}

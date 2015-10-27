package codecomplete

import (
	//	"go/ast"
	"go/parser"
	"go/token"
)

func getAllTypes(cur cursor) []string {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "", cur.src, 0)
	list := []string{}
	for _, k := range node.Scope.Objects {
		if k.Kind.String() == "type" {
			list = append(list, k.Name)
		}
	}
	return list
}

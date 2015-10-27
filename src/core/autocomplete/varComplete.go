package autocomplete

import (
	"core/codeasist"

	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"fmt"
)

func getVarCandidates(src string, offset int) {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "", src, 0)
	ptr := GetPrevNode(node, offset)
	if ptr != nil {
		if node, ok := ptr.(*ast.Ident); ok {
			node := getTypeNode(node)
			nameNode := codeasist.GetExprName(node)
			fmt.Println(nameNode)
			if strings.Index(nameNode, ".") > 0 {
				name := strings.Split(nameNode, ".")
				fmt.Println(name)
			}
			//ast.Print(fset, typen)
		}
	}
}

func getTypeNode(root *ast.Ident) ast.Expr {
	node := root
	if node, ok := node.Obj.Decl.(*ast.AssignStmt); ok { //if assign
		for i := 0; i < len(node.Lhs); i++ {
			name := codeasist.GetExprName(node.Lhs[i])
			if name == root.Name && node.Rhs != nil && len(node.Rhs) > 0 {
				if i < len(node.Rhs) {
					return node.Rhs[i]
				}
				return node.Rhs[0]
			}
		}
	}
	return nil
}

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

func getOffset(src string, line, col int) int {
	l, c := 1, 1
	for i, ch := range src {
		//		fmt.Println(i, l, c, line, col)
		if l >= line && c >= col {
			return i
		}
		if ch == '\n' {
			l++
			c = 0
		}
		c++
	}
	return -1
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

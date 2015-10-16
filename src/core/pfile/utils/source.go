package utils

import (
	"go/ast"
	"go/parser"
	"go/token"

	"fmt"
	//	"strconv"
	//	"reflect"
)

var (
	types map[string]bool
)

func Get(src string) {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "", src, 0)
	/*ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.FuncType:
			s = src[x.Pos()-1 : x.End()]
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}

		//		if n != nil {
		//			s = reflect.TypeOf(n).String()
		//			types[s] = true
		//			//			fmt.Println(s)
		//		}
		return true
	})*/

	v := MyVisitor{}

	ast.Walk(v, node)
}

type MyVisitor struct {
}

func (v MyVisitor) Visit(node ast.Node) ast.Visitor {
	if c, ok := node.(*ast.FuncDecl); ok {
		fmt.Println(c.Recv)
		fmt.Println(c, "\n")
	}
	return ast.Visitor(v)
}

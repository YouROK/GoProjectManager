package autocomplete

import (
	"fmt"
	"go/ast"
	"go/token"
)

type prevFinder struct {
	Off int
	Ptr ast.Node
}

func (v *prevFinder) Visit(node ast.Node) ast.Visitor {
	if _, ok := node.(*ast.File); ok {
		return ast.Visitor(v)
	}
	if node != nil && node.Pos() != token.NoPos {
		fmt.Println(node.Pos(), v.Off)
		fmt.Println(node)
		if int(node.Pos()) < v.Off && int(node.End()) < v.Off {
			ok := false
			if v.Ptr == nil {
				ok = true
			} else {
				if v.Off-int(node.End()) < v.Off-int(v.Ptr.End()) {
					ok = true
				}
			}
			if ok {
				v.Ptr = node
				return nil
			}
		}
	}
	return ast.Visitor(v)
}

func GetPrevNode(root ast.Node, off int) ast.Node {
	mv := prevFinder{}
	mv.Off = off
	ast.Walk(&mv, root)
	return mv.Ptr
}

func GetOutNode(root ast.Node, off int) ast.Node {
	var ptr ast.Node
	ast.Inspect(root, func(node ast.Node) bool {
		if _, ok := node.(*ast.File); ok {
			return true
		}
		if node != nil && node.Pos() != token.NoPos {
			if ptr == nil && off > int(node.Pos()) && off < int(node.End()) {
				ptr = node
			}
		}
		return true
	})
	return ptr
}

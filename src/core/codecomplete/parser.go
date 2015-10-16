package codecomplete

import (
	"go/ast"
)

type Node struct {
	parent  *Node
	childs  []*Node
	astNode ast.Node
}

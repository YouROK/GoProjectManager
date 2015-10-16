package codeasist

import (
	"fmt"
	"go/ast"
	"go/token"
)

func GetType(node ast.Node) string {
	switch n := node.(type) {
	case *ast.ArrayType:
		return "array type"
	case *ast.AssignStmt:
		return "assignment"
	case *ast.BadDecl:
		return "bad declaration"
	case *ast.BadExpr:
		return "bad expression"
	case *ast.BadStmt:
		return "bad statement"
	case *ast.BasicLit:
		return "basic literal"
	case *ast.BinaryExpr:
		return fmt.Sprintf("binary %s operation", n.Op)
	case *ast.BlockStmt:
		return "block"
	case *ast.BranchStmt:
		switch n.Tok {
		case token.BREAK:
			return "break statement"
		case token.CONTINUE:
			return "continue statement"
		case token.GOTO:
			return "goto statement"
		case token.FALLTHROUGH:
			return "fall-through statement"
		}
	case *ast.CallExpr:
		return "function call (or conversion)"
	case *ast.CaseClause:
		return "case clause"
	case *ast.ChanType:
		return "channel type"
	case *ast.CommClause:
		return "communication clause"
	case *ast.Comment:
		return "comment"
	case *ast.CommentGroup:
		return "comment group"
	case *ast.CompositeLit:
		return "composite literal"
	case *ast.DeclStmt:
		return GetType(n.Decl) + " statement"
	case *ast.DeferStmt:
		return "defer statement"
	case *ast.Ellipsis:
		return "ellipsis"
	case *ast.EmptyStmt:
		return "empty statement"
	case *ast.ExprStmt:
		return "expression statement"
	case *ast.Field:
		// Can be any of these:
		// struct {x, y int}  -- struct field(s)
		// struct {T}         -- anon struct field
		// interface {I}      -- interface embedding
		// interface {f()}    -- interface method
		// func (A) func(B) C -- receiver, param(s), result(s)
		return "field/method/parameter"
	case *ast.FieldList:
		return "field/method/parameter list"
	case *ast.File:
		return "source file"
	case *ast.ForStmt:
		return "for loop"
	case *ast.FuncDecl:
		return "function declaration"
	case *ast.FuncLit:
		return "function literal"
	case *ast.FuncType:
		return "function type"
	case *ast.GenDecl:
		switch n.Tok {
		case token.IMPORT:
			return "import declaration"
		case token.CONST:
			return "constant declaration"
		case token.TYPE:
			return "type declaration"
		case token.VAR:
			return "variable declaration"
		}
	case *ast.GoStmt:
		return "go statement"
	case *ast.Ident:
		return "identifier"
	case *ast.IfStmt:
		return "if statement"
	case *ast.ImportSpec:
		return "import specification"
	case *ast.IncDecStmt:
		if n.Tok == token.INC {
			return "increment statement"
		}
		return "decrement statement"
	case *ast.IndexExpr:
		return "index expression"
	case *ast.InterfaceType:
		return "interface type"
	case *ast.KeyValueExpr:
		return "key/value association"
	case *ast.LabeledStmt:
		return "statement label"
	case *ast.MapType:
		return "map type"
	case *ast.Package:
		return "package"
	case *ast.ParenExpr:
		return "parenthesized " + GetType(n.X)
	case *ast.RangeStmt:
		return "range loop"
	case *ast.ReturnStmt:
		return "return statement"
	case *ast.SelectStmt:
		return "select statement"
	case *ast.SelectorExpr:
		return "selector"
	case *ast.SendStmt:
		return "channel send"
	case *ast.SliceExpr:
		return "slice expression"
	case *ast.StarExpr:
		return "*-operation" // load/store expr or pointer type
	case *ast.StructType:
		return "struct type"
	case *ast.SwitchStmt:
		return "switch statement"
	case *ast.TypeAssertExpr:
		return "type assertion"
	case *ast.TypeSpec:
		return "type specification"
	case *ast.TypeSwitchStmt:
		return "type switch"
	case *ast.UnaryExpr:
		return fmt.Sprintf("unary %s operation", n.Op)
	case *ast.ValueSpec:
		return "value specification"
	}
	return ""
}

func GetExprName(node ast.Expr) string {
	tmp := ""
	switch x := node.(type) {
	case *ast.Ident:
		tmp = x.String()
	case *ast.ParenExpr:
		tmp = GetExprName(x.X)
	case *ast.InterfaceType:
		tmp = "interface"
	case *ast.StructType:
		tmp = "struct"
	case *ast.UnaryExpr:
		tmp = "&" + GetExprName(x.X)
	case *ast.StarExpr:
		tmp = "*" + GetExprName(x.X)
	case *ast.ArrayType:
		tmp = "[]" + GetExprName(x.Elt)
	case *ast.MapType:
		tmp = "map[" + GetExprName(x.Key) + "]" + GetExprName(x.Value)
	case *ast.ChanType:
		tmp = "chan " + GetExprName(x.Value)
	case *ast.CompositeLit:
		tmp = GetExprName(x.Type)
	case *ast.SelectorExpr:
		tmp = GetExprName(x.X) + "." + GetExprName(x.Sel)
	}
	return tmp
}

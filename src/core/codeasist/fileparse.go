package codeasist

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func fillFile(f *File, src string) {
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "", src, 0)
	for _, i := range node.Imports {
		imp := Unit{}
		if i.Path != nil {
			imp.Name = i.Path.Value
		}
		if i.Doc != nil {
			imp.Doc = i.Doc.Text()
		}

		imp.Line, imp.Column = getPosition(fset, i)
		imp.FileName = f.Path
		f.Imports = append(f.Imports, imp)
	}

	//add types,vars,consts
	for _, decl := range node.Decls {
		if node, ok := decl.(*ast.GenDecl); ok {
			if len(node.Specs) > 0 {
				if node.Tok.String() == "type" {
					if node, ok := node.Specs[0].(*ast.TypeSpec); ok {
						elm := Type{}
						if fillType(node, src, &elm) {
							elm.Line, elm.Column = getPosition(fset, node)
							elm.FileName = f.Path
							f.Types = append(f.Types, elm)
						}
					}
				} else if node.Tok.String() == "const" || node.Tok.String() == "var" {
					for _, n := range node.Specs {
						if node, ok := n.(*ast.ValueSpec); ok {
							vars := fillVarsConst(node, fset, src, f.Path)
							f.Vars = append(f.Vars, vars...)
						}
					}
				}
			}
		}
	}

	//add functions
	for _, decl := range node.Decls {
		if node, ok := decl.(*ast.FuncDecl); ok {
			elm := Function{}
			if fillFunc(node, src, &elm) {
				elm.Line, elm.Column = getPosition(fset, node)
				elm.FileName = f.Path
				f.Functions = append(f.Functions, elm)
			}
		}
	}
}

func fillVarsConst(node *ast.ValueSpec, fset *token.FileSet, src string, path string) []Var {
	vars := make([]Var, 0, 1)
	for i, n := range node.Names {
		if n.Name == "_" { // exmpl: var current, _ = os.Getwd()
			continue
		}
		u := Var{}
		u.Name = n.Name
		u.Line, u.Column = getPosition(fset, node)
		u.FileName = path
		if node.Doc != nil {
			u.Doc = node.Doc.Text()
		}
		if node.Values != nil {
			if _, ok := node.Values[i].(*ast.BasicLit); ok {
				u.Type = node.Values[i].(*ast.BasicLit).Value
			}
		}
		if node.Type != nil {
			u.Type = src[node.Type.Pos()-1 : node.Type.End()]
		}
		vars = append(vars, u)
	}
	return vars
}

func fillType(node *ast.TypeSpec, src string, elm *Type) bool {
	if node.Name != nil {
		if !ast.IsExported(node.Name.Name) {
			return false
		}
		elm.Name = node.Name.Name
	}
	if node.Doc != nil {
		elm.Doc = node.Doc.Text()
	}

	elm.Type = GetExprName(node.Type)
	return true
}

func fillFunc(node *ast.FuncDecl, src string, fnct *Function) bool {
	if node.Name != nil {
		if !ast.IsExported(node.Name.Name) {
			return false
		}
		fnct.Name = node.Name.Name
	} else {
		return false
	}

	if node.Doc != nil {
		fnct.Doc = node.Doc.Text()
	}
	if node.Type != nil {
		if node.Type.Params != nil {
			fnct.Args = src[node.Type.Params.Pos() : node.Type.Params.End()-2]
		}
		if node.Type.Results != nil {
			fnct.Type = src[node.Type.Results.Pos()-1 : node.Type.Results.End()]
		}
	}

	if node.Recv != nil && len(node.Recv.List) > 0 {
		expr := node.Recv.List[0].Type
		recv := GetExprName(expr)
		if recv != "" {
			fnct.Recv = recv
		}
	}
	return true
}

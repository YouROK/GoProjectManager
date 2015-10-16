package autocomplete

import (
	"go/ast"
	"path/filepath"
	"strings"

	"fmt"
)

func ProcessKeywords(cur *Cursor) Candidates {
	keywords := Candidates{}

	keywords = keywords.AddStrs(go_keywords...)
	pkgs := cur.Asist.GetPackages()
	for _, p := range pkgs {
		keywords = keywords.AddStr(p.Name)
	}

	if cur.IsLastLetter() {
		lastword := cur.GetLastWord()
		keywords.SetPercent(lastword)
		keywords.Sort()
	}
	return keywords
}

func ProcessPackage(cur *Cursor) Candidates {
	if cur.LeftWords != nil && len(cur.LeftWords) > 0 {
		last := strings.ToLower(cur.GetLastWord())
		if cur.LastLetter == ' ' && last == "package" {
			return Candidates{Candidate{Name: filepath.Base(filepath.Dir(cur.File.GetPath()))}}
		}
	}
	return nil
}

func ProcessFunc(cur *Cursor, node *ast.FuncDecl) Candidates {

	ast.Print(cur.FSet, GetOutNode(cur.Node, cur.Offset))

	if cur.LeftWords != nil && len(cur.LeftWords) > 0 {
		cands := Candidates{}
		lastw := strings.ToLower(cur.GetLastWord())
		pkg := cur.Asist.GetPackage(cur.PackagePath)
		if lastw == "func" && cur.LastLetter == ' ' { //begin of func
			file := pkg.GetFile(cur.File.GetPath())
			if file == nil {
				return nil
			}
			for _, t := range file.Types {
				fmt.Println(t.Name, t.Type)
				tmp := "(" + strings.ToLower(string(t.Name[0])) + " " + t.Name + ")"
				cand := Candidate{Name: tmp}
				if t.Type == "struct" {
					cand.Perc = 50
				}
				cands = append(cands, cand)
			}
			return cands
		}
	}
	return nil
}

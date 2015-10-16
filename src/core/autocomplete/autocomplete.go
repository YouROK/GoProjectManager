package autocomplete

import (
	"core/codeasist"
	"core/pfile"

	"go/ast"
	"go/parser"
	"go/token"

	"path/filepath"
	"strings"
	"unicode"

	"fmt"
)

func GetCandidates(file *pfile.GoFile, linenum, columnnum int, asist *codeasist.Asist) []string {
	srcline, err := file.GetLine(linenum)
	if err != nil {
		return nil
	}

	src := file.GetBody()
	fset := token.NewFileSet()
	node, _ := parser.ParseFile(fset, "", src, 0)
	off := getOffset(src, linenum, columnnum)

	cur := &Cursor{}
	cur.PackagePath = filepath.Dir(file.GetPath())
	cur.File = file
	cur.Asist = asist
	cur.LastLetter = ' '
	cur.Node = node
	cur.Source = src
	cur.Offset = off
	cur.FSet = fset

	if len(srcline) < columnnum {
		columnnum = len(srcline)
	}
	if columnnum > 0 {
		cur.LastLetter = []rune(srcline)[columnnum-1]
	} else if columnnum == 0 {
		cur.LastLetter = ' '
	}
	cur.LeftWords, cur.RightWords = GetLastWord(srcline, columnnum)
	fmt.Println()

	fmt.Println(cur.LeftWords, cur.RightWords)
	fmt.Println(string(cur.LastLetter))
	fmt.Println(process(cur))

	return nil
}

func process(cur *Cursor) Candidates {
	ptr := GetOutNode(cur.Node, cur.Offset)
	if ptr != nil {
		if node, ok := ptr.(*ast.FuncDecl); ok {
			return ProcessFunc(cur, node)
		}
	}
	return nil
}

func GetLastWord(linesrc string, off int) ([]string, []string) {

	wordsLeft := strings.FieldsFunc(linesrc[:off], func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r))
	})

	wordsRight := strings.FieldsFunc(linesrc[off:], func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r))
	})

	return wordsLeft, wordsRight
}

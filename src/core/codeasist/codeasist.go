package codeasist

import (
	"core/project"

	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

type Asist struct {
	pkgs []Package
	proj *project.Project
}

func NewAsist(proj *project.Project) *Asist {
	a := &Asist{}
	a.proj = proj
	return a
}

func (a *Asist) GetPackages() []Package {
	return a.pkgs
}

func (a *Asist) GetPackage(pkgName string) *Package {
	for _, p := range a.pkgs {
		lpdir := len(p.Dir)
		lpn := len(pkgName)
		if lpn > lpdir {
			continue
		}
		if p.Dir[lpdir-lpn:] == pkgName {
			return &p
		}
	}
	return nil
}

func (a *Asist) Scan() {
	pkgs := a.proj.GetPackages()
	a.pkgs = []Package{}
	for _, p := range pkgs {
		pkg := Package{}
		pkg.Dir = p.GetPath()
		pkg.Name = filepath.Base(pkg.Dir)
		pkgFileList := a.proj.GetPackageFiles(p)
		for _, f := range pkgFileList {
			if !f.IsOpen() {
				f.Open()
			}
			file := File{}
			file.Path = f.GetPath()
			fillFile(&file, f.GetBody())
			pkg.Files = append(pkg.Files, file)
		}
		a.pkgs = append(a.pkgs, pkg)
	}
}

func (p Package) GetStructFunctions(structName string) []Function {
	functions := make([]Function, 0)
	for _, file := range p.Files {
		for _, fun := range file.Functions {
			if fun.Recv == structName {
				functions = append(functions, fun)
			}
		}
	}
	return functions
}

func (p Package) GetFile(name string) *File {
	for _, f := range p.Files {
		if strings.Contains(f.Path, name) {
			return &f
		}
	}
	return nil
}

func (p Package) GetTypes() []Type {
	typeList := []Type{}
	for _, f := range p.Files {
		typeList = append(typeList, f.Types...)
	}
	return typeList
}

func getPosition(fset *token.FileSet, node ast.Node) (line int, col int) {
	position := fset.Position(node.Pos())
	line = position.Line
	col = position.Column
	return
}

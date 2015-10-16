package project

import (
	"core/pfile"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func (p *Project) scanProject() error {
	src := filepath.Join(p.projectDir, "src")
	p.files = make([]pfile.IFile, 0)
	walk := func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		fn, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		var file pfile.IFile
		if info.IsDir() {
			file = pfile.NewDir(fn)
		} else if filepath.Ext(info.Name()) == ".go" {
			file = pfile.NewGoFile(fn)
		} else {
			file = pfile.NewAnyFile(fn)
		}
		p.files = append(p.files, file)
		return nil
	}
	return filepath.Walk(src, walk)
}

func (p *Project) GetFileList() []pfile.IFile {
	if len(p.files) == 0 {
		p.scanProject()
	}
	return p.files
}

func (p *Project) GetFile(name string) pfile.IFile {
	if len(p.files) == 0 {
		p.scanProject()
	}
	for _, f := range p.files {
		if strings.HasSuffix(f.GetPath(), name) {
			return f
		}
	}
	return nil
}

func (p *Project) GetPackages() []pfile.Dir {
	if len(p.files) == 0 {
		p.scanProject()
	}
	pkgs := make([]pfile.Dir, 0)
	for _, f := range p.files {
		if dir, ok := f.(*pfile.Dir); ok {
			if len(p.GetPackageFiles(*dir)) > 0 {
				pkgs = append(pkgs, *dir)
			}
		}
	}
	return pkgs
}

func (p *Project) GetPackageFiles(dir pfile.Dir) []*pfile.GoFile {
	if len(p.files) == 0 {
		p.scanProject()
	}
	gfiles := make([]*pfile.GoFile, 0)
	pkgFindName := filepath.Base(dir.GetPath())
	for _, f := range p.files {
		if gf, ok := f.(*pfile.GoFile); ok {
			filePkgName := filepath.Base(filepath.Dir(gf.GetPath()))
			if filePkgName == pkgFindName {
				gfiles = append(gfiles, gf)
			}
		}
	}
	return gfiles
}

func (p *Project) AddNewFile(path string) error {
	if filepath.Ext(path) != ".go" {
		return errors.New("File ext not .go")
	}

	fn, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(p.projectDir, filepath.Base(fn)) {
		return errors.New("File not in project dir")
	}

	file := pfile.NewGoFile(fn)
	p.files = append(p.files, file)
	return nil
}

func (p *Project) NewGoFile(path string) error {
	if filepath.Ext(path) != ".go" {
		return errors.New("File ext not .go")
	}

	fn, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(p.projectDir, filepath.Base(fn)) {
		return errors.New("File not in project dir")
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	f.Close()

	file := pfile.NewGoFile(fn)
	p.files = append(p.files, file)
	return nil
}

//func (p *Project)

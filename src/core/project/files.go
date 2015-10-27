package project

import (
	"core/project/file"
	"os"
	"path/filepath"
	"strings"
)

func (p *Project) GetFileList() []*file.File {
	return p.files
}

func (p *Project) GetFile(name string) *file.File {
	for _, f := range p.files {
		if strings.HasSuffix(f.GetPath(), name) {
			return f
		}
	}
	return nil
}

func (p *Project) AddFile(relpath string) error {
	ppath := filepath.Join(p.projectDir, relpath)
	if _, err := os.Lstat(ppath); err != nil {
		return err
	}

	if p.GetFile(relpath) != nil {
		return nil
	}

	file := file.NewFile(relpath)
	p.files = append(p.files, file)
	return nil
}

func (p *Project) NewFile(relpath string) error {
	ppath := filepath.Join(p.projectDir, relpath)
	f, err := os.Create(ppath)
	if err != nil {
		return err
	}

	f.Close()

	if p.GetFile(relpath) != nil {
		return nil
	}
	file := file.NewFile(relpath)
	p.files = append(p.files, file)
	return nil
}

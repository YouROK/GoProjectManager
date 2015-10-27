package project

import (
	"os"
	"path/filepath"
)

func (p *Project) SaveTempProject(path string) error {
	_, err := os.Lstat(path)
	var isExist = false
	if err == nil {
		isExist = true
	}
	for _, f := range p.GetFileList() {
		fn := filepath.Join(path, f.GetPath())
		err = os.MkdirAll(filepath.Dir(fn), 0755)
		if err != nil {
			break
		}
		_, err = f.SaveAs(fn)
		if err != nil {
			break
		}
	}
	if err != nil && !isExist {
		os.RemoveAll(path)
	}
	return err
}

package project

import (
	"core/pfile"
	"path/filepath"
)

type Project struct {
	projectFile string
	projectDir  string
	cfg         *Config
	files       []pfile.IFile
}

func OpenProject(fileName string) (*Project, error) {
	var err error
	p := &Project{}
	p.projectFile = fileName
	p.projectDir = filepath.Dir(fileName)
	p.cfg, err = loadConfig(fileName)
	p.scanProject()
	return p, err
}

func (p *Project) SaveProject() error {
	return saveConfig(p.projectFile, p.cfg)
}

func (p *Project) GetProjectPath() string {
	return p.projectDir
}

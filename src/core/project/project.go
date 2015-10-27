package project

import (
	"core/project/file"
	"path/filepath"
)

type Project struct {
	projectFile string
	projectDir  string
	cfg         *Config
	files       []*file.File
}

func OpenProject(fileName string) (*Project, error) {
	var err error
	p := &Project{}
	p.projectFile, _ = filepath.Abs(fileName)
	p.projectDir = filepath.Dir(p.projectFile)
	p.cfg, err = loadConfig(fileName)
	for _, fn := range p.cfg.Files {
		fp := file.NewFile(fn)
		p.files = append(p.files, fp)
		fp.Open(p.projectDir)
	}

	return p, err
}

func (p *Project) SaveProject() error {
	p.cfg.Files = make([]string, 0)
	for _, f := range p.files {
		p.cfg.Files = append(p.cfg.Files, f.GetPath())
	}
	return saveConfig(p.projectFile, p.cfg)
}

func (p *Project) GetProjectPath() string {
	return p.projectDir
}

func (p *Project) GetEnvs() []string {

	envs := make([]string, 0)

	if p.cfg.GoArch != "" {
		envs = append(envs, "GOARCH="+p.cfg.GoArch)
	}
	if p.cfg.GoOS != "" {
		envs = append(envs, "GOOS="+p.cfg.GoOS)
	}
	if p.cfg.GoRoot != "" {
		envs = append(envs, "GOROOT="+p.cfg.GoRoot)
	}
	gopath, _ := filepath.Abs(p.GetProjectPath())
	envs = append(envs, "GOPATH="+gopath)
	return envs
}

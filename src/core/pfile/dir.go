package pfile

type Dir struct {
	IFile
	path string
}

func NewDir(path string) *Dir {
	f := &Dir{}
	f.path = path
	return f
}

func (d *Dir) GetPath() string {
	return d.path
}

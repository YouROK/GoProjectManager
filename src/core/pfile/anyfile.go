package pfile

import (
	"io"
	"os"
)

type AnyFile struct {
	io.ReadWriter
	IFile
	path string
}

func NewAnyFile(path string) *AnyFile {
	f := &AnyFile{}
	f.path = path
	return f
}

func (f *AnyFile) Read(p []byte) (n int, err error) {
	ff, err := os.Open(f.path)
	if err != nil {
		return 0, err
	}
	defer ff.Close()
	return ff.Read(p)
}

func (f *AnyFile) Write(p []byte) (n int, err error) {
	ff, err := os.Create(f.path)
	if err != nil {
		return 0, err
	}
	defer ff.Close()
	return ff.Write(p)
}

func (f *AnyFile) GetPath() string {
	return f.path
}

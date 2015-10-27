package file

import (
	"bufio"
	"bytes"
	"errors"
	"go/format"
	"os"
	"path/filepath"
)

type File struct {
	path     string
	isOpen   bool
	isModify bool
	isSource bool

	srcLines []string
}

func NewFile(path string) *File {
	f := &File{}
	f.path = path
	return f
}

func (f *File) GetPath() string {
	return f.path
}

func (f *File) GetOffset(x, y int) int {
	if y == 0 {
		if x >= 0 {
			return x
		} else {
			return 0
		}
	}

	off := 0
	for cl, s := range f.srcLines {
		if cl >= y {
			break
		}
		off += len(s) + 1
	}
	return off + x
}

func (f *File) GetLine(n int) (string, error) {
	if n > 0 && n <= len(f.srcLines) {
		return f.srcLines[n-1], nil
	}
	return "", errors.New("Out of range")
}

func (f *File) SetLine(n int, line string) error {
	if n > 0 && n <= len(f.srcLines) {
		f.srcLines[n-1] = line
		f.isModify = true
		return nil
	}
	return errors.New("Out of range")
}

func (f *File) GetLines() []string {
	return f.srcLines
}

func (f *File) SetLines(lines []string) {
	f.srcLines = lines
	f.isModify = true
}

func (f *File) GetBody() string {
	ret := ""
	for _, l := range f.srcLines {
		ret += l + "\n"
	}
	return ret
}

func (f *File) SetBody(body string) {
	buf := bytes.NewBufferString(body)
	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)

	f.isModify = true
	f.srcLines = make([]string, 0, 10)
	for scanner.Scan() {
		f.srcLines = append(f.srcLines, scanner.Text())
	}
}

func (f *File) Open(projectPath string) error {
	apath := filepath.Join(projectPath, f.path)
	ff, err := os.Open(apath)
	if err != nil {
		return err
	}
	defer ff.Close()
	scanner := bufio.NewScanner(ff)
	scanner.Split(bufio.ScanLines)

	f.srcLines = make([]string, 0, 10)
	for scanner.Scan() {
		f.srcLines = append(f.srcLines, scanner.Text())
	}
	f.isOpen = true
	f.isModify = false
	f.isSource = false
	if filepath.Ext(f.path) == ".go" {
		f.isSource = true
	}
	return nil
}

func (f *File) Save() error {
	ff, err := os.Create(f.path)
	if err != nil {
		return err
	}
	defer ff.Close()

	_, err = ff.WriteString(f.GetBody())
	if err != nil {
		return err
	}
	f.isModify = false
	return nil
}

func (f *File) SaveAs(filename string) (*File, error) {
	fn := NewFile(filename)
	fn.srcLines = make([]string, len(f.srcLines))
	copy(fn.srcLines, f.srcLines)
	err := fn.Save()
	if err != nil {
		os.Remove(filename)
		return nil, err
	}
	return fn, nil
}

func (f *File) IsOpen() bool {
	return f.isOpen
}

func (f *File) IsModify() bool {
	return f.isModify
}

func (f *File) Format() {
	if len(f.srcLines) > 0 {
		formatted, err := format.Source([]byte(f.GetBody()))
		if err == nil {
			f.SetBody(string(formatted))
		}
	}
}

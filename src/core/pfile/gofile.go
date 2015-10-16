package pfile

import (
	"bufio"
	"bytes"
	"errors"
	"os"
)

type GoFile struct {
	IFile
	path     string
	isOpen   bool
	isModify bool

	currentBody []string
}

func NewGoFile(path string) *GoFile {
	f := &GoFile{}
	f.path = path
	return f
}

func (f *GoFile) GetPath() string {
	return f.path
}

func (f *GoFile) GetLine(n int) (string, error) {
	if n > 0 && n <= len(f.currentBody) {
		return f.currentBody[n-1], nil
	}
	return "", errors.New("Out of range")
}

func (f *GoFile) SetLine(n int, line string) error {
	if n > 0 && n <= len(f.currentBody) {
		f.currentBody[n-1] = line
		f.isModify = true
		return nil
	}
	return errors.New("Out of range")
}

func (f *GoFile) GetLines() []string {
	return f.currentBody
}

func (f *GoFile) SetLines(lines []string) {
	f.currentBody = lines
	f.isModify = true
}

func (f *GoFile) GetBody() string {
	ret := ""
	for _, l := range f.currentBody {
		ret += l + "\n"
	}
	return ret
}

func (f *GoFile) SetBody(body string) {
	buf := bytes.NewBufferString(body)
	scanner := bufio.NewScanner(buf)
	scanner.Split(bufio.ScanLines)

	f.isModify = true
	f.currentBody = make([]string, 0, 10)
	for scanner.Scan() {
		f.currentBody = append(f.currentBody, scanner.Text())
	}
}

func (f *GoFile) Open() error {
	ff, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer ff.Close()
	scanner := bufio.NewScanner(ff)
	scanner.Split(bufio.ScanLines)

	f.currentBody = make([]string, 0, 10)
	for scanner.Scan() {
		f.currentBody = append(f.currentBody, scanner.Text())
	}
	f.isOpen = true
	f.isModify = false
	return nil
}

func (f *GoFile) Save() error {
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

func (f *GoFile) IsOpen() bool {
	return f.isOpen
}

func (f *GoFile) IsModify() bool {
	return f.isModify
}

func (f *GoFile) Format() {
	if len(f.currentBody) > 0 {
		formated, err := Format(f.GetBody())
		if err == nil {
			f.SetBody(formated)
		}
	}
}

// func (f*GoFile)

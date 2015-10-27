package codecomplete

import (
	"bytes"
	"core/project"
	"core/project/file"
	"core/run"
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func Init(p *project.Project) error {
	tmpPath := p.GetProjectPath() + "/.tmp"
	os.MkdirAll(tmpPath, 0755)
	return p.SaveTempProject(tmpPath)
}

func Complete(p *project.Project, pfile *file.File, off int) Candidates {
	tmpPath := p.GetProjectPath() + "/.tmp/" + pfile.GetPath()
	pfile.SaveAs(tmpPath)

	src := pfile.GetBody()
	cur := getCursor(src, off)
	//	fmt.Println(src[off-3:off] + "|" + src[off:off+3])

	//TODO
	cands := gocode(p, tmpPath, cur)
	cands = append(cands, keywords(cur)...)
	cands = append(cands, packages(cur, pfile)...)
	cands = append(cands, funcsRecv(cur)...)

	cands.Sort()

	ret := Candidates{}

	ret = append(ret, Candidate{Name: cur.partial, Perc: 200})
	i := 0

	for _, c := range cands {
		if c.Perc < 0 && i > 10 {
			break
		}
		ret = append(ret, c)
		i++
	}

	return ret
}

func gocode(p *project.Project, filename string, cur cursor) Candidates {
	out := bytes.NewBuffer([]byte{})
	cmd := run.NewExec("gocode", []string{"-f=json", "--in=" + filename, "autocomplete", strconv.Itoa(cur.offset)}, p.GetEnvs())
	cmd.StdOut = out
	cmd.Run()

	var cands Candidates
	pos := strings.Index(out.String()[1:], "[")
	buf := out.Bytes()
	if json.Unmarshal(buf[pos+1:len(buf)-1], &cands) == nil {
		if len(cur.partial) > 1 {
			cands.SetPercent(cur.partial)
		}
	}
	return cands
}

func funcsRecv(cur cursor) Candidates {
	if cur.token.token().tok == token.FUNC {
		fmt.Println(cur.token.token().Literal())
		cands := Candidates{}
		list := getAllTypes(cur)
		for _, t := range list {
			cands = append(cands, NewCandidate("", "("+strings.ToLower(string(t[0]))+" *"+t+")", ""))
		}
		return cands
	}
	return nil
}

func packages(cur cursor, f *file.File) Candidates {
	if cur.token.previous_token() {
		if cur.token.token().tok == token.PACKAGE {
			pkg := filepath.Base(filepath.Dir(f.GetPath()))
			cands := Candidates{NewCandidate("package", pkg, ""), NewCandidate("package", "main", "")}
			if len(cur.partial) > 1 {
				cands.SetPercent(cur.partial)
			}
			return cands
		}
	}
	return nil
}

func keywords(cur cursor) Candidates {
	if cur.token.token().tok != token.PERIOD {
		if cur.token.previous_token() && cur.token.token().tok == token.PERIOD {
			return nil
		}
		kc := getKeywordsCandidates()
		if cur.partial != "" {
			kc.SetPercent(cur.partial)
		}
		return kc
	}
	return nil
}

func isOut(cur cursor) bool {
	return cur.token.skip_to_bracket_pair()
}

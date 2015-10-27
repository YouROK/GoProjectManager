package codecomplete

import (
	"sort"
	"strings"
)

type Candidate struct {
	Class string `json:"class"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Perc  int    `json:"omitempty"`
}

type Candidates []Candidate

func NewCandidate(c, n, t string) Candidate {
	return Candidate{Class: c, Name: n, Type: t}
}

func (c Candidates) Add(cand Candidate) Candidates {
	return append(c, cand)
}

func (c Candidates) Adds(cand ...Candidate) Candidates {
	return append(c, cand...)
}

func (c Candidates) AddStr(cand string) Candidates {
	cc := Candidate{}
	cc.Name = cand
	return append(c, cc)
}

func (c Candidates) AddStrs(cand ...string) Candidates {
	cs := Candidates{}
	for _, cc := range cand {
		cs = cs.AddStr(cc)
	}
	return append(c, cs...)
}

func (c Candidates) Len() int {
	return len(c)
}

func (c Candidates) Less(i, j int) bool {
	return c[i].Perc > c[j].Perc
}

func (c Candidates) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Candidates) Sort() {
	sort.Sort(c)
}

func (c Candidates) SetPercent(word string) {
	for i := 0; i < len(c); i++ {
		c[i].Perc = wordCompare(word, c[i].Name)
	}
}

func getKeywordsCandidates() Candidates {
	ret := Candidates{}
	for _, k := range append(go_keywords, append(go_out_keywords, go_generic_types...)...) {
		ret = append(ret, Candidate{Class: "keyword", Type: "", Name: k})
	}
	return ret
}

func wordCompare(word, origWord string) int {
	word = strings.ToLower(word)
	origWord = strings.ToLower(origWord)

	perc := percentLetter([]rune(word), []rune(origWord))
	if perc < 50 {
		perc = Levenshtein(word, origWord)
	}
	return perc
}

func percentLetter(findLetters, origLetters []rune) int {
	letter := 0
	pos := 0
	for i, r := range findLetters {
		for j := pos; j < len(origLetters); j++ {
			if r == origLetters[j] {
				letter++
				if i == j {
					letter++
				}
				pos = j + 1
				break
			}
		}
		if pos >= len(origLetters) {
			break
		}
	}
	return len(origLetters) * letter / 200
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func Levenshtein(test, orig string) int {
	la := len(test)
	lb := len(orig)
	d := make([]int, la+1)
	var lastdiag, olddiag, temp int

	for i := 1; i <= la; i++ {
		d[i] = i
	}
	for i := 1; i <= lb; i++ {
		d[0] = i
		lastdiag = i - 1
		for j := 1; j <= la; j++ {
			olddiag = d[j]
			min := d[j] + 1
			if (d[j-1] + 1) < min {
				min = d[j-1] + 1
			}
			if test[j-1] == orig[i-1] {
				temp = 0
			} else {
				temp = 1
			}
			if (lastdiag + temp) < min {
				min = lastdiag + temp
			}
			d[j] = min
			lastdiag = olddiag
		}
	}
	return 100 - (120 * d[la] / lb)
}

var (
	go_out_keywords = []string{
		"func",
		"defer",
		"type",
		"struct",
		"interface",
		"import",
		"package",
	}

	go_keywords = []string{
		"return",
		"range",
		"select",
		"break",
		"continue",
		"default",
		"case",
		"go",
		"map",
		"chan",
		"else",
		"switch",
		"const",
		"var",
		"fallthrough",
		"if",
		"for",
		"goto",
	}

	go_generic_types = []string{
		"bool",
		"uint",
		"int",
		"uintptr",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
		"complex64",
		"complex128",
		"byte",
		"rune",
		"string",
	}
)

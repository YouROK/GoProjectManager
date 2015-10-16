package autocomplete

import (
	"sort"
)

type Candidate struct {
	Name string
	Perc int
}

type Candidates []Candidate

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
		c[i].Perc = WordCompare(word, c[i].Name)
	}
}

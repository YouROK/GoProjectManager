package autocomplete

import (
	"strings"
)

var (
	go_keywords = []string{
		"func",
		"defer",
		"type",
		"struct",
		"interface",
		"import",
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
		"package",
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

func WordCompare(word, origWord string) int {
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

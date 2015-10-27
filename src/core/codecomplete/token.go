package codecomplete

import (
	//	"go/ast"
	//	"go/parser"
	"go/scanner"
	"go/token"
)

type tokenItem struct {
	off int
	tok token.Token
	lit string
}

type tokenIterator struct {
	tokens []tokenItem
	index  int
}

func NewToken(src string, cursor int) tokenIterator {
	src = src[:cursor] + ";" + src[cursor:]
	tokens := make([]tokenItem, 0, 1000)
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, []byte(src), nil, 0)
	index := 0
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		off := fset.Position(pos).Offset

		tokens = append(tokens, tokenItem{
			off: off,
			tok: tok,
			lit: lit,
		})
		if cursor > off {
			index++
		}
	}
	return tokenIterator{
		tokens: tokens,
		index:  index,
	}
}

func (i tokenItem) Literal() string {
	if i.tok.IsLiteral() {
		return i.lit
	} else {
		return i.tok.String()
	}
	return ""
}

func (this *tokenIterator) token() tokenItem {
	return this.tokens[this.index]
}

func (this *tokenIterator) previous_token() bool {
	if this.index <= 0 {
		return false
	}
	this.index--
	return true
}

func (this *tokenIterator) next_token() bool {
	if this.index >= len(this.tokens)-1 {
		return false
	}
	this.index++
	return true
}

var g_bracket_pairs = map[token.Token]token.Token{
	token.RPAREN: token.LPAREN,
	token.RBRACK: token.LBRACK,
}

// when the cursor is at the ')' or ']', move the cursor to an opposite bracket
// pair, this functions takes inner bracker pairs into account
func (this *tokenIterator) skip_to_bracket_pair() bool {
	right := this.token().tok
	left := g_bracket_pairs[right]
	return this.skip_to_left_bracket(left, right)
}

func (this *tokenIterator) skip_to_left_bracket(left, right token.Token) bool {
	// TODO: Make this functin recursive.
	if this.token().tok == left {
		return true
	}
	balance := 1
	for balance != 0 {
		this.previous_token()
		if this.index == 0 {
			return false
		}
		switch this.token().tok {
		case right:
			balance++
		case left:
			balance--
		}
	}
	return true
}

// Move the cursor to the open brace of the current block, taking inner blocks
// into account.
func (this *tokenIterator) skip_to_open_brace() bool {
	return this.skip_to_left_bracket(token.LBRACE, token.RBRACE)
}

// try_extract_struct_init_expr tries to match the current cursor position as being inside a struct
// initialization expression of the form:
// &X{
// 	Xa: 1,
// 	Xb: 2,
// }
// Nested struct initialization expressions are handled correctly.
func (this *tokenIterator) try_extract_struct_init_expr() []byte {
	for this.index >= 0 {
		if !this.skip_to_open_brace() {
			return nil
		}

		if !this.previous_token() {
			return nil
		}

		return []byte(this.token().Literal())
	}
	return nil
}

// starting from the end of the 'file', move backwards and return a slice of a
// valid Go expression
func (this *tokenIterator) extract_go_expr() []byte {
	// TODO: Make this function recursive.
	orig := this.index

	// prev always contains the type of the previously scanned token (initialized with the token
	// right under the cursor). This is the token to the *right* of the current one.
	prev := this.token().tok
loop:
	for {
		this.previous_token()
		if this.index == 0 {
			return make_expr(this.tokens[:orig])
		}
		t := this.token().tok
		switch t {
		case token.PERIOD:
			if prev != token.IDENT {
				// Not ".ident".
				break loop
			}
		case token.IDENT:
			if prev == token.IDENT {
				// "ident ident".
				break loop
			}
		case token.RPAREN, token.RBRACK:
			if prev == token.IDENT {
				// ")ident" or "]ident".
				break loop
			}
			this.skip_to_bracket_pair()
		default:
			break loop
		}
		prev = t
	}
	exprT := this.tokens[this.index+1 : orig]
	return make_expr(exprT)
}

// Given a slice of tokenItem, reassembles them into the original literal expression.
func make_expr(tokens []tokenItem) []byte {
	e := ""
	for _, t := range tokens {
		e += t.Literal()
	}
	return []byte(e)
}

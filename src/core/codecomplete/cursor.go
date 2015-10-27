package codecomplete

import (
	"go/token"
)

type cursor struct {
	partial string
	token   tokenIterator
	offset  int
	src     string
}

func getCursor(src string, off int) cursor {
	c := cursor{}
	c.token = NewToken(src, off)
	c.offset = off
	c.src = src
	if c.token.previous_token() {
		if c.token.token().tok == token.IDENT {
			c.partial = c.token.token().Literal()
		}
	}
	return c
}

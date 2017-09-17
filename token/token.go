package token

import (
	"fmt"
)

type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF
	TOKEN_LEFT_BRACKET
	TOKEN_RIGHT_BRACKET
	TOKEN_EQUAL_SIGN
	TOKEN_NEWLINE
	TOKEN_SECTION
	TOKEN_KEY
	TOKEN_VALUE
)

const (
	EOF rune = 0
	LEFT_BRACKET string = "["
	RIGHT_BRACKET string = "]"
	EQUAL_SIGN string = "="
	NEWLINE string = "\n"
)

type Token struct {
	Type TokenType
	Value string
}

func (t *Token) String() string {

	switch t.Type {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_ERROR:
		return t.Value
	}

	return fmt.Sprintf("%q", t.Value)
}
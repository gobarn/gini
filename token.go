package gini

import (
	"fmt"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenLeftBracket
	tokenRightBracket
	tokenEqualSign
	tokenNewline
	tokenSection
	tokenKey
	tokenValue
)

const (
	eof rune = 0
	leftBracket string = "["
	rightBracket string = "]"
	equalSign string = "="
	newline string = "\n"
)

type Token struct {
	Type tokenType
	Value string
	Line int
}

func (t *Token) String() string {
	switch t.Type {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.Value
	}

	return fmt.Sprintf("%q", t.Value)
}
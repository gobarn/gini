package gini

import (
	"fmt"
	"unicode/utf8"
	"strings"
)

const (
	ErrorUnexpectedEOF string = "Unexpected EOF"
	ErrorMissingRightBracket string = "Missing a closing section bracket"
	ErrorMissingLabel string = "Your form is missing a label"
	ErrorMissingOpenParen string = "Missing open parenthesis after constraint"
	ErrorMissingClosingParen string = "Missing closing parenthesis after constraint"
	ErrorMissingType string = "Missing data type for this form element"
	ErrorInvalidConstrainti string = "Invalid constraint"
)

type StateFn func(*Lexer) StateFn

// Lexer basic spits back a stream of token on a channel. This lexer is
// based on lexing techniques outline by Rob Pike check out the link for 
// more details
// http://cuddle.googlecode.com/hg/talk/lex.html#landing-slide
type Lexer struct {
	Input string
	Line int
	Tokens chan Token
	State StateFn
	Start int
	Pos int
	Width int
}

// Create a new lexer
func lex(input string) *Lexer {
	return &Lexer{
		Input: input,
		Line: 1,
		State: LexTop,
		Tokens: make(chan Token, 10),
	}
}

func (lxr *Lexer) Next() Token {
	for {
		select {
		case token := <-lxr.Tokens:
			return token
		default: 
			lxr.State = lxr.State(lxr)
		}
	}
}

// Backup to the start of the last read token.
func (lxr *Lexer) Backup() {
	lxr.Pos -= lxr.Width
}

// Get a slice of current src from start
// to the current position
func (lxr *Lexer) Current() string {
	return lxr.Input[lxr.Start:lxr.Pos]
}

// Ignore current token
func (lxr *Lexer) Ignore() {
	lxr.Start = lxr.Pos
}

// Increment the position
func (lxr *Lexer) Inc() {
	lxr.Pos++
	if lxr.Pos >= utf8.RuneCountInString(lxr.Input) {
		lxr.Emit(tokenEOF)
	}
}
// Decrement the position
func (lxr *Lexer) Dec() {
	lxr.Pos--
}

func (lxr *Lexer) Emit(token tokenType) {
	lxr.Tokens <- Token{
		Type: token,
		Value: lxr.Input[lxr.Start:lxr.Pos],
		Line: lxr.Line,
	}
	lxr.Start = lxr.Pos
}

// Token as error
func (lxr *Lexer) ErrorF(format string, args ...interface{}) StateFn {
	lxr.Tokens <- Token {
		Type: tokenError,
		Value: fmt.Sprintf(format, args...),
	}

	return nil
}

func LexTop(lxr *Lexer) StateFn {
	isWhitespace()

	if strings.HasPrefix(InputToEnd(), leftBracket) {
		return LexLeftBracket
	} else {
		return LexKey
	}
}

func LexEqualSign(lxr *Lexer) StateFn {
	lxr.Pos += len(equalSign)
	lxr.Emit(tokenEqualSign)
	return LexValue
}

func LexKey(lxr *Lexer) StateFn {
	for {
		if strings.HasPrefix(InputToEnd(), equalSign) {
			lxr.Emit(tokenKey)
			return LexEqualSign
		}
		lxr.Inc()
		if lxr.IsEOF() {
			return lxr.Errorf(ErrorUnexpectedEOF)
		}
	}
}

func LexLeftBracket(lxr *Lexer) StateFn {
	lxr.Pos += len(leftBracket)
	lxr.Emit(tokenLeftBracket)
	return LexSection
}

func LexRightBracket(lxr *Lexer) StateFn {
	lxr.Pos += len(rightBracket)
	lxr.Emit(tokeRightBracket)
	return LexBegin
}

func LexSection(lxr *Lexer) StateFn {
	for {
		if lxr.IsEOF() {
			return lxr.Errorf(ErrorMissingRightBracket)
		}
		if strings.HasPrefix(InputToEnd(), rightBracket) {
			lxr.Emit(tokenSection)
			return LexRightBracket
		}
		lxr.Inc()
	}
}

func LexValue(lxr *Lexer) StateFn {
	for {
		if strings.HasPrefix(InputToEnd(), newline) {
			lxr.Emit(tokenValue)
			return LexBegin
		}
		lxr.Inc()
		if lxr.IsEOF() {
			return lxr.Errorf(ErrorMissingRightBracket)
		}
	}
}

func isWhitespace(t rune) bool {
	return t == ' ' || t == '\t'
}


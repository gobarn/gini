package lexer

import (
	"github.com/gobarn/gini/token"
)

type LexFn func(*Lexer) LexFn

// Lexer basic spits back a stream of token on a channel. This lexer is
// based on lexing techniques outline by Rob Pike check out the link for 
// more details
// http://cuddle.googlecode.com/hg/talk/lex.html#landing-slide
type Lexer struct {
	Name string
	Src string
	Tokens chan token.Token
	State LexFn
	Start int
	Pos int
	Width int
}

func Lex(name string, src string) *Lexer {
	return &Lexer {
		Name: name,
		Src: src,
		State: LexBegin,
		Tokens: make(chan Token, 5),
	}
}

func LexBegin(lxr *Lexer) LexFn {
	SkipWhitespace()

	if strings.HasPrefix(InputToEnd(), token.LEFT_BRACKET) {
		return LexLeftBracket
	} else {
		return LexKey
	}
}

func LexEqualSign(lxr *Lexer) LexFn {
	lxr.Pos += len(token.EQUAL_SIGN)
	lxr.Emit(token.TOKEN_EQUAL_SIGN)
	return LexValue
}

func LexKey(lxr *Lexer) LexFn {
	for {
		if strings.HasPrefix(InputToEnd(), token.EQUAL_SIGN) {
			lxr.Emit(token.TOKEN_KEY)
			return LexEqualSign
		}
		lxr.Inc()
		if lxr.IsEOF() {
			return lxr.Errorf(errors.UNEXPECTED_EOF)
		}
	}
}

func LexLeftBracket(lxr *Lexer) LexFn {
	lxr.Pos += len(token.LEFT_BRACKET)
	lxr.Emit(token.TOKEN_LEFT_BRACKET)
	return LexSection
}

func LexRightBracket(lxr *Lexer) LexFn {
	lxr.Pos += len(token.RIGHT_BRACKET)
	lxr.Emit(token.TOKEN_RIGHT_BRACKET)
	return LexBegin
}

func LexSection(lxr *Lexer) LexFn {
	for {
		if lxr.IsEOF() {
			return lxr.Errorf(errors.MISSING_RIGHT_BRACKET)
		}
		if strings.HasPrefix(InputToEnd(), token.RIGHT_BRACKET) {
			lxr.Emit(token.TOKEN_SECTION)
			return LexRightBracket
		}
		lxr.Inc()
	}
}

func LexValue(lxr *Lexer) LexFn {
	for {
		if strings.HasPrefix(InputToEnd(), token.NEWLINE) {
			lxr.Emit(token.TOKEN_VALUE)
			return LexBegin
		}
		lxr.Inc()
		if lxr.IsEOF() {
			return lxr.Errorf(errors.MISSING_RIGHT_BRACKET)
		}
	}
}



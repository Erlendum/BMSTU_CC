package lexer

import (
	"unicode"
)

const (
	TokenEOF = iota
	TokenERROR
	TokenIDENT
	TokenBOOL
	TokenLBRACE
	TokenRBRACE
	TokenSEMICOLON
	TokenASSIGN
	TokenNOT
	TokenAND
	TokenOR
	TokenEmpty
)

var (
	keywords = map[string]int{
		"true":  TokenBOOL,
		"false": TokenBOOL,
	}

	operators = map[string]int{
		"{": TokenLBRACE,
		"}": TokenRBRACE,
		";": TokenSEMICOLON,
		"=": TokenASSIGN,
		"~": TokenNOT,
		"&": TokenAND,
		"!": TokenOR,
	}
)

type Token struct {
	Type    int
	Literal string
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) Tokenize() []Token {
	tokens := []Token{}

	tok := Token{Type: TokenEmpty}
	for tok.Type != TokenEOF && tok.Type != TokenERROR {
		tok = l.NextToken()
		tokens = append(tokens, tok)
	}
	return tokens
}

func (l *Lexer) NextToken() Token {
	l.skipSpaces()

	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	if typ, ok := l.readOperator(); ok {
		return typ
	}

	ch := rune(l.input[l.pos])

	if unicode.IsLetter(ch) {
		return l.readIdentifierOrBool()
	}

	l.pos++
	return Token{Type: TokenERROR, Literal: string(ch)}
}

func (l *Lexer) readOperator() (Token, bool) {
	tok := l.input[l.pos : l.pos+1]
	if typ, ok := operators[tok]; ok {
		tok := Token{
			Type:    typ,
			Literal: tok,
		}
		l.pos++
		return tok, true
	}
	return Token{}, false
}

func (l *Lexer) skipSpaces() {
	if l.pos >= len(l.input) {
		return
	}

	for ; l.pos < len(l.input); l.pos++ {
		ch := rune(l.input[l.pos])
		if !unicode.IsSpace(ch) {
			return
		}
	}
}

func (l *Lexer) readIdentifierOrBool() Token {
	start := l.pos
	for l.pos < len(l.input) {
		ch := rune(l.input[l.pos])
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			break
		}
		l.pos++
	}

	literal := l.input[start:l.pos]

	if typ, ok := keywords[literal]; ok {
		return Token{Type: typ, Literal: literal}
	}
	return Token{Type: TokenIDENT, Literal: literal}
}

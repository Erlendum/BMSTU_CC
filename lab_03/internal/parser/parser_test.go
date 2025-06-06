package parser

import (
	"testing"

	"github.com/Erlendum/BMSTU_CC/lab_03/internal/lexer"
	"github.com/stretchr/testify/assert"
)

func token(t int, lit string) lexer.Token {
	return lexer.Token{
		Type:    t,
		Literal: lit,
	}
}

func TestParsePrimaryBool(t *testing.T) {
	tokens := []lexer.Token{
		token(lexer.TokenBOOL, "true"),
	}
	p := NewParser(tokens)
	_, ok := p.Parse()
	assert.False(t, ok)
}

func TestParseSimpleAssign(t *testing.T) {
	tokens := []lexer.Token{
		token(lexer.TokenLBRACE, "{"),
		token(lexer.TokenIDENT, "a"),
		token(lexer.TokenASSIGN, "="),
		token(lexer.TokenBOOL, "true"),
		token(lexer.TokenRBRACE, "}"),
	}
	p := NewParser(tokens)
	node, ok := p.Parse()
	assert.True(t, ok)
	assert.Equal(t, "Program", node.Type)
}

func TestParseWithOr(t *testing.T) {
	tokens := []lexer.Token{
		token(lexer.TokenLBRACE, "{"),
		token(lexer.TokenIDENT, "a"),
		token(lexer.TokenASSIGN, "="),
		token(lexer.TokenBOOL, "true"),
		token(lexer.TokenOR, "!"),
		token(lexer.TokenBOOL, "false"),
		token(lexer.TokenRBRACE, "}"),
	}
	p := NewParser(tokens)
	node, ok := p.Parse()
	assert.True(t, ok)
	assert.Equal(t, "Program", node.Type)
}

func TestParseMultipleStatements(t *testing.T) {
	tokens := []lexer.Token{
		token(lexer.TokenLBRACE, "{"),
		token(lexer.TokenIDENT, "a"),
		token(lexer.TokenASSIGN, "="),
		token(lexer.TokenBOOL, "true"),
		token(lexer.TokenSEMICOLON, ";"),
		token(lexer.TokenIDENT, "b"),
		token(lexer.TokenASSIGN, "="),
		token(lexer.TokenBOOL, "false"),
		token(lexer.TokenRBRACE, "}"),
	}
	p := NewParser(tokens)
	node, ok := p.Parse()
	assert.True(t, ok)
	assert.Equal(t, "Program", node.Type)
	assert.Equal(t, "Block", node.Children[0].Type)
	assert.Equal(t, "StmtList", node.Children[0].Children[0].Type)
}

func TestParseNotExpression(t *testing.T) {
	tokens := []lexer.Token{
		token(lexer.TokenLBRACE, "{"),
		token(lexer.TokenIDENT, "x"),
		token(lexer.TokenASSIGN, "="),
		token(lexer.TokenNOT, "~"),
		token(lexer.TokenBOOL, "false"),
		token(lexer.TokenRBRACE, "}"),
	}
	p := NewParser(tokens)
	_, ok := p.Parse()
	assert.True(t, ok)
}

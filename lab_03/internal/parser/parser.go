package parser

import (
	"fmt"
	"strings"

	"github.com/Erlendum/BMSTU_CC/lab_03/internal/lexer"
)

type ASTNode struct {
	Type     string
	Value    string
	Children []*ASTNode
	Token    lexer.Token
}

type Parser struct {
	tokens []lexer.Token
	pos    int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) CurrentToken() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) expect(tokenType int) (lexer.Token, bool) {
	tok := p.CurrentToken()
	if tok.Type == tokenType {
		p.incPos()
		return tok, true
	}
	return tok, false
}

func (p *Parser) parsePrimary() (*ASTNode, bool) {
	tok := p.CurrentToken()
	switch tok.Type {
	case lexer.TokenBOOL:
		p.incPos()
		return &ASTNode{
			Type:  "Primary",
			Value: tok.Literal,
			Token: tok,
		}, true
	case lexer.TokenIDENT:
		p.incPos()
		return &ASTNode{
			Type:  "Primary",
			Value: tok.Literal,
			Token: tok,
		}, true
	default:
		return nil, false
	}
}

func (p *Parser) parseSecondary() (*ASTNode, bool) {
	if p.isCurrentTokenMatchType(lexer.TokenNOT) {
		primary, ok := p.parsePrimary()
		if !ok {
			return nil, false
		}
		return &ASTNode{
			Type:     "Secondary",
			Children: []*ASTNode{primary},
			Value:    "~",
		}, true
	}
	return p.parsePrimary()
}

func (p *Parser) parseMonom() (*ASTNode, bool) {
	left, ok := p.parseSecondary()
	if !ok {
		return nil, false
	}

	for p.CurrentToken().Type == lexer.TokenAND {
		op := p.CurrentToken()
		p.incPos()
		right, ok := p.parseSecondary()
		if !ok {
			return nil, false
		}
		left = &ASTNode{
			Type:     "Monom",
			Children: []*ASTNode{left, right},
			Value:    op.Literal,
			Token:    op,
		}
	}
	return left, true
}

func (p *Parser) parseLogicalExpr() (*ASTNode, bool) {
	left, ok := p.parseMonom()
	if !ok {
		return nil, false
	}

	for p.CurrentToken().Type == lexer.TokenOR {
		op := p.CurrentToken()
		p.incPos()
		right, ok := p.parseMonom()
		if !ok {
			return nil, false
		}
		left = &ASTNode{
			Type:     "LogicalExpr",
			Children: []*ASTNode{left, right},
			Value:    op.Literal,
			Token:    op,
		}
	}
	return left, true
}

func (p *Parser) parseExpression() (*ASTNode, bool) {
	return p.parseLogicalExpr()
}

func (p *Parser) parseOperator() (*ASTNode, bool) {
	idTok, ok := p.expect(lexer.TokenIDENT)
	if !ok {
		return nil, false
	}

	if _, ok := p.expect(lexer.TokenASSIGN); !ok {
		return nil, false
	}

	expr, ok := p.parseExpression()
	if !ok {
		return nil, false
	}

	return &ASTNode{
		Type: "Operator",
		Children: []*ASTNode{
			{Type: "Identifier", Value: idTok.Literal, Token: idTok},
			expr,
		},
		Token: idTok,
	}, true
}

func (p *Parser) parseOperatorTail() (*ASTNode, bool) {
	if !p.isCurrentTokenMatchType(lexer.TokenSEMICOLON) {
		return &ASTNode{
			Type: "OperatorTail",
		}, true
	}

	op, ok := p.parseOperator()
	if !ok {
		return nil, false
	}

	tail, ok := p.parseOperatorTail()
	if !ok {
		return nil, false
	}

	return &ASTNode{
		Type:     "OperatorTail",
		Children: []*ASTNode{op, tail},
	}, true
}

func (p *Parser) parseStmtList() (*ASTNode, bool) {
	first, ok := p.parseOperator()
	if !ok {
		return nil, false
	}

	tail, ok := p.parseOperatorTail()
	if !ok {
		return nil, false
	}

	return &ASTNode{
		Type:     "StmtList",
		Children: []*ASTNode{first, tail},
	}, true
}

func (p *Parser) parseBlock() (*ASTNode, bool) {
	if !p.isCurrentTokenMatchType(lexer.TokenLBRACE) {
		return nil, false
	}

	stmts, ok := p.parseStmtList()
	if !ok {
		return nil, false
	}

	if !p.isCurrentTokenMatchType(lexer.TokenRBRACE) {
		return nil, false
	}

	return &ASTNode{
		Type:     "Block",
		Children: []*ASTNode{stmts},
	}, true
}

func (p *Parser) parseProgram() (*ASTNode, bool) {
	block, ok := p.parseBlock()
	if !ok {
		return nil, false
	}

	return &ASTNode{
		Type:     "Program",
		Children: []*ASTNode{block},
	}, true
}

func (p *Parser) Parse() (*ASTNode, bool) {
	return p.parseProgram()
}

func (p *Parser) incPos() {
	p.pos++
}

func (p *Parser) isCurrentTokenMatchType(tokenType int) bool {
	if p.CurrentToken().Type == tokenType {
		p.incPos()
		return true
	}
	return false
}

func (n *ASTNode) ToDot() string {
	var builder strings.Builder
	builder.WriteString("digraph AST {\n")
	builder.WriteString("  node [shape=box, fontname=\"Courier\", fontsize=10];\n")
	builder.WriteString("  edge [fontname=\"Courier\", fontsize=10];\n\n")

	var nodeCounter int
	generateDOTNode(&builder, n, &nodeCounter)

	builder.WriteString("}\n")
	return builder.String()
}

func generateDOTNode(builder *strings.Builder, node *ASTNode, counter *int) int {
	if node == nil {
		return -1
	}

	currentID := *counter
	*counter++

	label := node.Type
	if node.Value != "" {
		label += fmt.Sprintf("\\n%s", node.Value)
	}

	builder.WriteString(fmt.Sprintf("  node%d [label=\"%s\"];\n", currentID, label))

	for _, child := range node.Children {
		childID := generateDOTNode(builder, child, counter)
		if childID >= 0 {
			builder.WriteString(fmt.Sprintf("  node%d -> node%d;\n", currentID, childID))
		}
	}

	return currentID
}

package parser

import (
	"fmt"
	"strings"

	"github.com/Erlendum/BMSTU_CC/lab_04/internal/lexer"
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

func (p *Parser) parseProgram() (*ASTNode, bool) {
	expr, ok := p.parseExpression()
	if !ok {
		return nil, false
	}

	if p.CurrentToken().Type != lexer.TokenEOF {
		return nil, false
	}

	return expr, true
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

func (n *ASTNode) ToReversePolishNotation() []string {
	var result []string
	n.traversePostOrder(&result)
	return result
}

func (n *ASTNode) traversePostOrder(result *[]string) {
	if n == nil {
		return
	}

	for _, child := range n.Children {
		child.traversePostOrder(result)
	}

	if n.Value != "" {
		*result = append(*result, n.Value)
	}
}

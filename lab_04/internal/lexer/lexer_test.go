package lexer

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "empty input",
			input: "",
			expected: []Token{
				{TokenEOF, ""},
			},
		},
		{
			name:  "single identifier",
			input: "x",
			expected: []Token{
				{TokenIDENT, "x"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "boolean values",
			input: "true false",
			expected: []Token{
				{TokenBOOL, "true"},
				{TokenBOOL, "false"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "logical operations",
			input: "a & b ! ~c",
			expected: []Token{
				{TokenIDENT, "a"},
				{TokenAND, "&"},
				{TokenIDENT, "b"},
				{TokenOR, "!"},
				{TokenNOT, "~"},
				{TokenIDENT, "c"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "assignment",
			input: "x = true",
			expected: []Token{
				{TokenIDENT, "x"},
				{TokenASSIGN, "="},
				{TokenBOOL, "true"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "block with semicolons",
			input: "{ x = true; y = false; }",
			expected: []Token{
				{TokenLBRACE, "{"},
				{TokenIDENT, "x"},
				{TokenASSIGN, "="},
				{TokenBOOL, "true"},
				{TokenSEMICOLON, ";"},
				{TokenIDENT, "y"},
				{TokenASSIGN, "="},
				{TokenBOOL, "false"},
				{TokenSEMICOLON, ";"},
				{TokenRBRACE, "}"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "whitespace handling",
			input: "  x      \n=\t     ~y  ",
			expected: []Token{
				{TokenIDENT, "x"},
				{TokenASSIGN, "="},
				{TokenNOT, "~"},
				{TokenIDENT, "y"},
				{TokenEOF, ""},
			},
		},
		{
			name:  "error token",
			input: "x @ y",
			expected: []Token{
				{TokenIDENT, "x"},
				{TokenERROR, "@"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d, expected: %v, got: %v", len(tt.expected), len(tokens), tt.expected, tokens)
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type || tok.Literal != tt.expected[i].Literal {
					t.Errorf("token %d mismatch: expected %v, got %v",
						i, tt.expected[i], tok)
				}
			}
		})
	}
}

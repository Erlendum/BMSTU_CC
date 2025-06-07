package grammar

import (
	"testing"
)

func TestEliminateLeftRecursion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Непосредственная левая рекурсия",
			input: `1
A
2
a b
2
A -> Aa
A -> b
A`,
			expected: `2
A A'
2
a b
3
A -> bA'
A' -> aA'
A' -> ε
A
`,
		},
		{
			name: "Косвенная левая рекурсия",
			input: `2
A B
2
a b
3
A -> Ba
B -> Ab
B -> c
A
`,
			expected: `3
A B B'
2
a b
4
A -> Ba
B -> cB'
B' -> abB'
B' -> ε
A
`},
		{
			name: "циклы (алгоритм не справляется)",
			input: `2
A B
1
b
3
A -> B
B -> A
B -> b
A
`,
			expected: `3
A B B'
1
b
4
A -> B
B -> bB'
B' -> B'
B' -> ε
A
`,
		},
		{
			name: "Многоуровневая косвенная левая рекурсия",
			input: `4
S A B C
4
a b c d
8
S -> Sa
S -> Ab
A -> Ac
A -> Bd
B -> d
B -> Sa
B -> Cc
C -> Sa
S
`,
			expected: `8
A A' B B' C C' S S'
4
a b c d
13
A -> BdA'
A' -> cA'
A' -> ε
B -> CcB'
B -> dB'
B' -> dA'bS'aB'
B' -> ε
C -> dB'dA'bS'aC'
C' -> cB'dA'bS'aC'
C' -> ε
S -> AbS'
S' -> aS'
S' -> ε
S
`,
		},
		{
			name: "Без рекурсии",
			input: `2
S A
2
a b
3
S -> Ab
A -> a
A -> b
S
`,
			expected: `2
A S
2
a b
3
A -> a
A -> b
S -> Ab
S
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGrammarFromString(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if err := g.EliminateLeftRecursion(); err != nil {
				t.Fatal(err)
			}

			result := g.ToString()

			if result != tt.expected {
				t.Errorf("expected: %s, actual: %s", tt.expected, result)
			}
		})
	}
}

func TestEliminateChainRules(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Простые цепные правила",
			input: `2
S A
2
a b
3
S -> A
A -> a
A -> b
S
`,
			expected: `2
A S
2
a b
4
A -> a
A -> b
S -> a
S -> b
S
`,
		},
		{
			name: "Многоуровневые цепные правила",
			input: `3
S A B
2
a b
5
S -> A
A -> B
B -> a
B -> b
A -> a
S
`,
			expected: `3
A B S
2
a b
6
A -> a
A -> b
B -> a
B -> b
S -> a
S -> b
S
`,
		},
		{
			name: "Нет цепных правил",
			input: `2
S A
2
a b
3
S -> ab
A -> a
A -> b
S
`,
			expected: `2
A S
2
a b
3
A -> a
A -> b
S -> ab
S
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGrammarFromString(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if err := g.EliminateChainRules(); err != nil {
				t.Fatal(err)
			}

			result := g.ToString()
			if result != tt.expected {
				t.Errorf("ecpected: %s, actual: %s", tt.expected, result)
			}
		})
	}
}

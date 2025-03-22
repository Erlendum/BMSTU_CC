package infixToPostix

import "testing"

func TestFillConcateneteChars(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"(ab)*c", "(a.b)*.c"},
		{"(a(b|d))*", "(a.(b|d))*"},
		{"a(bb)+c", "a.(b.b)+.c"},
		{"abc", "a.b.c"},
		{"((a.b.c))", "((a.b.c))"},
		{"(a(b|c)*d)*((ad)*c)", "(a.(b|c)*.d)*.((a.d)*.c)"},
		{"((0|1)(0|1)(0|1))*", "((0|1).(0|1).(0|1))*"},
	}

	for _, tt := range tests {
		actual := fillConcatenateChars(tt.input)
		if actual != tt.expected {
			t.Errorf("Input: %s, Expected: %s, Actual %s", tt.input, tt.expected, actual)
		}
	}
}

func TestTransform(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"(ab)*c", "ab.*c."},
		{"(a(b|d))*", "abd|.*"},
		{"a(bb)+c", "abb.+.c."},
		{"abc", "ab.c."},
		{"((a.b.c))", "ab.c."},
		{"a.(b.b)+.c", "abb.+.c."},
		{"(a(b|c)*d)*((ad)*c)", "abc|*.d.*ad.*c.."},
		{"((0|1).(0|1).(0|1))*", "01|01|.01|.*"},
	}

	for _, tt := range tests {
		actual := Transform(tt.input)
		if actual != tt.expected {
			t.Errorf("Input: %s, Expected: %s, Actual %s", tt.input, tt.expected, actual)
		}
	}
}

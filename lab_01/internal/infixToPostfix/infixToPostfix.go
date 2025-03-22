package infixToPostix

import "strings"

const (
	maxPriority = 4
)

func fillConcatenateChars(infix string) string {
	var result strings.Builder
	n := len(infix)

	for i := 0; i < n; i++ {
		result.WriteByte(infix[i])

		if i+1 < n && shouldAddConcatenateChar(infix[i], infix[i+1]) {
			result.WriteByte('.')
		}
	}

	return result.String()
}

func shouldAddConcatenateChar(a, b byte) bool {
	return (isLetterOrDigit(a) && isLetterOrDigit(b)) ||
		(isLetterOrDigit(a) && b == '(') ||
		(a == ')' && isLetterOrDigit(b)) ||
		((a == '*' || a == '+') && (isLetterOrDigit(b) || b == '(')) || (a == ')' && b == '(')
}

func isLetterOrDigit(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

var specialCharsPriorityMap = map[rune]int{
	'(': 1,
	'|': 2,
	'.': 3,
	'?': 4,
	'*': 4,
	'+': 4,
}

func priorityOf(r rune) int {
	if priority, ok := specialCharsPriorityMap[r]; ok {
		return priority
	}
	return maxPriority + 1
}

func Transform(infix string) string {
	infix = fillConcatenateChars(infix)

	postfix := []rune{}
	stack := []rune{}

	for _, r := range infix {
		switch r {
		case '(':
			stack = append(stack, r)
		case ')':
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		default:
			for len(stack) > 0 && priorityOf(stack[len(stack)-1]) >= priorityOf(r) {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, r)
		}

	}

	for len(stack) > 0 {
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return string(postfix)
}

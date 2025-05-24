package dfa

import (
	"testing"

	nfa_pkg "github.com/Erlendum/BMSTU_CC/lab_01/internal/nfa"
)

type testCase struct {
	name     string
	input    string
	expected expectedDFA
}

type expectedDFA struct {
	startStateID int
	states       map[int]stateInfo
}

type stateInfo struct {
	isFinal     bool
	transitions map[rune]int
}

func TestBuild(t *testing.T) {
	tests := []testCase{
		{
			input: "ab.",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: false,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
				},
			},
		},
		{
			input: "ab|",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: false,
						transitions: map[rune]int{
							'a': 1,
							'b': 2,
						},
					},
					1: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
					2: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
				},
			},
		},
		{
			input: "ab.*",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
				},
			},
		},
		{
			input: "ab.?",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
				},
			},
		},
		{
			input: "ab.+",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: false,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		nfa := nfa_pkg.Build(tt.input)
		dfa := Build(nfa)
		checkDFA(t, dfa, tt.expected)
	}
}

func TestMinimize(t *testing.T) {
	tests := []testCase{
		{
			name:  "ab.",
			input: "ab.",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {

						isFinal: false,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
				},
			},
		},
		{
			name:  "ab|",
			input: "ab|",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: false,
						transitions: map[rune]int{
							'a': 1,
							'b': 1,
						},
					},
					1: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
				},
			},
		},

		{
			name:  "ab.*",
			input: "ab.*",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 0,
						},
					},
				},
			},
		},
		{
			name:  "ab.?",
			input: "ab.?",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal:     true,
						transitions: map[rune]int{},
					},
				},
			},
		},
		{
			name:  "ab.+",
			input: "ab.+",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: false,
						transitions: map[rune]int{
							'a': 1,
						},
					},
					1: {
						isFinal: false,
						transitions: map[rune]int{
							'b': 2,
						},
					},
					2: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 1,
						},
					},
				},
			},
		},
		{
			name:  "ab|*",
			input: "ab|*",
			expected: expectedDFA{
				startStateID: 0,
				states: map[int]stateInfo{
					0: {
						isFinal: true,
						transitions: map[rune]int{
							'a': 0,
							'b': 0,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nfa := nfa_pkg.Build(tt.input)
			dfa := Build(nfa)
			minimizedDFA := dfa.Minimize()
			checkDFA(t, minimizedDFA, tt.expected)
		})
	}
}

func checkDFA(t *testing.T, dfa *DFA, expected expectedDFA) {
	if dfa.Start != expected.startStateID {
		t.Errorf("ожидалось начальное состояние %d, получено - %d", expected.startStateID, dfa.Start)
	}

	if len(dfa.States) != len(expected.states) {
		t.Errorf("ожидаемое количество состояний - %d, получено - %d", len(expected.states), len(dfa.States))
	}

	for stateID, expectedState := range expected.states {
		state, ok := dfa.States[stateID]
		if !ok {
			t.Errorf("ожидалось состояние с ID %d, но оно не найдено", stateID)
			continue
		}

		if state.IsFinal != expectedState.isFinal {
			t.Errorf("для состояния %d ожидалось isFinal=%v, получено isFinal=%v", stateID, expectedState.isFinal, state.IsFinal)
		}

		if len(state.Transitions) != len(expectedState.transitions) {
			t.Errorf("для состояния %d ожидалось %d переходов, получено %d", stateID, len(expectedState.transitions), len(state.Transitions))
		}

		for symbol, expectedNextStateID := range expectedState.transitions {
			nextStateID, ok := state.Transitions[symbol]
			if !ok {
				t.Errorf("для состояния %d ожидался переход по символу %c, но он не найден", stateID, symbol)
				continue
			}

			if nextStateID != expectedNextStateID {
				t.Errorf("для состояния %d ожидался переход по символу %c в состояние %d, получено состояние %d", stateID, symbol, expectedNextStateID, nextStateID)
			}
		}
	}
}

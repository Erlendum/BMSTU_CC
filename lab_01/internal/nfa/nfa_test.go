package nfa

import (
	"testing"
)

type testCase struct {
	input    string
	expected expectedNFA
}

type expectedNFA struct {
	startStateID int
	endStateID   int
	transitions  map[int]transMap
}

type transMap map[rune][]int

func TestBuild(t *testing.T) {
	tests := []testCase{
		{
			input: "ab.",
			expected: expectedNFA{
				startStateID: 0,
				endStateID:   3,
				transitions: map[int]transMap{
					0: {'a': {1}},
					1: {EPS: {2}},
					2: {'b': {3}},
				},
			},
		},
		{
			input: "ab|",
			expected: expectedNFA{
				startStateID: 4,
				endStateID:   5,
				transitions: map[int]transMap{
					4: {EPS: {0, 2}},
					0: {'a': {1}},
					1: {EPS: {5}},
					2: {'b': {3}},
					3: {EPS: {5}},
				},
			},
		},
		{
			input: "ab.*",
			expected: expectedNFA{
				startStateID: 4,
				endStateID:   5,
				transitions: map[int]transMap{
					4: {EPS: {0, 5}},
					0: {'a': {1}},
					1: {EPS: {2}},
					2: {'b': {3}},
					3: {EPS: {0, 5}},
				},
			},
		},
		{
			input: "ab.?",
			expected: expectedNFA{
				startStateID: 4,
				endStateID:   5,
				transitions: map[int]transMap{
					4: {EPS: {0, 5}},
					0: {'a': {1}},
					1: {EPS: {2}},
					2: {'b': {3}},
					3: {EPS: {5}},
				},
			},
		},
		{
			input: "ab.+",
			expected: expectedNFA{
				startStateID: 4,
				endStateID:   5,
				transitions: map[int]transMap{
					4: {EPS: {0}},
					0: {'a': {1}},
					1: {EPS: {2}},
					2: {'b': {3}},
					3: {EPS: {0, 5}},
				},
			},
		},
	}

	for _, tt := range tests {
		nfa := Build(tt.input)
		checkNFA(t, nfa, tt.expected)
	}
}

func checkNFA(t *testing.T, nfa *NFA, expected expectedNFA) {
	if nfa.Start.ID != expected.startStateID {
		t.Errorf("ожидалось начальное состояние %d, получено - %d", expected.startStateID, nfa.Start.ID)
	}

	if nfa.End.ID != expected.endStateID {
		t.Errorf("ожидалось конечное состояние -  %d, получено - %d", expected.endStateID, nfa.End.ID)
	}

	visited := make(map[int]bool)
	checkTransitions(t, nfa.Start, expected.transitions, visited)
}

func checkTransitions(t *testing.T, state *State, expected map[int]transMap, visited map[int]bool) {
	if visited[state.ID] {
		return
	}
	visited[state.ID] = true

	expectedTransitions, ok := expected[state.ID]
	if !ok && len(state.Transitions) == 0 {
		return
	}

	if len(expectedTransitions) != len(state.Transitions) {
		t.Errorf("ожидаемое количество переходов - %d, получено - %d", len(expectedTransitions), len(state.Transitions))
		return
	}

	for symbol, nextStates := range state.Transitions {
		expectedNextStates, ok := expectedTransitions[symbol]
		if !ok {
			t.Errorf("не ожидался переход по символу %c из состояния %d", symbol, state.ID)
			continue
		}

		if len(nextStates) != len(expectedNextStates) {
			t.Errorf("ожидаемо колво переходов по символу %c из состояния %d - %d, получено - %d",
				symbol, state.ID, len(expectedNextStates), len(nextStates))
			continue
		}

		for i, nextState := range nextStates {
			if nextState.ID != expectedNextStates[i] {
				t.Errorf("ожидаемый переход %d, актуал переход %d", expectedNextStates[i], nextState.ID)
			}
		}
	}

	for _, nextStates := range state.Transitions {
		for _, nextState := range nextStates {
			checkTransitions(t, nextState, expected, visited)
		}
	}
}

package dfa

import (
	"fmt"

	nfa_pkg "github.com/Erlendum/BMSTU_CC/internal/nfa"
)

type State struct {
	ID          int
	NFAStates   map[int]bool
	Transitions map[rune]int
	IsFinal     bool
}

type DFA struct {
	Start    int
	States   map[int]*State
	Alphabet []rune
}

func NewState(id int, nfaStates map[int]bool, isFinal bool) *State {
	return &State{
		ID:          id,
		NFAStates:   nfaStates,
		Transitions: make(map[rune]int),
		IsFinal:     isFinal,
	}
}

func Build(nfa *nfa_pkg.NFA) *DFA {
	alphabet := nfa.ExtractAlphabet()

	dfa := &DFA{
		States:   make(map[int]*State),
		Alphabet: alphabet,
	}

	startedStates := make(map[int]bool)
	for _, state := range nfa.StartStates {
		startedStates[state.ID] = true
	}

	startNFAStates := nfa.EpsilonClosure(startedStates)
	dfa.Start = 0
	dfa.States[0] = NewState(0, startNFAStates, nfa.IsFinalState(startNFAStates))

	queue := []int{0}
	processed := make(map[int]bool)

	stateID := 1

	for len(queue) > 0 {
		currentStateID := queue[0]
		queue = queue[1:]

		if processed[currentStateID] {
			continue
		}
		processed[currentStateID] = true

		currentState := dfa.States[currentStateID]

		for _, symbol := range alphabet {
			nextNFAStates := make(map[int]bool)
			for nfaStateID := range currentState.NFAStates {
				state := nfa.StateByID(nfaStateID)
				for _, nextState := range state.Transitions[symbol] {
					nextNFAStates[nextState.ID] = true
				}
			}

			nextNFAStates = nfa.EpsilonClosure(nextNFAStates)

			if len(nextNFAStates) == 0 {
				continue
			}

			found := false
			var nextStateID int
			for id, state := range dfa.States {
				if statesEqual(state.NFAStates, nextNFAStates) {
					found = true
					nextStateID = id
					break
				}
			}

			if !found {
				nextStateID = stateID
				dfa.States[nextStateID] = NewState(nextStateID, nextNFAStates, nfa.IsFinalState(nextNFAStates))
				queue = append(queue, nextStateID)
				stateID++
			}

			currentState.Transitions[symbol] = nextStateID
		}
	}

	return dfa
}

func statesEqual(a, b map[int]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for stateID := range a {
		if !b[stateID] {
			return false
		}
	}
	return true
}

func (dfa *DFA) ToGraphviz() string {
	graph := "digraph DFA {\n"
	graph += "  rankdir=LR;\n"
	graph += "  node [shape = circle];\n"

	graph += fmt.Sprintf("  start [shape = point];\n")
	graph += fmt.Sprintf("  start -> %d;\n", dfa.Start)

	for _, state := range dfa.States {
		if state.IsFinal {
			graph += fmt.Sprintf("  %d [shape = doublecircle];\n", state.ID)
		} else {
			graph += fmt.Sprintf("  %d [shape = circle];\n", state.ID)
		}
	}

	for _, state := range dfa.States {
		for symbol, nextStateID := range state.Transitions {
			graph += fmt.Sprintf("  %d -> %d [label=\"%c\"];\n", state.ID, nextStateID, symbol)
		}
	}

	graph += "}\n"
	return graph
}

func (dfa *DFA) Minimize() *DFA {
	invertedNFA := dfa.invert()
	intermediateDFA := Build(invertedNFA)
	invertedNFA2 := intermediateDFA.invert()
	minimizedDFA := Build(invertedNFA2)
	return minimizedDFA
}

func (dfa *DFA) invert() *nfa_pkg.NFA {
	stateMap := make(map[int]*nfa_pkg.State)
	startStates := make([]*nfa_pkg.State, 0)
	startRandomID := 0
	for id, state := range dfa.States {
		stateMap[id] = &nfa_pkg.State{ID: id, Transitions: map[rune][]*nfa_pkg.State{}}
		if state.IsFinal {
			startRandomID = state.ID
			startStates = append(startStates, stateMap[id])
		}
	}

	stateMap[dfa.Start].IsFinal = true

	for fromID, state := range dfa.States {
		for symbol, toID := range state.Transitions {
			if _, exists := stateMap[toID]; !exists {
				stateMap[toID] = &nfa_pkg.State{ID: toID, Transitions: map[rune][]*nfa_pkg.State{}}
			}
			stateMap[toID].Transitions[symbol] = append(stateMap[toID].Transitions[symbol], stateMap[fromID])
		}
	}

	nfa := &nfa_pkg.NFA{
		Start:       stateMap[startRandomID],
		End:         stateMap[dfa.Start],
		StartStates: startStates,
	}

	return nfa
}

func (dfa *DFA) SimulateDFA(input string) ([]string, bool) {
	var steps []string
	currentStateID := dfa.Start
	currentState := dfa.States[currentStateID]

	steps = append(steps, dfa.ToGraphvizWithHighlight(currentStateID, "Start"))

	for i, symbol := range input {
		if nextStateID, exists := currentState.Transitions[symbol]; exists {
			currentStateID = nextStateID
			currentState = dfa.States[currentStateID]
			steps = append(steps, dfa.ToGraphvizWithHighlight(currentStateID, fmt.Sprintf("Step %d: Symbol '%c'", i+1, symbol)))
		} else {
			steps = append(steps, dfa.ToGraphvizWithError(currentStateID, symbol))
			return steps, false
		}
	}

	isAccepted := currentState.IsFinal
	if isAccepted {
		steps = append(steps, dfa.ToGraphvizWithHighlight(currentStateID, "Accepted"))
	} else {
		steps = append(steps, dfa.ToGraphvizWithHighlight(currentStateID, "Rejected"))
	}

	return steps, isAccepted
}

func (dfa *DFA) ToGraphvizWithHighlight(currentStateID int, description string) string {
	graph := "digraph DFA {\n"
	graph += "  rankdir=LR;\n"
	graph += "  node [shape = circle];\n"

	graph += fmt.Sprintf("  start [shape = point];\n")
	graph += fmt.Sprintf("  start -> %d;\n", dfa.Start)

	for _, state := range dfa.States {
		if state.IsFinal {
			graph += fmt.Sprintf("  %d [shape = doublecircle];\n", state.ID)
		} else {
			graph += fmt.Sprintf("  %d [shape = circle];\n", state.ID)
		}
	}

	graph += fmt.Sprintf("  %d [color=red, fontcolor=red];\n", currentStateID)

	graph += fmt.Sprintf("  labelloc=\"t\";\n")
	graph += fmt.Sprintf("  label=\"%s\";\n", description)

	for _, state := range dfa.States {
		for symbol, nextStateID := range state.Transitions {
			graph += fmt.Sprintf("  %d -> %d [label=\"%c\"];\n", state.ID, nextStateID, symbol)
		}
	}

	graph += "}\n"
	return graph
}

func (dfa *DFA) ToGraphvizWithError(currentStateID int, symbol rune) string {
	graph := "digraph DFA {\n"
	graph += "  rankdir=LR;\n"
	graph += "  node [shape = circle];\n"

	graph += fmt.Sprintf("  start [shape = point];\n")
	graph += fmt.Sprintf("  start -> %d;\n", dfa.Start)

	for _, state := range dfa.States {
		if state.IsFinal {
			graph += fmt.Sprintf("  %d [shape = doublecircle];\n", state.ID)
		} else {
			graph += fmt.Sprintf("  %d [shape = circle];\n", state.ID)
		}
	}

	graph += fmt.Sprintf("  %d [color=red, fontcolor=red];\n", currentStateID)
	graph += fmt.Sprintf("  %d -> error [label=\"%c\"];\n", currentStateID, symbol)
	graph += "  error [shape=box, color=red, fontcolor=red];\n"

	graph += fmt.Sprintf("  labelloc=\"t\";\n")
	graph += fmt.Sprintf("  label=\"Error: No transition for symbol '%c'\";\n", symbol)

	for _, state := range dfa.States {
		for symbol, nextStateID := range state.Transitions {
			graph += fmt.Sprintf("  %d -> %d [label=\"%c\"];\n", state.ID, nextStateID, symbol)
		}
	}

	graph += "}\n"
	return graph
}

package nfa

import (
	"fmt"
	"sort"
)

const EPS = 'ε'

type State struct {
	ID          int
	Transitions map[rune][]*State
	IsFinal     bool
}

type NFA struct {
	Start       *State
	End         *State
	StartStates []*State // у NFA по постронию одно состояние start, впихиваю сюда массив для алгоритма Бржозовского, так как там после инверта мб несколько стартов
}

func (a *NFA) ExtractAlphabet() []rune {
	alphabetMap := make(map[rune]bool)

	var traverse func(state *State)
	traverse = func(state *State) {
		for symbol := range state.Transitions {
			if symbol != EPS {
				alphabetMap[symbol] = true
			}
		}
	}

	visited := make(map[int]bool)
	stack := []*State{a.Start}
	stack = append(stack, a.StartStates...)

	for len(stack) > 0 {
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if state == nil {
			continue
		}

		if visited[state.ID] {
			continue
		}
		visited[state.ID] = true

		traverse(state)

		for _, nextStates := range state.Transitions {
			for _, nextState := range nextStates {
				stack = append(stack, nextState)
			}
		}
	}

	alphabet := make([]rune, 0, len(alphabetMap))
	for symbol := range alphabetMap {
		alphabet = append(alphabet, symbol)
	}

	sort.Slice(alphabet, func(i, j int) bool {
		return alphabet[i] < alphabet[j]
	})

	return alphabet
}

func NewState(id int) *State {
	return &State{
		ID:          id,
		Transitions: make(map[rune][]*State),
	}
}

func New(start, end *State) *NFA {
	return &NFA{Start: start, End: end, StartStates: []*State{}}
}

func Build(postfix string) *NFA {
	stack := []*NFA{}
	stateID := 0

	for _, char := range postfix {
		switch char {
		case '.':
			nfa2 := stack[len(stack)-1]
			nfa1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			nfa1.End.Transitions[EPS] = append(nfa1.End.Transitions[EPS], nfa2.Start)

			stack = append(stack, New(nfa1.Start, nfa2.End))
		case '|':
			nfa2 := stack[len(stack)-1]
			nfa1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			start := NewState(stateID)
			stateID++
			end := NewState(stateID)
			stateID++

			start.Transitions[EPS] = append(start.Transitions[EPS], nfa1.Start, nfa2.Start)

			nfa1.End.Transitions[EPS] = append(nfa1.End.Transitions[EPS], end)
			nfa2.End.Transitions[EPS] = append(nfa2.End.Transitions[EPS], end)

			stack = append(stack, New(start, end))
		case '?':
			nfa := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			start := NewState(stateID)
			stateID++
			end := NewState(stateID)
			stateID++

			start.Transitions[EPS] = append(start.Transitions[EPS], nfa.Start, end)
			nfa.End.Transitions[EPS] = append(nfa.End.Transitions[EPS], end)

			stack = append(stack, New(start, end))
		case '*':
			nfa := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			start := NewState(stateID)
			stateID++
			end := NewState(stateID)
			stateID++

			start.Transitions[EPS] = append(start.Transitions[EPS], nfa.Start, end)
			nfa.End.Transitions[EPS] = append(nfa.End.Transitions[EPS], nfa.Start)
			nfa.End.Transitions[EPS] = append(nfa.End.Transitions[EPS], end)

			stack = append(stack, New(start, end))
		case '+':
			nfa := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			start := NewState(stateID)
			stateID++
			end := NewState(stateID)
			stateID++

			start.Transitions[EPS] = append(start.Transitions[EPS], nfa.Start)
			nfa.End.Transitions[EPS] = append(nfa.End.Transitions[EPS], nfa.Start)
			nfa.End.Transitions[EPS] = append(nfa.End.Transitions[EPS], end)

			stack = append(stack, New(start, end))
		default:
			start := NewState(stateID)
			stateID++
			end := NewState(stateID)
			stateID++

			start.Transitions[char] = append(start.Transitions[char], end)
			stack = append(stack, New(start, end))
		}
	}

	stack[0].StartStates = append(stack[0].StartStates, stack[0].Start)
	stack[0].End.IsFinal = true
	return stack[0]
}

func (a *NFA) ToGraphviz() string {
	graph := "digraph NFA {\n"
	graph += "  rankdir=LR;\n"
	graph += "  node [shape = circle];\n"

	graph += "  start [shape = point];\n"
	graph += fmt.Sprintf("  start -> %d;\n", a.Start.ID)

	graph += fmt.Sprintf("  %d [shape = doublecircle];\n", a.End.ID)

	visited := make(map[*State]bool)
	stack := []*State{a.Start}

	for len(stack) > 0 {
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[state] {
			continue
		}
		visited[state] = true

		graph += fmt.Sprintf("  %d [label=\"%d\"];\n", state.ID, state.ID)

		for char, nextStates := range state.Transitions {
			for _, nextState := range nextStates {
				graph += fmt.Sprintf("  %d -> %d [label=\"%c\"];\n", state.ID, nextState.ID, char)
				stack = append(stack, nextState)
			}
		}
	}

	graph += "}\n"
	return graph
}

func (a *NFA) StateByID(stateID int) *State {
	visited := make(map[int]bool)
	stack := []*State{a.Start}
	stack = append(stack, a.StartStates...)

	for len(stack) > 0 {
		state := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[state.ID] {
			continue
		}
		visited[state.ID] = true

		if state.ID == stateID {
			return state
		}

		for _, nextStates := range state.Transitions {
			for _, nextState := range nextStates {
				stack = append(stack, nextState)
			}
		}
	}

	return nil
}

func (a *NFA) IsFinalState(states map[int]bool) bool {
	for stateID := range states {
		state := a.StateByID(stateID)
		if state.IsFinal {
			return true
		}
	}
	return false
}

func (a *NFA) EpsilonClosure(states map[int]bool) map[int]bool {
	closure := make(map[int]bool)
	for stateID := range states {
		closure[stateID] = true
	}

	stack := make([]int, 0, len(states))
	for stateID := range states {
		stack = append(stack, stateID)
	}

	for len(stack) > 0 {
		currentStateID := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		state := a.StateByID(currentStateID)
		for _, nextState := range state.Transitions[EPS] {
			if !closure[nextState.ID] {
				closure[nextState.ID] = true
				stack = append(stack, nextState.ID)
			}
		}
	}

	return closure
}

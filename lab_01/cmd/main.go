package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Erlendum/BMSTU_CC/internal/dfa"
	infixToPostix "github.com/Erlendum/BMSTU_CC/internal/infixToPostfix"
	"github.com/Erlendum/BMSTU_CC/internal/nfa"
)

const (
	nfaFileName    = "nfa.dot"
	dfaFileName    = "dfa.dot"
	minDFAFileName = "min_dfa.dot"
	stepsDir       = "./steps"
)

func prepareStepsDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("ошибка удаления папки: %w", err)
	}

	err = os.Mkdir(dir, 0755)
	if err != nil {
		return fmt.Errorf("ошибка создания папки: %w", err)
	}

	return nil
}

func main() {
	mode := flag.String("mode", "nfa", "Режим работы (nfa, dfa, minDFA, modeling), по умолчанию будет nfa (построение НКА)")
	regex := flag.String("regex", "(ab)*c", "Регулярное выражение, по умолчанию будет (ab)*c")
	input := flag.String("input", "abc", "Входная строка для режима modeling, по умолчанию будет abc")
	flag.Parse()

	switch *mode {
	// case "test":
	// 	a := dfa.NewState(0, nil, false)
	// 	b := dfa.NewState(1, nil, false)
	// 	c := dfa.NewState(2, nil, false)
	// 	d := dfa.NewState(3, nil, false)
	// 	e := dfa.NewState(4, nil, true)
	// 	f := dfa.NewState(5, nil, true)

	// 	a.Transitions = map[rune]int{'0': 1, '1': 2}
	// 	b.Transitions = map[rune]int{'0': 4, '1': 5}
	// 	c.Transitions = map[rune]int{'0': 0, '1': 0}
	// 	d.Transitions = map[rune]int{'0': 5, '1': 4}
	// 	e.Transitions = map[rune]int{'0': 3, '1': 5}
	// 	f.Transitions = map[rune]int{'0': 3, '1': 4}
	// 	dfa := dfa.DFA{Start: 0, Alphabet: []rune{'0', '1'}, States: map[int]*dfa.State{0: a, 1: b, 2: c, 3: d, 4: e, 5: f}}
	// 	graph := dfa.ToGraphviz()
	// 	err := os.WriteFile(dfaFileName, []byte(graph), 0644)
	// 	if err != nil {
	// 		fmt.Println("ошибка при записи файла:", err)
	// 		return
	// 	}
	// 	fmt.Printf("DFA сохранен в файл: %s\n", dfaFileName)
	// 	dfa.Minimize()
	case "nfa":
		postfix := infixToPostix.Transform(*regex)

		graphvizNFA := nfa.Build(postfix).ToGraphviz()

		err := os.WriteFile(nfaFileName, []byte(graphvizNFA), 0644)
		if err != nil {
			fmt.Println("ошибка при записи файла:", err)
			return
		}
		fmt.Printf("NFA сохранен в файл: %s\n", nfaFileName)
	case "dfa":
		postfix := infixToPostix.Transform(*regex)
		graphvizDFA := dfa.Build(nfa.Build(postfix)).ToGraphviz()
		err := os.WriteFile(dfaFileName, []byte(graphvizDFA), 0644)
		if err != nil {
			fmt.Println("ошибка при записи файла:", err)
			return
		}
		fmt.Printf("DFA сохранен в файл: %s\n", dfaFileName)
	case "minDFA":
		postfix := infixToPostix.Transform(*regex)
		graphvizMinDFA := dfa.Build(nfa.Build(postfix)).Minimize().ToGraphviz()
		err := os.WriteFile(minDFAFileName, []byte(graphvizMinDFA), 0644)
		if err != nil {
			fmt.Println("ошибка при записи файла:", err)
			return
		}
		fmt.Printf("Min DFA сохранен в файл: %s\n", minDFAFileName)
	case "modeling":
		postfix := infixToPostix.Transform(*regex)
		minDFA := dfa.Build(nfa.Build(postfix)).Minimize()
		steps, accepted := minDFA.SimulateDFA(*input)

		err := prepareStepsDir(stepsDir)
		if err != nil {
			fmt.Printf("ошибка подготовки папки: %v\n", err)
			return
		}

		for i, step := range steps {
			filename := fmt.Sprintf(stepsDir+"/step_%d.dot", i+1)
			err := os.WriteFile(filename, []byte(step), 0644)
			if err != nil {
				fmt.Printf("ошибка при записи файла %s: %v\n", filename, err)
				return
			}
			fmt.Printf("Step %d сохранен как %s\n", i+1, filename)
		}

		if accepted {
			fmt.Printf("Строка %s допускается ДКА", *input)
		} else {
			fmt.Printf("Строка %s НЕ допускается ДКА", *input)
		}
	default:
		fmt.Println("Режим не поддерживается. Доступные режим: nfa, dfa, minDFA, modeling")
	}
}

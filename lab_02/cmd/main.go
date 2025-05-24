package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Erlendum/BMSTU_CC/lab_02/internal/grammar"
)

func main() {
	mode := flag.String("mode", "eliminate_left_recursion", "Режим работы (eliminate_left_recursion, eliminate_chain), по умолчанию будет eliminate_left_recursion (устранение левой рекурсии)")
	inputFile := flag.String("input_file", "./data/eliminate_left_recursion/in_01.txt", "Путь до входного файла")
	outputFile := flag.String("output_file", "./data/eliminate_left_recursion/out_01.txt", "Путь до выходного файла")
	flag.Parse()

	switch *mode {
	case "eliminate_left_recursion":
		inputBytes, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Println("Ошибка чтения файла:", err)
			return
		}
		inputString := string(inputBytes)
		g, err := grammar.NewGrammarFromString(inputString)
		if err != nil {
			fmt.Println("Ошибка построения грамматики:", err)
			return
		}

		err = g.EliminateLeftRecursion()
		if err != nil {
			fmt.Println("Ошибка устранения левой рекурсии:", err)
			return
		}

		output := g.ToString()

		err = os.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Printf("ошибка при записи файла %s: %v\n", *outputFile, err)
			return
		}
		fmt.Printf("Выходной файл сохранен как %s\n", *outputFile)
	case "eliminate_chain":
		inputBytes, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Println("Ошибка чтения файла:", err)
			return
		}
		inputString := string(inputBytes)
		g, err := grammar.NewGrammarFromString(inputString)
		if err != nil {
			fmt.Println("Ошибка построения грамматики:", err)
			return
		}

		err = g.EliminateChainRules()
		if err != nil {
			fmt.Println("Ошибка устранения левой рекурсии:", err)
			return
		}

		output := g.ToString()

		err = os.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Printf("ошибка при записи файла %s: %v\n", *outputFile, err)
			return
		}
		fmt.Printf("Выходной файл сохранен как %s\n", *outputFile)
	default:
		fmt.Println("Режим не поддерживается. Доступные режимы: eliminate_left_recursion, eliminate_chain")
	}
}

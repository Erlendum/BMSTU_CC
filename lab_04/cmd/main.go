package main

import (
	"fmt"
	"os"

	"github.com/Erlendum/BMSTU_CC/lab_04/internal/lexer"
	"github.com/Erlendum/BMSTU_CC/lab_04/internal/parser"
)

const (
	defaultOutFileName = "ast.dot"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("go run ./cmd/main.go [filepath]")
		os.Exit(1)
	}

	filename := os.Args[1]
	inputBytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("read file error:", err)
		os.Exit(1)
	}

	lex := lexer.NewLexer(string(inputBytes))
	tokens := lex.Tokenize()

	pars := parser.NewParser(tokens)
	ast, ok := pars.Parse()

	if !ok {
		tok := pars.CurrentToken()
		fmt.Printf("unexpected token: %s\n", tok.Literal)
	} else {
		out := ast.ToDot()
		err = os.WriteFile(defaultOutFileName, []byte(out), 0644)
		if err != nil {
			fmt.Printf("write file %s error: %v\n", defaultOutFileName, err)
			os.Exit(1)
		}
		fmt.Printf("ast file: %s\n", defaultOutFileName)

		fmt.Printf("RPN: %v", ast.ToReversePolishNotation())
	}
}

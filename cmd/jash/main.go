package main

import (
	"fmt"
	"os"

	"github.com/qwantuum/jash/pkg/evaluator"
	"github.com/qwantuum/jash/pkg/lexer"
	"github.com/qwantuum/jash/pkg/parser"
)

func main() {
	args := os.Args[1:]

	stacktrace := false
	filename := ""
	for _, a := range args {
		if a == "--stacktrace" {
			stacktrace = true
		} else if filename == "" {
			filename = a
		}
	}

	if filename == "" {
		fmt.Println("Usage: jash [--stacktrace] <file.jash>")
		os.Exit(1)
	}

	evaluator.StacktraceEnabled = stacktrace
	evaluator.InitJIT()

	source, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(source))
	tokens, errs := l.Tokenize()
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "Lexer error: %s\n", e)
		}
		os.Exit(1)
	}

	p := parser.New(tokens)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			fmt.Fprintf(os.Stderr, "Parser error: %s\n", e)
		}
		os.Exit(1)
	}

	env := evaluator.NewEnvironment()
	result := evaluator.Eval(program, env)

	if result != nil {
		switch r := result.(type) {
		case *evaluator.Error:
			if trace := evaluator.GetStackTrace(); trace != "" {
				fmt.Fprint(os.Stderr, trace)
			}
			evaluator.ClearCallStack()
			fmt.Fprintf(os.Stderr, "Runtime error: %s\n", r.Message)
			os.Exit(1)
		default:
			if result.Type() != evaluator.NULL_OBJ && result.Type() != evaluator.RETURN_OBJ {
				fmt.Println(result.Inspect())
			}
		}
	}
}

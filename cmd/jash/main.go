package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/qwantuum/jash/pkg/evaluator"
	"github.com/qwantuum/jash/pkg/lexer"
	"github.com/qwantuum/jash/pkg/parser"
)

func runFile(filename string, stacktrace bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Panic: %v\n", r)
			os.Exit(1)
		}
	}()

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

func runREPL() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "REPL panic: %v\n", r)
		}
	}()

	evaluator.InitJIT()
	env := evaluator.NewEnvironment()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Jash REPL")
	fmt.Println("Enter blank line to evaluate, or 'exit' to quit")

	var lines []string

	for {
		fmt.Print("  ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" {
			break
		}

		lines = append(lines, line)

		if strings.TrimSpace(line) == "" {
			input := strings.Join(lines, "\n")
			lines = nil

			input = strings.TrimSpace(input)
			if input == "" {
				continue
			}

			l := lexer.New(input)
			tokens, errs := l.Tokenize()
			if len(errs) > 0 {
				for _, e := range errs {
					fmt.Fprintln(os.Stderr, e)
				}
				continue
			}

			p := parser.New(tokens)
			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				for _, e := range p.Errors() {
					fmt.Fprintln(os.Stderr, e)
				}
				continue
			}

			result := evaluator.Eval(program, env)
			if result != nil && result.Type() != evaluator.NULL_OBJ && result.Type() != evaluator.RETURN_OBJ {
				switch r := result.(type) {
				case *evaluator.Error:
					fmt.Fprintf(os.Stderr, "Error: %s\n", r.Message)
				default:
					fmt.Println(result.Inspect())
				}
			}
		}
	}
}

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
		runREPL()
		return
	}

	runFile(filename, stacktrace)
}

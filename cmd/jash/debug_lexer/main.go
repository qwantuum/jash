package main

import (
	"fmt"
	"os"

	"github.com/qwantuum/jash/pkg/lexer"
)

func main() {
	fmt.Fprintln(os.Stderr, "=== DEBUG LEXER START ===")

	source := "def foo()\n    print(\"hi\")"
	l := lexer.New(source)
	fmt.Fprintln(os.Stderr, "before Tokenize()")
	tokens, errs := l.Tokenize()
	fmt.Fprintln(os.Stderr, "after Tokenize()")

	fmt.Fprintf(os.Stderr, "tokens=%d errors=%d\n", len(tokens), len(errs))
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "err: %s\n", e)
	}
	for _, t := range tokens {
		fmt.Fprintf(os.Stderr, "  %s %q\n", t.Type, t.Literal)
	}
	fmt.Fprintln(os.Stderr, "=== DONE ===")
}

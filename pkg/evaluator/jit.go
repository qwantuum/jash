package evaluator

import (
	"fmt"
	"os"

	"github.com/qwantuum/jash/pkg/ast"
)

var JITEnabled bool

type CompiledFunc struct {
	instructions []Instruction
	constants    []interface{}
}

type JITManager struct {
	callCounts map[string]int
	compiled   map[string]*CompiledFunc
	threshold  int
}

var GlobalJIT *JITManager

func InitJIT() {
	GlobalJIT = &JITManager{
		callCounts: make(map[string]int),
		compiled:   make(map[string]*CompiledFunc),
		threshold:  3,
	}
}

func (jm *JITManager) RecordCall(name string) {
	if !JITEnabled {
		return
	}
	jm.callCounts[name]++
}

func (jm *JITManager) IsCompiled(name string) bool {
	_, ok := jm.compiled[name]
	return ok
}

func (jm *JITManager) Compile(name string, params []*ast.Identifier, body *ast.BlockStatement) bool {
	if !JITEnabled {
		return false
	}

	c := newCompiler()
	for _, p := range params {
		c.getVarIndex(p.Value)
	}

	instructions, constants, err := c.Compile(body)
	if err != nil {
		return false
	}

	jm.compiled[name] = &CompiledFunc{
		instructions: instructions,
		constants:    constants,
	}

	if os.Getenv("JASH_JIT_DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, "[JIT] compiled '%s' (%d instructions)\n", name, len(instructions))
	}

	return true
}

func (jm *JITManager) Execute(name string, args []Object) Object {
	cf, ok := jm.compiled[name]
	if !ok {
		return &Error{Message: fmt.Sprintf("JIT: %s not compiled", name)}
	}

	globals := make([]Object, len(cf.instructions)+len(args)+10)
	for i, arg := range args {
		globals[i] = arg
	}

	vm := newVM(cf.instructions, cf.constants, globals)
	result := vm.Run()

	if result == nil {
		return NULL
	}
	return result
}

func (jm *JITManager) ShouldCompile(name string) bool {
	count, ok := jm.callCounts[name]
	if !ok {
		return false
	}
	if jm.IsCompiled(name) {
		return false
	}
	return count >= jm.threshold
}

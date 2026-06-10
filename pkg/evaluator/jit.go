package evaluator

import (
	"fmt"
	"os"

	"github.com/qwantuum/jash/pkg/ast"
)

type CompiledFunc struct {
	instructions []Instruction
	constants    []interface{}
	varNames     map[int]string
}

type JITManager struct {
	compiled map[string]*CompiledFunc
}

var GlobalJIT *JITManager

func InitJIT() {
	GlobalJIT = &JITManager{
		compiled: make(map[string]*CompiledFunc),
	}
}

func (jm *JITManager) IsCompiled(name string) bool {
	_, ok := jm.compiled[name]
	return ok
}

func (jm *JITManager) Compile(name string, params []*ast.Identifier, body *ast.BlockStatement) bool {
	c := newCompiler()
	for _, p := range params {
		c.getVarIndex(p.Value)
	}

	instructions, constants, varNames, err := c.Compile(body)
	if err != nil {
		return false
	}

	jm.compiled[name] = &CompiledFunc{
		instructions: instructions,
		constants:    constants,
		varNames:     varNames,
	}

	if os.Getenv("JASH_JIT_DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, "[JIT] compiled '%s' (%d instructions)\n", name, len(instructions))
	}

	return true
}

func (jm *JITManager) Execute(name string, args []Object, env *Environment) Object {
	cf, ok := jm.compiled[name]
	if !ok {
		return &Error{Message: fmt.Sprintf("JIT: %s not compiled", name)}
	}

	globals := make([]Object, len(cf.instructions)+len(args)+10)
	for i, arg := range args {
		globals[i] = arg
	}

	vm := newVM(cf.instructions, cf.constants, cf.varNames, globals, env)
	result := vm.Run()

	if result == nil {
		return NULL
	}
	return result
}

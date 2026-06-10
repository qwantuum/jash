package evaluator

import "github.com/qwantuum/jash/pkg/ast"

type Opcode byte

const (
	OpConstant Opcode = iota
	OpLoad
	OpStore
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpEq
	OpNeq
	OpLt
	OpGt
	OpLte
	OpGte
	OpAnd
	OpOr
	OpNot
	OpMinus
	OpCall
	OpReturn
	OpReturnVal
	OpJump
	OpJumpIfFalse
	OpPop
	OpNewArray
	OpNewObject
	OpSetMember
	OpGetMember
	OpNull
	OpTrue
	OpFalse
	OpDefFunc
)

type Instruction struct {
	Op    Opcode
	Arg   int
	Arg2  int
	Const interface{}
}

type funcDef struct {
	Name   string
	Params []*ast.Identifier
	Body   *ast.BlockStatement
}

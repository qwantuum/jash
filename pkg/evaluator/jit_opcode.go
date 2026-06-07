package evaluator

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
)

type Instruction struct {
	Op    Opcode
	Arg   int
	Arg2  int
	Const interface{}
}

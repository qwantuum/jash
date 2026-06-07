package evaluator

import (
	"fmt"

	"github.com/qwantuum/jash/pkg/ast"
)

type compiler struct {
	instructions []Instruction
	constants    []interface{}
	varIndex     map[string]int
	varCount     int
}

func newCompiler() *compiler {
	return &compiler{
		varIndex: make(map[string]int),
	}
}

func (c *compiler) Compile(node ast.Node) ([]Instruction, []interface{}, error) {
	c.instructions = nil
	c.constants = nil
	c.varIndex = make(map[string]int)
	c.varCount = 0

	switch n := node.(type) {
	case *ast.BlockStatement:
		c.compileBlock(n)
	default:
		c.compileNode(node)
	}

	c.emit(OpReturn)

	return c.instructions, c.constants, nil
}

func (c *compiler) compileBlock(block *ast.BlockStatement) {
	for _, stmt := range block.Statements {
		c.compileStatement(stmt)
	}
}

func (c *compiler) compileStatement(stmt ast.Statement) {
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		c.compileNode(s.Expression)
		c.emit(OpPop)
	case *ast.ReturnStatement:
		if s.Value != nil {
			c.compileNode(s.Value)
			c.emit(OpReturnVal)
		} else {
			c.emit(OpReturn)
		}
	case *ast.AssignStatement:
		c.compileNode(s.Value)
		idx := c.getVarIndex(s.Name.Value)
		c.emit(OpStore, idx)
	case *ast.IfStatement:
		c.compileIf(s)
	case *ast.ForStatement:
		c.compileFor(s)
	case *ast.WhileStatement:
		c.compileWhile(s)
	case *ast.FunctionStatement:
	}
}

func (c *compiler) compileIf(stmt *ast.IfStatement) {
	c.compileNode(stmt.Condition)
	jumpElse := len(c.instructions)
	c.emit(OpJumpIfFalse, 0)

	c.compileBlock(stmt.Body)

	if stmt.ElseBody != nil {
		jumpEnd := len(c.instructions)
		c.emit(OpJump, 0)

		c.instructions[jumpElse].Arg = len(c.instructions)
		c.compileBlock(stmt.ElseBody)

		c.instructions[jumpEnd].Arg = len(c.instructions)
	} else {
		c.instructions[jumpElse].Arg = len(c.instructions)
	}
}

func (c *compiler) compileFor(stmt *ast.ForStatement) {
	start := len(c.instructions)
	c.compileNode(stmt.Iterable)
	c.emit(OpLoad, c.getVarIndex(stmt.Variable.Value))
	c.compileBlock(stmt.Body)
	c.emit(OpJump, start)
}

func (c *compiler) compileWhile(stmt *ast.WhileStatement) {
	start := len(c.instructions)
	c.compileNode(stmt.Condition)
	jumpExit := len(c.instructions)
	c.emit(OpJumpIfFalse, 0)
	c.compileBlock(stmt.Body)
	c.emit(OpJump, start)
	c.instructions[jumpExit].Arg = len(c.instructions)
}

func (c *compiler) compileNode(node ast.Node) {
	switch n := node.(type) {
	case *ast.Identifier:
		idx, ok := c.varIndex[n.Value]
		if !ok {
			idx = c.getVarIndex(n.Value)
		}
		c.emit(OpLoad, idx)
	case *ast.NumberLiteral:
		idx := c.addConstant(n.Value)
		c.emit(OpConstant, idx)
	case *ast.StringLiteral:
		idx := c.addConstant(n.Value)
		c.emit(OpConstant, idx)
	case *ast.BooleanLiteral:
		if n.Value {
			c.emit(OpTrue)
		} else {
			c.emit(OpFalse)
		}
	case *ast.NullLiteral:
		c.emit(OpNull)
	case *ast.InfixExpression:
		c.compileNode(n.Left)
		c.compileNode(n.Right)
		c.emit(c.infixOp(n.Operator))
	case *ast.PrefixExpression:
		c.compileNode(n.Right)
		switch n.Operator {
		case "-":
			c.emit(OpMinus)
		case "not":
			c.emit(OpNot)
		}
	case *ast.CallExpression:
		for i := len(n.Arguments) - 1; i >= 0; i-- {
			c.compileNode(n.Arguments[i])
		}
		c.compileNode(n.Function)
		c.emit(OpCall, len(n.Arguments))
	case *ast.MemberAccess:
		c.compileNode(n.Object)
		idx := c.addConstant(n.Member.Value)
		c.emit(OpGetMember, idx)
	case *ast.JSONArray:
		c.emit(OpNewArray)
		for _, elem := range n.Elements {
			c.compileNode(elem)
			c.emit(OpSetMember)
		}
	case *ast.JSONObject:
		c.emit(OpNewObject)
		for k, v := range n.Pairs {
			idx := c.addConstant(k)
			c.compileNode(v)
			c.emit(OpSetMember, idx)
		}
	default:
		panic(fmt.Sprintf("JIT: unsupported node %T", n))
	}
}

func (c *compiler) emit(op Opcode, args ...int) {
	inst := Instruction{Op: op}
	if len(args) > 0 {
		inst.Arg = args[0]
	}
	if len(args) > 1 {
		inst.Arg2 = args[1]
	}
	c.instructions = append(c.instructions, inst)
}

func (c *compiler) addConstant(val interface{}) int {
	c.constants = append(c.constants, val)
	return len(c.constants) - 1
}

func (c *compiler) getVarIndex(name string) int {
	if idx, ok := c.varIndex[name]; ok {
		return idx
	}
	idx := c.varCount
	c.varCount++
	c.varIndex[name] = idx
	return idx
}

func (c *compiler) infixOp(op string) Opcode {
	switch op {
	case "+":
		return OpAdd
	case "-":
		return OpSub
	case "*":
		return OpMul
	case "/":
		return OpDiv
	case "==":
		return OpEq
	case "!=":
		return OpNeq
	case "<":
		return OpLt
	case ">":
		return OpGt
	case "<=":
		return OpLte
	case ">=":
		return OpGte
	case "and":
		return OpAnd
	case "or":
		return OpOr
	default:
		panic("unknown operator: " + op)
	}
}

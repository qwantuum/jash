package evaluator

import (
	"fmt"
)

type vm struct {
	instructions []Instruction
	constants    []interface{}
	varNames     map[int]string
	ip           int
	stack        []Object
	sp           int
	globals      []Object
	env          *Environment
}

func newVM(instructions []Instruction, constants []interface{}, varNames map[int]string, globals []Object, env *Environment) *vm {
	return &vm{
		instructions: instructions,
		constants:    constants,
		varNames:     varNames,
		stack:        make([]Object, 1024),
		sp:           0,
		globals:      globals,
		env:          env,
	}
}

func (vm *vm) Run() Object {
	for vm.ip < len(vm.instructions) {
		inst := vm.instructions[vm.ip]
		vm.ip++

		switch inst.Op {
		case OpConstant:
			vm.push(vm.toJashObject(vm.constants[inst.Arg]))

		case OpLoad:
			obj := vm.globals[inst.Arg]
			if obj == nil {
				name := vm.varNames[inst.Arg]
				if name != "" {
					obj, _ = vm.env.Get(name)
				}
			}
			vm.push(obj)

		case OpStore:
			val := vm.pop()
			name := vm.varNames[inst.Arg]
			if name != "" {
				vm.env.Set(name, val)
			}
			vm.globals[inst.Arg] = val

		case OpAdd:
			r := vm.pop()
			l := vm.pop()
			vm.push(evalArith("+", l, r))

		case OpSub:
			r := vm.pop()
			l := vm.pop()
			vm.push(evalArith("-", l, r))

		case OpMul:
			r := vm.pop()
			l := vm.pop()
			vm.push(evalArith("*", l, r))

		case OpDiv:
			r := vm.pop()
			l := vm.pop()
			vm.push(evalArith("/", l, r))

		case OpEq:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(eq(l, r)))

		case OpNeq:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(!eq(l, r)))

		case OpLt:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(cmp(l, r) < 0))

		case OpGt:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(cmp(l, r) > 0))

		case OpLte:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(cmp(l, r) <= 0))

		case OpGte:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(cmp(l, r) >= 0))

		case OpAnd:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(isTruthy(l) && isTruthy(r)))

		case OpOr:
			r := vm.pop()
			l := vm.pop()
			vm.push(nativeBoolToBooleanObject(isTruthy(l) || isTruthy(r)))

		case OpNot:
			val := vm.pop()
			vm.push(nativeBoolToBooleanObject(!isTruthy(val)))

		case OpMinus:
			val := vm.pop()
			switch v := val.(type) {
			case *Integer:
				vm.push(&Integer{Value: -v.Value})
			case *Float:
				vm.push(&Float{Value: -v.Value})
			default:
				vm.push(&Error{Message: "JIT: unsupported -"})
			}

		case OpCall:
			fn := vm.pop()
			args := make([]Object, inst.Arg)
			for i := inst.Arg - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			result := applyFunction(fn, args)
			if isError(result) {
				return result
			}
			vm.push(result)

		case OpReturn:
			return nil

		case OpReturnVal:
			return vm.pop()

		case OpJump:
			vm.ip = inst.Arg

		case OpJumpIfFalse:
			val := vm.pop()
			if !isTruthy(val) {
				vm.ip = inst.Arg
			}

		case OpPop:
			vm.pop()

		case OpNewArray:
			vm.push(&JSONArray{Elements: []Object{}})

		case OpNewObject:
			vm.push(&JSONObject{Pairs: make(map[string]Object)})

		case OpSetMember:
			val := vm.pop()
			obj := vm.pop()
			switch o := obj.(type) {
			case *JSONArray:
				o.Elements = append(o.Elements, val)
				vm.push(o)
			case *JSONObject:
				key := vm.constants[inst.Arg].(string)
				o.Pairs[key] = val
				vm.push(o)
			}

		case OpGetMember:
			obj := vm.pop()
			key := vm.constants[inst.Arg].(string)
			switch o := obj.(type) {
			case *JSONObject:
				if v, ok := o.Pairs[key]; ok {
					vm.push(v)
				} else {
					return &Error{Message: fmt.Sprintf("JIT: key not found: %s", key)}
				}
			case *JSONArray:
				if key == "length" {
					vm.push(&Integer{Value: int64(len(o.Elements))})
				} else {
					return &Error{Message: fmt.Sprintf("JIT: array has no member: %s", key)}
				}
			default:
				return &Error{Message: fmt.Sprintf("JIT: cannot access member on %T", o)}
			}

		case OpNull:
			vm.push(NULL)

		case OpTrue:
			vm.push(TRUE)

		case OpFalse:
			vm.push(FALSE)
		}
	}
	return nil
}

func (vm *vm) push(obj Object) {
	vm.stack[vm.sp] = obj
	vm.sp++
}

func (vm *vm) pop() Object {
	vm.sp--
	return vm.stack[vm.sp]
}

func (vm *vm) toJashObject(v interface{}) Object {
	switch val := v.(type) {
	case string:
		return &String{Value: val}
	case int64:
		return &Integer{Value: val}
	case float64:
		return &Float{Value: val}
	case bool:
		if val {
			return TRUE
		}
		return FALSE
	case nil:
		return NULL
	default:
		return &String{Value: fmt.Sprintf("%v", v)}
	}
}

func evalArith(op string, l, r Object) Object {
	li, lInt := l.(*Integer)
	ri, rInt := r.(*Integer)
	lf, lFloat := l.(*Float)
	rf, rFloat := r.(*Float)

	if lInt && rInt {
		switch op {
		case "+":
			return &Integer{Value: li.Value + ri.Value}
		case "-":
			return &Integer{Value: li.Value - ri.Value}
		case "*":
			return &Integer{Value: li.Value * ri.Value}
		case "/":
			if ri.Value == 0 {
				return &Error{Message: "division by zero"}
			}
			return &Integer{Value: li.Value / ri.Value}
		}
	}

	var lv, rv float64
	if lInt {
		lv = float64(li.Value)
	} else if lFloat {
		lv = lf.Value
	} else if s, ok := l.(*String); ok && op == "+" {
		rs, rok := r.(*String)
		if rok {
			return &String{Value: s.Value + rs.Value}
		}
		return &Error{Message: "JIT: type mismatch"}
	} else {
		return &Error{Message: "JIT: type mismatch"}
	}
	if rInt {
		rv = float64(ri.Value)
	} else if rFloat {
		rv = rf.Value
	} else {
		return &Error{Message: "JIT: type mismatch"}
	}

	switch op {
	case "+":
		return &Float{Value: lv + rv}
	case "-":
		return &Float{Value: lv - rv}
	case "*":
		return &Float{Value: lv * rv}
	case "/":
		if rv == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Float{Value: lv / rv}
	}
	return &Error{Message: "JIT: unknown operator"}
}

func eq(a, b Object) bool {
	switch la := a.(type) {
	case *Integer:
		if lb, ok := b.(*Integer); ok {
			return la.Value == lb.Value
		}
	case *Float:
		if lb, ok := b.(*Float); ok {
			return la.Value == lb.Value
		}
	case *String:
		if lb, ok := b.(*String); ok {
			return la.Value == lb.Value
		}
	case *Boolean:
		if lb, ok := b.(*Boolean); ok {
			return la.Value == lb.Value
		}
	}
	return false
}

func cmp(a, b Object) int {
	switch la := a.(type) {
	case *Integer:
		if lb, ok := b.(*Integer); ok {
			if la.Value < lb.Value {
				return -1
			} else if la.Value > lb.Value {
				return 1
			}
			return 0
		}
	case *Float:
		if lb, ok := b.(*Float); ok {
			if la.Value < lb.Value {
				return -1
			} else if la.Value > lb.Value {
				return 1
			}
			return 0
		}
	}
	return 0
}

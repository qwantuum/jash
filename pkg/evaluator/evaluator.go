package evaluator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/qwantuum/jash/pkg/ast"
)

var StacktraceEnabled bool
var callStack []string

func PushCallStack(name string) {
	if StacktraceEnabled {
		callStack = append(callStack, name)
		fmt.Printf("%s> enter %s\n", indent(len(callStack)-1), name)
	}
}

func PopCallStack() {
	if StacktraceEnabled && len(callStack) > 0 {
		fmt.Printf("%s< exit %s\n", indent(len(callStack)-1), callStack[len(callStack)-1])
		callStack = callStack[:len(callStack)-1]
	}
}

func indent(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("  ", n)
}

func GetStackTrace() string {
	if !StacktraceEnabled || len(callStack) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Stack trace (most recent call last):\n")
	for i := len(callStack) - 1; i >= 0; i-- {
		b.WriteString(fmt.Sprintf("  %s()\n", callStack[i]))
	}
	return b.String()
}

func ClearCallStack() {
	callStack = nil
}

type ObjectType string

const (
	INTEGER_OBJ    ObjectType = "INTEGER"
	FLOAT_OBJ      ObjectType = "FLOAT"
	STRING_OBJ     ObjectType = "STRING"
	BOOLEAN_OBJ    ObjectType = "BOOLEAN"
	NULL_OBJ       ObjectType = "NULL"
	JSON_OBJECT_OBJ ObjectType = "JSON_OBJECT"
	JSON_ARRAY_OBJ  ObjectType = "JSON_ARRAY"
	FUNCTION_OBJ   ObjectType = "FUNCTION"
	BUILTIN_OBJ    ObjectType = "BUILTIN"
	RETURN_OBJ     ObjectType = "RETURN"
	ERROR_OBJ      ObjectType = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type JSONObject struct {
	Pairs map[string]Object
}

func (jo *JSONObject) Type() ObjectType { return JSON_OBJECT_OBJ }
func (jo *JSONObject) Inspect() string {
	var out bytes.Buffer
	out.WriteString("{")
	i := 0
	for k, v := range jo.Pairs {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString("\"")
		out.WriteString(k)
		out.WriteString("\": ")
		out.WriteString(v.Inspect())
		i++
	}
	out.WriteString("}")
	return out.String()
}

type JSONArray struct {
	Elements []Object
}

func (ja *JSONArray) Type() ObjectType { return JSON_ARRAY_OBJ }
func (ja *JSONArray) Inspect() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, e := range ja.Elements {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(e.Inspect())
	}
	out.WriteString("]")
	return out.String()
}

type Function struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	return fmt.Sprintf("def %s(%s) { ... }", f.Name, f.Parameters)
}

type Builtin struct {
	Name string
	Fn   func(args ...Object) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "Error: " + e.Message }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	env := &Environment{
		store: make(map[string]Object),
	}
	env.loadBuiltins()
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := &Environment{
		store: make(map[string]Object),
		outer: outer,
	}
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) loadBuiltins() {
	e.store["print"] = &Builtin{Name: "print", Fn: printFunc}
	e.store["len"] = &Builtin{Name: "len", Fn: lenFunc}
	e.store["serve"] = &Builtin{Name: "serve", Fn: serveFunc}
	e.store["type"] = &Builtin{Name: "type", Fn: typeFunc}

	aiObj := &JSONObject{
		Pairs: map[string]Object{
			"predict": &Builtin{Name: "ai.predict", Fn: aiPredictFunc},
			"ollama":  &Builtin{Name: "ai.ollama", Fn: ollamaFunc},
		},
	}
	e.store["ai"] = aiObj

	jashUIObj := &JSONObject{
		Pairs: map[string]Object{
			"window": &Builtin{Name: "jash_ui.window", Fn: uiWindowFunc},
		},
	}
	e.store["jash_ui"] = jashUIObj

	imageObj := &JSONObject{
		Pairs: map[string]Object{
			"ascii": &Builtin{Name: "image.ascii", Fn: imageASCIIFunc},
		},
	}
	e.store["image"] = imageObj
}

var (
	NULL    = &Null{}
	TRUE    = &Boolean{Value: true}
	FALSE   = &Boolean{Value: false}
)

func Eval(node ast.Node, env *Environment) Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.BlockStatement:
		return evalBlockStatement(n, env)
	case *ast.FunctionStatement:
		return evalFunctionStatement(n, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(n, env)
	case *ast.AssignStatement:
		return evalAssignStatement(n, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IfStatement:
		return evalIfStatement(n, env)
	case *ast.ForStatement:
		return evalForStatement(n, env)
	case *ast.WhileStatement:
		return evalWhileStatement(n, env)
	case *ast.RepeatStatement:
		return evalRepeatStatement(n, env)
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.NumberLiteral:
		return evalNumberLiteral(n)
	case *ast.StringLiteral:
		return &String{Value: n.Value}
	case *ast.BooleanLiteral:
		if n.Value {
			return TRUE
		}
		return FALSE
	case *ast.NullLiteral:
		return NULL
	case *ast.JSONObject:
		return evalJSONObject(n, env)
	case *ast.JSONArray:
		return evalJSONArray(n, env)
	case *ast.CallExpression:
		return evalCallExpression(n, env)
	case *ast.MemberAccess:
		return evalMemberAccess(n, env)
	case *ast.InfixExpression:
		return evalInfixExpression(n, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(n, env)
	}
	return NULL
}

func evalProgram(program *ast.Program, env *Environment) Object {
	if GlobalJIT != nil {
		c := newCompiler()
		instructions, constants, varNames, err := c.Compile(program)
		if err == nil {
			globals := make([]Object, len(instructions)+256)
			vm := newVM(instructions, constants, varNames, globals, env)
			result := vm.Run()
			if result != nil {
				return result
			}
			return NULL
		}
	}

	var result Object
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		switch r := result.(type) {
		case *ReturnValue:
			return r.Value
		case *Error:
			return r
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *Environment) Object {
	var result Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalFunctionStatement(node *ast.FunctionStatement, env *Environment) Object {
	fn := &Function{
		Name:       node.Name.Value,
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}

	if GlobalJIT != nil {
		GlobalJIT.Compile(node.Name.Value, node.Parameters, node.Body)
	}

	return env.Set(node.Name.Value, fn)
}

func evalReturnStatement(node *ast.ReturnStatement, env *Environment) Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	return &ReturnValue{Value: val}
}

func evalAssignStatement(node *ast.AssignStatement, env *Environment) Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	return env.Set(node.Name.Value, val)
}

func evalIfStatement(node *ast.IfStatement, env *Environment) Object {
	condition := Eval(node.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(node.Body, env)
	} else if node.ElseBody != nil {
		return Eval(node.ElseBody, env)
	}
	return NULL
}

func evalForStatement(node *ast.ForStatement, env *Environment) Object {
	iterable := Eval(node.Iterable, env)
	if isError(iterable) {
		return iterable
	}

	var result Object = NULL

	switch it := iterable.(type) {
	case *JSONArray:
		for _, elem := range it.Elements {
			env.Set(node.Variable.Value, elem)
			r := Eval(node.Body, env)
			if r != nil {
				if r.Type() == RETURN_OBJ {
					return r
				}
				if isError(r) {
					return r
				}
				result = r
			}
		}
	case *String:
		for _, ch := range it.Value {
			env.Set(node.Variable.Value, &String{Value: string(ch)})
			r := Eval(node.Body, env)
			if r != nil {
				if r.Type() == RETURN_OBJ {
					return r
				}
				if isError(r) {
					return r
				}
				result = r
			}
		}
	case *JSONObject:
		for key := range it.Pairs {
			env.Set(node.Variable.Value, &String{Value: key})
			r := Eval(node.Body, env)
			if r != nil {
				if r.Type() == RETURN_OBJ {
					return r
				}
				if isError(r) {
					return r
				}
				result = r
			}
		}
	default:
		return &Error{Message: fmt.Sprintf("cannot iterate over %s", iterable.Type())}
	}

	return result
}

func evalWhileStatement(node *ast.WhileStatement, env *Environment) Object {
	var result Object = NULL
	for {
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		r := Eval(node.Body, env)
		if r != nil {
			if r.Type() == RETURN_OBJ {
				return r
			}
			if isError(r) {
				return r
			}
			result = r
		}
	}
	return result
}

func evalRepeatStatement(node *ast.RepeatStatement, env *Environment) Object {
	countObj := Eval(node.Count, env)
	if isError(countObj) {
		return countObj
	}

	count, ok := countObj.(*Integer)
	if !ok {
		return &Error{Message: fmt.Sprintf("repeat() requires an integer argument, got %s", countObj.Type())}
	}

	var result Object = NULL
	for i := int64(0); i < count.Value; i++ {
		r := Eval(node.Body, env)
		if r != nil {
			if r.Type() == RETURN_OBJ {
				return r
			}
			if isError(r) {
				return r
			}
			result = r
		}
	}
	return result
}

func evalIdentifier(node *ast.Identifier, env *Environment) Object {
	if node.Value == "true" {
		return TRUE
	}
	if node.Value == "false" {
		return FALSE
	}
	if node.Value == "null" {
		return NULL
	}

	val, ok := env.Get(node.Value)
	if !ok {
		return &Error{Message: fmt.Sprintf("undefined variable: %s", node.Value)}
	}
	return val
}

func evalNumberLiteral(node *ast.NumberLiteral) Object {
	val := node.Value
	if strings.Contains(val, ".") {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return &Error{Message: fmt.Sprintf("invalid float: %s", val)}
		}
		return &Float{Value: f}
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return &Error{Message: fmt.Sprintf("invalid integer: %s", val)}
	}
	return &Integer{Value: i}
}

func evalJSONObject(node *ast.JSONObject, env *Environment) Object {
	pairs := make(map[string]Object)
	for k, v := range node.Pairs {
		val := Eval(v, env)
		if isError(val) {
			return val
		}
		pairs[k] = val
	}
	return &JSONObject{Pairs: pairs}
}

func evalJSONArray(node *ast.JSONArray, env *Environment) Object {
	elements := make([]Object, len(node.Elements))
	for i, elem := range node.Elements {
		val := Eval(elem, env)
		if isError(val) {
			return val
		}
		elements[i] = val
	}
	return &JSONArray{Elements: elements}
}

func evalCallExpression(node *ast.CallExpression, env *Environment) Object {
	function := Eval(node.Function, env)
	if isError(function) {
		return function
	}

	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(function, args)
}

func evalMemberAccess(node *ast.MemberAccess, env *Environment) Object {
	obj := Eval(node.Object, env)
	if isError(obj) {
		return obj
	}

	switch o := obj.(type) {
	case *JSONObject:
		val, ok := o.Pairs[node.Member.Value]
		if !ok {
			return &Error{Message: fmt.Sprintf("key not found: %s", node.Member.Value)}
		}
		return val
	case *Builtin:
		return obj
	case *JSONArray:
		if node.Member.Value == "length" {
			return &Integer{Value: int64(len(o.Elements))}
		}
		return &Error{Message: fmt.Sprintf("array has no member: %s", node.Member.Value)}
	default:
		return &Error{Message: fmt.Sprintf("cannot access member on %s", obj.Type())}
	}
}

func evalInfixExpression(node *ast.InfixExpression, env *Environment) Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(node.Operator, left, right)
	case left.Type() == FLOAT_OBJ && right.Type() == FLOAT_OBJ:
		return evalFloatInfixExpression(node.Operator, left, right)
	case left.Type() == INTEGER_OBJ && right.Type() == FLOAT_OBJ:
		l := &Float{Value: float64(left.(*Integer).Value)}
		return evalFloatInfixExpression(node.Operator, l, right)
	case left.Type() == FLOAT_OBJ && right.Type() == INTEGER_OBJ:
		r := &Float{Value: float64(right.(*Integer).Value)}
		return evalFloatInfixExpression(node.Operator, left, r)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(node.Operator, left, right)
	case node.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return &Error{Message: fmt.Sprintf("type mismatch: %s %s %s",
			left.Type(), node.Operator, right.Type())}
	}
}

func evalPrefixExpression(node *ast.PrefixExpression, env *Environment) Object {
	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch node.Operator {
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "not":
		return evalNotPrefixOperatorExpression(right)
	default:
		return &Error{Message: fmt.Sprintf("unknown prefix operator: %s", node.Operator)}
	}
}

func evalIntegerInfixExpression(operator string, left, right Object) Object {
	l := left.(*Integer).Value
	r := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: l + r}
	case "-":
		return &Integer{Value: l - r}
	case "*":
		return &Integer{Value: l * r}
	case "/":
		if r == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Integer{Value: l / r}
	case "<":
		return nativeBoolToBooleanObject(l < r)
	case ">":
		return nativeBoolToBooleanObject(l > r)
	case "<=":
		return nativeBoolToBooleanObject(l <= r)
	case ">=":
		return nativeBoolToBooleanObject(l >= r)
	case "==":
		return nativeBoolToBooleanObject(l == r)
	case "!=":
		return nativeBoolToBooleanObject(l != r)
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s", operator)}
	}
}

func evalFloatInfixExpression(operator string, left, right Object) Object {
	l := left.(*Float).Value
	r := right.(*Float).Value

	switch operator {
	case "+":
		return &Float{Value: l + r}
	case "-":
		return &Float{Value: l - r}
	case "*":
		return &Float{Value: l * r}
	case "/":
		if r == 0 {
			return &Error{Message: "division by zero"}
		}
		return &Float{Value: l / r}
	case "<":
		return nativeBoolToBooleanObject(l < r)
	case ">":
		return nativeBoolToBooleanObject(l > r)
	case "<=":
		return nativeBoolToBooleanObject(l <= r)
	case ">=":
		return nativeBoolToBooleanObject(l >= r)
	case "==":
		return nativeBoolToBooleanObject(l == r)
	case "!=":
		return nativeBoolToBooleanObject(l != r)
	default:
		return &Error{Message: fmt.Sprintf("unknown operator: %s", operator)}
	}
}

func evalStringInfixExpression(operator string, left, right Object) Object {
	l := left.(*String).Value
	r := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: l + r}
	case "==":
		return nativeBoolToBooleanObject(l == r)
	case "!=":
		return nativeBoolToBooleanObject(l != r)
	default:
		return &Error{Message: fmt.Sprintf("unknown string operator: %s", operator)}
	}
}

func evalMinusPrefixOperatorExpression(right Object) Object {
	switch r := right.(type) {
	case *Integer:
		return &Integer{Value: -r.Value}
	case *Float:
		return &Float{Value: -r.Value}
	default:
		return &Error{Message: fmt.Sprintf("unknown prefix - for %s", right.Type())}
	}
}

func evalNotPrefixOperatorExpression(right Object) Object {
	return nativeBoolToBooleanObject(!isTruthy(right))
}

func evalExpressions(exps []ast.Expression, env *Environment) []Object {
	result := make([]Object, 0, len(exps))
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn Object, args []Object) Object {
	switch f := fn.(type) {
	case *Function:
		PushCallStack(f.Name)

		if GlobalJIT != nil && GlobalJIT.IsCompiled(f.Name) {
			env := NewEnclosedEnvironment(f.Env)
			for i, param := range f.Parameters {
				if i < len(args) {
					env.Set(param.Value, args[i])
				}
			}
			result := GlobalJIT.Execute(f.Name, args, env)
			PopCallStack()
			return result
		}

		env := NewEnclosedEnvironment(f.Env)
		for i, param := range f.Parameters {
			if i < len(args) {
				env.Set(param.Value, args[i])
			}
		}
		result := Eval(f.Body, env)
		if isError(result) {
			PopCallStack()
			return result
		}
		PopCallStack()
		if rv, ok := result.(*ReturnValue); ok {
			return rv.Value
		}
		return result
	case *Builtin:
		PushCallStack(f.Name)
		result := f.Fn(args...)
		PopCallStack()
		return result
	default:
		return &Error{Message: fmt.Sprintf("not a function: %s", fn.Type())}
	}
}

func ApplyFunction(fn Object, args []Object) Object {
	return applyFunction(fn, args)
}

func nativeBoolToBooleanObject(b bool) *Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj Object) bool {
	switch o := obj.(type) {
	case *Boolean:
		return o.Value
	case *Null:
		return false
	case *Integer:
		return o.Value != 0
	case *Float:
		return o.Value != 0
	case *String:
		return o.Value != ""
	default:
		return true
	}
}

func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func printFunc(args ...Object) Object {
	parts := make([]string, len(args))
	for i, arg := range args {
		parts[i] = arg.Inspect()
	}
	fmt.Println(strings.Join(parts, " "))
	return NULL
}

func lenFunc(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: "len() requires exactly 1 argument"}
	}
	switch o := args[0].(type) {
	case *String:
		return &Integer{Value: int64(len(o.Value))}
	case *JSONArray:
		return &Integer{Value: int64(len(o.Elements))}
	case *JSONObject:
		return &Integer{Value: int64(len(o.Pairs))}
	default:
		return &Error{Message: fmt.Sprintf("len() not supported for %s", args[0].Type())}
	}
}

func typeFunc(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: "type() requires exactly 1 argument"}
	}
	return &String{Value: string(args[0].Type())}
}

func aiPredictFunc(args ...Object) Object {
	input := ""
	if len(args) > 0 {
		input = args[0].Inspect()
	}

	return &JSONObject{
		Pairs: map[string]Object{
			"prediction": &String{Value: "positive"},
			"confidence": &Float{Value: 0.9532},
			"model":      &String{Value: "Jash-AI v1.0"},
			"input":      &String{Value: input},
		},
	}
}

func ollamaFunc(args ...Object) Object {
	if len(args) != 1 {
		return &Error{Message: "ollama() requires exactly 1 argument: the Ollama server URL"}
	}

	baseURL, ok := args[0].(*String)
	if !ok {
		return &Error{Message: "ollama() argument must be a string URL"}
	}

	url := strings.TrimRight(baseURL.Value, "/")

	return &JSONObject{
		Pairs: map[string]Object{
			"generate": &Builtin{Fn: makeOllamaGenerate(url)},
			"chat":     &Builtin{Fn: makeOllamaChat(url)},
			"list":     &Builtin{Fn: makeOllamaList(url)},
		},
	}
}

func makeOllamaGenerate(baseURL string) func(args ...Object) Object {
	return func(args ...Object) Object {
		if len(args) < 2 {
			return &Error{Message: "ollama.generate() requires at least 2 arguments: model and prompt"}
		}

		model, ok := args[0].(*String)
		if !ok {
			return &Error{Message: "first argument to ollama.generate() must be a string (model name)"}
		}

		prompt, ok := args[1].(*String)
		if !ok {
			return &Error{Message: "second argument to ollama.generate() must be a string (prompt)"}
		}

		body := map[string]interface{}{
			"model":  model.Value,
			"prompt": prompt.Value,
			"stream": false,
		}

		return ollamaRequest(baseURL+"/api/generate", body)
	}
}

func makeOllamaChat(baseURL string) func(args ...Object) Object {
	return func(args ...Object) Object {
		if len(args) < 2 {
			return &Error{Message: "ollama.chat() requires at least 2 arguments: model and messages array"}
		}

		model, ok := args[0].(*String)
		if !ok {
			return &Error{Message: "first argument to ollama.chat() must be a string (model name)"}
		}

		messages, ok := args[1].(*JSONArray)
		if !ok {
			return &Error{Message: "second argument to ollama.chat() must be an array of message objects"}
		}

		body := map[string]interface{}{
			"model":    model.Value,
			"messages": jashToGoObject(messages),
			"stream":   false,
		}

		return ollamaRequest(baseURL+"/api/chat", body)
	}
}

func makeOllamaList(baseURL string) func(args ...Object) Object {
	return func(args ...Object) Object {
		resp, err := http.Get(baseURL + "/api/tags")
		if err != nil {
			return &Error{Message: fmt.Sprintf("ollama list request failed: %s", err)}
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return &Error{Message: fmt.Sprintf("failed to read ollama response: %s", err)}
		}

		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return &Error{Message: fmt.Sprintf("failed to parse ollama response: %s", err)}
		}

		return goToJashObject(result)
	}
}

func ollamaRequest(endpoint string, body map[string]interface{}) Object {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return &Error{Message: fmt.Sprintf("failed to marshal request: %s", err)}
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return &Error{Message: fmt.Sprintf("ollama request failed: %s", err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Error{Message: fmt.Sprintf("failed to read ollama response: %s", err)}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return &Error{Message: fmt.Sprintf("failed to parse ollama response: %s", err)}
	}

	return goToJashObject(result)
}

func goToJashObject(v interface{}) Object {
	switch val := v.(type) {
	case string:
		return &String{Value: val}
	case float64:
		if val == float64(int64(val)) {
			return &Integer{Value: int64(val)}
		}
		return &Float{Value: val}
	case bool:
		return nativeBoolToBooleanObject(val)
	case nil:
		return NULL
	case map[string]interface{}:
		pairs := make(map[string]Object, len(val))
		for k, vv := range val {
			pairs[k] = goToJashObject(vv)
		}
		return &JSONObject{Pairs: pairs}
	case []interface{}:
		elems := make([]Object, len(val))
		for i, vv := range val {
			elems[i] = goToJashObject(vv)
		}
		return &JSONArray{Elements: elems}
	default:
		return &String{Value: fmt.Sprintf("%v", v)}
	}
}

func jashToGoObject(obj Object) interface{} {
	switch o := obj.(type) {
	case *Integer:
		return o.Value
	case *Float:
		return o.Value
	case *String:
		return o.Value
	case *Boolean:
		return o.Value
	case *Null:
		return nil
	case *JSONObject:
		m := make(map[string]interface{}, len(o.Pairs))
		for k, v := range o.Pairs {
			m[k] = jashToGoObject(v)
		}
		return m
	case *JSONArray:
		arr := make([]interface{}, len(o.Elements))
		for i, v := range o.Elements {
			arr[i] = jashToGoObject(v)
		}
		return arr
	default:
		return o.Inspect()
	}
}

func serveFunc(args ...Object) Object {
	if len(args) != 2 {
		return &Error{Message: "serve() requires exactly 2 arguments: port and handler"}
	}

	port, ok := args[0].(*Integer)
	if !ok {
		return &Error{Message: "first argument to serve() must be an integer port number"}
	}

	handler, ok := args[1].(*Function)
	if !ok {
		return &Error{Message: "second argument to serve() must be a function"}
	}

	addr := fmt.Sprintf(":%d", port.Value)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		bodyBytes := make([]byte, 0)
		if r.Body != nil {
			var readErr error
			bodyBytes, readErr = readBody(r)
			if readErr != nil {
				http.Error(w, readErr.Error(), http.StatusInternalServerError)
				return
			}
		}

		reqObj := &JSONObject{
			Pairs: map[string]Object{
				"method": &String{Value: r.Method},
				"path":   &String{Value: r.URL.Path},
				"body":   &String{Value: string(bodyBytes)},
				"query":  &String{Value: r.URL.RawQuery},
			},
		}

		result := applyFunction(handler, []Object{reqObj})
		if isError(result) {
			http.Error(w, result.Inspect(), http.StatusInternalServerError)
			return
		}

		jsonBytes := objectToJSON(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	})

	fmt.Printf("Jash server listening on http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		return &Error{Message: fmt.Sprintf("server error: %s", err)}
	}

	return NULL
}

func readBody(r *http.Request) ([]byte, error) {
	defer r.Body.Close()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func objectToJSON(obj Object) []byte {
	var buf bytes.Buffer
	writeJSON(&buf, obj)
	return buf.Bytes()
}

func writeJSON(buf *bytes.Buffer, obj Object) {
	switch o := obj.(type) {
	case *Integer:
		buf.WriteString(fmt.Sprintf("%d", o.Value))
	case *Float:
		if math.IsInf(o.Value, 0) || math.IsNaN(o.Value) {
			buf.WriteString("null")
			return
		}
		buf.WriteString(fmt.Sprintf("%g", o.Value))
	case *String:
		d, _ := json.Marshal(o.Value)
		buf.Write(d)
	case *Boolean:
		buf.WriteString(fmt.Sprintf("%t", o.Value))
	case *Null:
		buf.WriteString("null")
	case *JSONObject:
		buf.WriteByte('{')
		i := 0
		for k, v := range o.Pairs {
			if i > 0 {
				buf.WriteByte(',')
			}
			kd, _ := json.Marshal(k)
			buf.Write(kd)
			buf.WriteByte(':')
			writeJSON(buf, v)
			i++
		}
		buf.WriteByte('}')
	case *JSONArray:
		buf.WriteByte('[')
		for i, v := range o.Elements {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeJSON(buf, v)
		}
		buf.WriteByte(']')
	default:
		buf.WriteString(fmt.Sprintf("\"%s\"", o.Inspect()))
	}
}

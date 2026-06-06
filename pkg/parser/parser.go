package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/qwantuum/jash/pkg/ast"
	"github.com/qwantuum/jash/pkg/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	ASSIGN
	OR
	AND
	EQUALS
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	CALL
	MEMBER
)

type Parser struct {
	tokens    []token.Token
	position  int
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(tokens []token.Token) *Parser {
	p := &Parser{
		tokens:   tokens,
		position: 0,
		errors:   []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.NUMBER, p.parseNumberLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.NULL, p.parseNullLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACE, p.parseJSONObject)
	p.registerPrefix(token.LBRACKET, p.parseJSONArray)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.STAR, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.DOT, p.parseMemberAccess)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.position++
	if p.position >= len(p.tokens) {
		p.peekToken = token.Token{Type: token.EOF}
	} else {
		p.peekToken = p.tokens[p.position]
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expect(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors,
		fmt.Sprintf("line %d: expected %s, got %s (%s)",
			p.curToken.Line, t, p.peekToken.Type, p.peekToken.Literal))
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(t token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t token.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		if p.curToken.Type == token.NEWLINE || p.curToken.Type == token.DEDENT {
			p.nextToken()
			continue
		}

		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.DEF:
		return p.parseFunctionStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.IDENT:
		if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{}
	p.nextToken()

	stmt.Name = &ast.Identifier{Value: p.curToken.Literal}

	if !p.expect(token.LPAREN) {
		return nil
	}
	p.nextToken()

	stmt.Parameters = []*ast.Identifier{}
	for p.curToken.Type != token.RPAREN && p.curToken.Type != token.EOF {
		if p.curToken.Type == token.COMMA {
			p.nextToken()
			continue
		}
		param := &ast.Identifier{Value: p.curToken.Literal}
		stmt.Parameters = append(stmt.Parameters, param)
		p.nextToken()
	}

	if !p.expect(token.COLON) {
		return nil
	}

	stmt.Body = p.parseBlockBody()
	p.nextToken()
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.NEWLINE || p.peekToken.Type == token.DEDENT {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{
		Name: &ast.Identifier{Value: p.curToken.Literal},
	}
	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.NEWLINE || p.peekToken.Type == token.DEDENT {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{}
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expect(token.COLON) {
		return nil
	}

	stmt.Body = p.parseBlockBody()
	p.nextToken()

	if p.curTokenIs(token.ELIF) {
		elseStmt := p.parseIfStatement()
		stmt.ElseBody = &ast.BlockStatement{
			Statements: []ast.Statement{
				elseStmt,
			},
		}
	} else if p.curTokenIs(token.ELSE) {
		p.nextToken()
		stmt.ElseBody = p.parseBlockBody()
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{}
	p.nextToken()

	stmt.Variable = &ast.Identifier{Value: p.curToken.Literal}
	if !p.expect(token.IN) {
		return nil
	}
	p.nextToken()

	stmt.Iterable = p.parseExpression(LOWEST)

	if !p.expect(token.COLON) {
		return nil
	}

	stmt.Body = p.parseBlockBody()
	p.nextToken()
	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{}
	p.nextToken()

	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expect(token.COLON) {
		return nil
	}

	stmt.Body = p.parseBlockBody()
	p.nextToken()
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.NEWLINE || p.peekToken.Type == token.DEDENT {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockBody() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Statements: []ast.Statement{},
	}

	if !p.expect(token.NEWLINE) {
		return block
	}
	if !p.expect(token.INDENT) {
		return block
	}
	p.nextToken()

	for {
		for p.curTokenIs(token.NEWLINE) {
			p.nextToken()
		}

		if p.curTokenIs(token.DEDENT) || p.curTokenIs(token.EOF) {
			break
		}

		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors,
			fmt.Sprintf("line %d: unexpected token %s (%s)",
				p.curToken.Line, p.curToken.Type, p.curToken.Literal))
		return nil
	}
	leftExp := prefix()

	for !p.curTokenIs(token.NEWLINE) && !p.curTokenIs(token.EOF) &&
		precedence < p.peekPrecedence() {

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Value: p.curToken.Literal}
}

func (p *Parser) parseNumberLiteral() ast.Expression {
	return &ast.NumberLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	val, _ := strconv.ParseBool(p.curToken.Literal)
	return &ast.BooleanLiteral{Value: val}
}

func (p *Parser) parseNullLiteral() ast.Expression {
	return &ast.NullLiteral{}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expect(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseJSONObject() ast.Expression {
	obj := &ast.JSONObject{
		Pairs: make(map[string]ast.Expression),
	}

	p.nextToken()

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
		return obj
	}

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.COMMA) || p.curTokenIs(token.NEWLINE) {
			p.nextToken()
			continue
		}

		if p.curToken.Type != token.STRING {
			p.errors = append(p.errors,
				fmt.Sprintf("line %d: expected string key in JSON object, got %s",
					p.curToken.Line, p.curToken.Type))
			return obj
		}
		key := p.curToken.Literal
		p.nextToken()

		if !p.curTokenIs(token.COLON) {
			p.errors = append(p.errors,
				fmt.Sprintf("line %d: expected ':' after JSON object key", p.curToken.Line))
			return obj
		}
		p.nextToken()

		value := p.parseExpression(LOWEST)
		obj.Pairs[key] = value

		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RBRACE) {
		p.nextToken()
	}

	return obj
}

func (p *Parser) parseJSONArray() ast.Expression {
	arr := &ast.JSONArray{
		Elements: []ast.Expression{},
	}

	p.nextToken()

	if p.curTokenIs(token.RBRACKET) {
		p.nextToken()
		return arr
	}

	for !p.curTokenIs(token.RBRACKET) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.COMMA) || p.curTokenIs(token.NEWLINE) {
			p.nextToken()
			continue
		}

		elem := p.parseExpression(LOWEST)
		arr.Elements = append(arr.Elements, elem)

		p.nextToken()
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RBRACKET) {
		p.nextToken()
	}

	return arr
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Left:     left,
		Operator: p.curToken.Literal,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Function:  function,
		Arguments: []ast.Expression{},
	}

	p.nextToken()

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return expr
	}

	for !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.COMMA) {
			p.nextToken()
			continue
		}

		arg := p.parseExpression(LOWEST)
		expr.Arguments = append(expr.Arguments, arg)
		p.nextToken()

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
	}

	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	return expr
}

func (p *Parser) parseMemberAccess(obj ast.Expression) ast.Expression {
	expr := &ast.MemberAccess{
		Object: obj,
	}
	p.nextToken()
	expr.Member = &ast.Identifier{Value: p.curToken.Literal}
	return expr
}

func (p *Parser) peekPrecedence() int {
	return precedenceOf(p.peekToken.Type)
}

func (p *Parser) curPrecedence() int {
	return precedenceOf(p.curToken.Type)
}

func precedenceOf(t token.TokenType) int {
	switch t {
	case token.OR:
		return OR
	case token.AND:
		return AND
	case token.EQ, token.NEQ:
		return EQUALS
	case token.LT, token.GT, token.LTE, token.GTE:
		return COMPARISON
	case token.PLUS, token.MINUS:
		return SUM
	case token.STAR, token.SLASH:
		return PRODUCT
	case token.LPAREN:
		return CALL
	case token.DOT:
		return MEMBER
	default:
		return LOWEST
	}
}



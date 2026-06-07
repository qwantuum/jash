package lexer

import (
	"fmt"
	"strings"

	"github.com/qwantuum/jash/pkg/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int

	indentStack []int
	tokens      []token.Token
	errors      []string
}

func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		indentStack: []int{0},
	}
	return l
}

func (l *Lexer) Tokenize() ([]token.Token, []string) {
	l.readChar()

	beginningOfLine := true

	for l.ch != 0 {
		if beginningOfLine {
			indent := l.countIndent()

			if l.ch == '\n' || l.ch == '\r' {
				l.readChar()
				beginningOfLine = true
				continue
			}

			if l.ch == '#' {
				l.skipComment()
				beginningOfLine = true
				continue
			}

			l.handleIndent(indent)
			beginningOfLine = false
		}

		switch {
		case l.ch == ' ' || l.ch == '\t':
			l.readChar()
		case l.ch == '#':
			l.skipComment()
			beginningOfLine = true
		case l.ch == '\n' || l.ch == '\r':
			l.emitToken(token.NEWLINE, "\\n")
			if l.ch == '\r' {
				l.readChar()
			}
			l.readChar()
			beginningOfLine = true
		case isLetter(l.ch):
			l.readIdentifier()
		case isDigit(l.ch):
			l.readNumber()
		case l.ch == '"':
			l.readString()
		case l.ch == '\'':
			l.readSingleQuotedString()
		case l.ch == '=':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				l.emitToken(token.EQ, string(ch)+string(l.ch))
				l.readChar()
			} else {
				l.emitToken(token.ASSIGN, string(l.ch))
				l.readChar()
			}
		case l.ch == '+':
			l.emitToken(token.PLUS, string(l.ch))
			l.readChar()
		case l.ch == '-':
			l.emitToken(token.MINUS, string(l.ch))
			l.readChar()
		case l.ch == '*':
			l.emitToken(token.STAR, string(l.ch))
			l.readChar()
		case l.ch == '/':
			l.emitToken(token.SLASH, string(l.ch))
			l.readChar()
		case l.ch == '(':
			l.emitToken(token.LPAREN, string(l.ch))
			l.readChar()
		case l.ch == ')':
			l.emitToken(token.RPAREN, string(l.ch))
			l.readChar()
		case l.ch == '{':
			l.emitToken(token.LBRACE, string(l.ch))
			l.readChar()
		case l.ch == '}':
			l.emitToken(token.RBRACE, string(l.ch))
			l.readChar()
		case l.ch == '[':
			l.emitToken(token.LBRACKET, string(l.ch))
			l.readChar()
		case l.ch == ']':
			l.emitToken(token.RBRACKET, string(l.ch))
			l.readChar()
		case l.ch == ':':
			l.emitToken(token.COLON, string(l.ch))
			l.readChar()
		case l.ch == ',':
			l.emitToken(token.COMMA, string(l.ch))
			l.readChar()
		case l.ch == '.':
			l.emitToken(token.DOT, string(l.ch))
			l.readChar()
		case l.ch == '<':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				l.emitToken(token.LTE, string(ch)+string(l.ch))
				l.readChar()
			} else {
				l.emitToken(token.LT, string(l.ch))
				l.readChar()
			}
		case l.ch == '>':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				l.emitToken(token.GTE, string(ch)+string(l.ch))
				l.readChar()
			} else {
				l.emitToken(token.GT, string(l.ch))
				l.readChar()
			}
		case l.ch == '!':
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				l.emitToken(token.NEQ, string(ch)+string(l.ch))
				l.readChar()
			} else {
				l.errors = append(l.errors,
					fmt.Sprintf("line %d: expected '!=' got '!'", l.line))
				l.readChar()
			}
		default:
			l.errors = append(l.errors,
				fmt.Sprintf("line %d: unexpected character '%c'", l.line, l.ch))
			l.readChar()
		}
	}

	for len(l.indentStack) > 1 {
		l.emitDedent()
	}
	l.emitToken(token.EOF, "")

	return l.tokens, l.errors
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) countIndent() int {
	indent := 0
	for l.ch == ' ' || l.ch == '\t' {
		if l.ch == ' ' {
			indent++
		} else {
			indent += 4
		}
		l.readChar()
	}
	return indent
}

func (l *Lexer) handleIndent(indent int) {
	top := l.indentStack[len(l.indentStack)-1]

	if indent > top {
		l.indentStack = append(l.indentStack, indent)
		l.emitToken(token.INDENT, fmt.Sprintf("%d", indent))
	} else if indent < top {
		for len(l.indentStack) > 1 && indent < l.indentStack[len(l.indentStack)-1] {
			l.emitDedent()
		}
		if indent != l.indentStack[len(l.indentStack)-1] {
			l.errors = append(l.errors,
				fmt.Sprintf("line %d: inconsistent indentation (expected %d spaces, got %d)",
					l.line, l.indentStack[len(l.indentStack)-1], indent))
		}
	}
}

func (l *Lexer) emitToken(t token.TokenType, literal string) {
	l.tokens = append(l.tokens, token.Token{
		Type:    t,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	})
}

func (l *Lexer) emitDedent() {
	l.tokens = append(l.tokens, token.Token{
		Type:    token.DEDENT,
		Literal: "",
		Line:    l.line,
		Column:  l.column,
	})
	l.indentStack = l.indentStack[:len(l.indentStack)-1]
}

func (l *Lexer) skipComment() {
	for l.ch != 0 && l.ch != '\n' && l.ch != '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() {
	start := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	ident := l.input[start:l.position]
	tokType := token.LookupIdent(ident)
	l.emitToken(tokType, ident)
}

func (l *Lexer) readNumber() {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	literal := l.input[start:l.position]
	l.emitToken(token.NUMBER, literal)
}

func (l *Lexer) readString() {
	l.readChar()
	var buf strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '"':
				buf.WriteByte('"')
			case '\\':
				buf.WriteByte('\\')
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case 'r':
				buf.WriteByte('\r')
			default:
				buf.WriteByte('\\')
				buf.WriteByte(l.ch)
			}
		} else {
			buf.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar()
	}
	l.emitToken(token.STRING, buf.String())
}

func (l *Lexer) readSingleQuotedString() {
	l.readChar()
	var buf strings.Builder
	for l.ch != '\'' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case '\'':
				buf.WriteByte('\'')
			case '\\':
				buf.WriteByte('\\')
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			default:
				buf.WriteByte('\\')
				buf.WriteByte(l.ch)
			}
		} else {
			buf.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '\'' {
		l.readChar()
	}
	l.emitToken(token.STRING, buf.String())
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

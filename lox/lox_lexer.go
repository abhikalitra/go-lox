package lox

import "strconv"

type TokenType int

const (
	LeftParen TokenType = iota
	RightParen
	LeftBrace
	RightBrace
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BangEqual
	EQUAL
	EqualEqual
	GREATER
	GreaterEqual
	LESS
	LessEqual

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
	ERROR
)

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	keywords["and"] = AND
	keywords["class"] = CLASS
	keywords["else"] = ELSE
	keywords["false"] = FALSE
	keywords["for"] = FOR
	keywords["fun"] = FUN
	keywords["if"] = IF
	keywords["nil"] = NIL
	keywords["or"] = OR
	keywords["print"] = PRINT
	keywords["return"] = RETURN
	keywords["super"] = SUPER
	keywords["this"] = THIS
	keywords["true"] = TRUE
	keywords["var"] = VAR
	keywords["while"] = WHILE
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   interface{}
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func NewErrorToken() Token {
	return Token{TokenType: ERROR}
}

type Scanner struct {
	Source string
	Tokens []Token

	start   int
	current int
	line    int
}

func NewScanner() Scanner {
	return Scanner{current: 0, start: 0, line: 1}
}

func (s *Scanner) ScanTokens() {
	for {
		s.start = s.current
		s.scanToken()
		if s.AtEnd() {
			break
		}
	}
	s.addToken(EOF)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen)
	case ')':
		s.addToken(RightParen)
	case '{':
		s.addToken(LeftBrace)
	case '}':
		s.addToken(RightBrace)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		s.addTokenWithDual(s.match('='), BangEqual, BANG)
	case '=':
		s.addTokenWithDual(s.match('='), EqualEqual, EQUAL)
	case '<':
		s.addTokenWithDual(s.match('='), LessEqual, LESS)
	case '>':
		s.addTokenWithDual(s.match('='), GreaterEqual, GREATER)
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for {
				if s.peek() == '\n' || s.AtEnd() {
					break
				} else {
					s.advance()
				}
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
		break
	case '"':
		s.string()
	default:
		if s.IsDigit(c) {
			s.number()
		} else if s.IsAlpha(c) {
			s.identifier()
		} else {
			s.error(s.line, "Unexpected character.")
		}

	}
}

func (s *Scanner) advance() uint8 {
	r := s.Source[s.current]
	s.current++
	return r
}

func (s *Scanner) addToken(tt TokenType) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, NewToken(tt, text, nil, s.line))
}

func (s *Scanner) addTokenWithDual(match bool, first TokenType, second TokenType) {
	if match {
		s.addToken(first)
	} else {
		s.addToken(second)
	}
}

func (s *Scanner) addTokenWithLiteral(tt TokenType, literal interface{}) {
	text := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, NewToken(tt, text, literal, s.line))
}

func (s *Scanner) match(expected int32) bool {
	if s.AtEnd() {
		return false
	}
	r := s.Source[s.current]
	if r != uint8(expected) {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) AtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) error(line int, msg string) {

}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.AtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.AtEnd() {
		s.error(s.line, "Unterminated string..")
	}
	s.advance()
	value := s.Source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) number() {
	for s.IsDigit(s.peek()) {
		s.advance()
	}
	println(s.peek(), s.peekNext(), s.IsDigit(s.peekNext()))
	if s.peek() == '.' && s.IsDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.IsDigit(s.peek()) {
			s.advance()
		}
	}

	value := s.Source[s.start:s.current]
	number, _ := strconv.ParseFloat(value, 64)
	s.addTokenWithLiteral(NUMBER, number)
}

func (s *Scanner) peek() uint8 {
	if s.AtEnd() {
		return 0
	}
	return s.Source[s.current]
}

func (s *Scanner) peekNext() uint8 {
	if s.current+1 >= len(s.Source) {
		return 0
	}
	return s.Source[s.current+1]
}

func (s *Scanner) IsDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) IsAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) IsAlphaNumeric(c uint8) bool {
	return s.IsAlpha(c) || s.IsDigit(c)
}

func (s *Scanner) identifier() {
	for s.IsAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.Source[s.start:s.current]
	result, exits := keywords[text]
	if exits {
		s.addToken(result)
	} else {
		s.addToken(IDENTIFIER)
	}
}

func (s *Scanner) Eval(expr string) {
	s.reset()
	s.Source = expr
	s.ScanTokens()
}

func (s *Scanner) reset() {
	s.current = 0
	s.start = 0
	s.line = 1
	s.Source = ""
	s.Tokens = nil
}

package pip

import (
	"errors"

	"github.com/harryzcy/snuuze/types"
)

var (
	ErrInvalidSyntax = errors.New("invalid syntax")
)

type Parser struct {
	scanner *Scanner
	dep     *types.Dependency
}

func NewParser(file, line string) *Parser {
	return &Parser{
		scanner: NewScanner(line),
		dep: &types.Dependency{
			File:  file,
			Extra: map[string]interface{}{},
		},
	}
}

func (p *Parser) Parse() (*types.Dependency, error) {
	var err error
	p.dep.Name, err = p.parseDependencyName()
	if err != nil {
		return nil, err
	}

	tokenType, tokenValue, err := p.scanner.Scan()
	if err != nil {
		return nil, err
	}
	// optional extras
	if tokenType == TokenTypeBracketLeft {
		p.dep.Extra["extras"], err = p.parseExtras()
		if err != nil {
			return nil, err
		}

		tokenType, tokenValue, err = p.scanner.Scan()
		if err != nil {
			return nil, err
		}
	}

	if tokenType == TokenTypeOperator {
		var constraints [][2]string
		p.dep.Version, constraints, err = p.parseConstraints(tokenValue)
		if err != nil {
			return nil, err
		}

		p.dep.Extra["constraints"] = constraints
	}

	return p.dep, nil
}

func (p *Parser) parseDependencyName() (string, error) {
	tokenType, tokenValue, err := p.scanner.Scan()
	if err != nil {
		return "", err
	}

	if tokenType != TokenTypeIdentifier {
		return "", ErrInvalidSyntax
	}

	return tokenValue, nil
}

func (p *Parser) parseExtras() (string, error) {
	extras := ""
	tokenType, tokenValue, err := p.scanner.Scan()
	if err != nil {
		return "", err
	}
	if tokenType == TokenTypeIdentifier {
		extras += tokenValue
	} else {
		return "", ErrInvalidSyntax
	}

	for {
		tokenType, _, err := p.scanner.Scan()
		if err != nil {
			return "", err
		}

		if tokenType == TokenTypeBracketRight {
			break
		}

		if tokenType == TokenTypeComma {
			tokenType, tokenValue, err := p.scanner.Scan()
			if err != nil {
				return "", err
			}
			if tokenType != TokenTypeIdentifier {
				return "", ErrInvalidSyntax
			}
			extras += ", " + tokenValue
		}
	}

	return extras, nil
}

func (p *Parser) parseConstraints(initialOperator string) (string, [][2]string, error) {
	constraints := make([][2]string, 0)
	version := initialOperator

	tokenType, tokenValue, err := p.scanner.Scan()
	if err != nil {
		return "", nil, err
	}
	if tokenType != TokenTypeIdentifier {
		return "", nil, ErrInvalidSyntax
	}
	version += " " + tokenValue
	constraints = append(constraints, [2]string{initialOperator, tokenValue})

	for {
		tokenType, tokenValue, err = p.scanner.Scan()
		if err != nil {
			return "", nil, err
		}
		if tokenType != TokenTypeComma {
			break
		}

		if tokenType == TokenTypeComma {
			operator := ""
			tokenType, tokenValue, err = p.scanner.Scan()
			if err != nil {
				return "", nil, err
			}
			if tokenType != TokenTypeOperator {
				return "", nil, ErrInvalidSyntax
			}
			operator = tokenValue

			tokenType, tokenValue, err = p.scanner.Scan()
			if err != nil {
				return "", nil, err
			}
			if tokenType != TokenTypeIdentifier {
				return "", nil, ErrInvalidSyntax
			}
			version += ", " + operator + " " + tokenValue
			constraints = append(constraints, [2]string{operator, tokenValue})
		}
	}

	return version, constraints, nil
}

type TokenType int

const (
	TokenTypeError TokenType = iota
	TokenTypeEOT
	TokenTypeIdentifier
	TokenTypeBracketLeft
	TokenTypeBracketRight
	TokenTypeOperator
	TokenTypeComma
	TokenTypeSemicolon
)

type Scanner struct {
	line     string
	position int
}

func NewScanner(line string) *Scanner {
	return &Scanner{
		line:     line,
		position: 0,
	}
}

func (s *Scanner) Scan() (TokenType, string, error) {
	tokenType, tokenValue, err := s.scanToken()
	if err != nil {
		return TokenTypeError, "", err
	}

	return tokenType, tokenValue, nil
}

func (s *Scanner) scanToken() (TokenType, string, error) {
	for !s.isEOT() && s.isWhitespace() {
		s.position++
	}

	if s.isEOT() {
		return TokenTypeEOT, "", nil
	}
	if s.isLetterOrDigit() {
		return s.scanIdentifier()
	}
	if s.isOperator() {
		return s.scanOperator()
	}
	if s.line[s.position] == '[' {
		s.position++
		return TokenTypeBracketLeft, "[", nil
	}
	if s.line[s.position] == ']' {
		s.position++
		return TokenTypeBracketRight, "]", nil
	}
	if s.line[s.position] == ',' {
		s.position++
		return TokenTypeComma, ",", nil
	}
	if s.line[s.position] == ';' {
		s.position++
		return TokenTypeSemicolon, ";", nil
	}

	return TokenTypeEOT, "", ErrInvalidSyntax
}

func (s *Scanner) scanIdentifier() (TokenType, string, error) {
	identifier := ""
	char := s.line[s.position]
	for s.isLetterOrDigit() ||
		char == '_' || char == '-' || char == '.' || char == '+' || char == '*' || char == '!' {
		identifier += string(s.line[s.position])
		s.position++
		if s.isEOT() {
			break
		}
		char = s.line[s.position]
	}

	return TokenTypeIdentifier, identifier, nil
}

func (s *Scanner) scanOperator() (TokenType, string, error) {
	operator := ""
	for !s.isEOT() && s.isOperator() {
		operator += string(s.line[s.position])
		s.position++
	}

	return TokenTypeOperator, operator, nil
}

func (s *Scanner) isLetter() bool {
	char := s.line[s.position]
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func (s *Scanner) isDigit() bool {
	char := s.line[s.position]
	return char >= '0' && char <= '9'
}

func (s *Scanner) isLetterOrDigit() bool {
	return s.isLetter() || s.isDigit()
}

func (s *Scanner) isOperator() bool {
	char := s.line[s.position]
	return char == '=' || char == '!' || char == '<' || char == '>' || char == '~'
}

func (s *Scanner) isWhitespace() bool {
	char := s.line[s.position]
	return char == ' ' || char == '\t'
}

func (s *Scanner) isEOT() bool {
	return s.position >= len(s.line)
}

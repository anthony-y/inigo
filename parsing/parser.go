package parsing

import (
	"fmt"
	"strconv"
)

// Rough Backus Naur form
/*
program ::= assignment
         | section
         ;

value ::= identifier
		| STRING
		| NUMBER
		;

assignment ::= identifier "=" value

section ::= "[" identifier "]"
*/

type (
	iniParser struct {
		tokens   []Token
		tokenPos int

		line int

		//currentSection string
	}

	// ParseError represents an error that occurrs during INI parsing
	ParseError struct {
		Message string
		Line    int
	}
)

var currentSection = ""

func (pe ParseError) Error() string {
	return fmt.Sprintf("INIGO: %s (line %d)", pe.Message, pe.Line)
}

func (i *iniParser) previous() Token {
	if i.tokenPos == 0 {
		return i.tokens[0]
	}

	return i.tokens[i.tokenPos-1]
}

func (i *iniParser) next() {
	i.tokenPos++
}

func (i *iniParser) current() Token {
	return i.tokens[i.tokenPos]
}

func (i *iniParser) peek() Token {
	return i.tokens[i.tokenPos+1]
}

func (i *iniParser) makeError(message string) ParseError {
	return ParseError{
		Message: message,
		Line:    i.line,
	}
}

// ParseLine parses a line of ini into an expression
func ParseLine(tokens []Token) (Expression, error) {
	parser := &iniParser{
		tokens:   tokens,
		tokenPos: 0,
		line:     tokens[0].LineNum,
	}

	e, err := program(parser)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func program(i *iniParser) (Expression, error) {
	tok := i.current()

	if tok.Type == OpenBrace {
		return section(i)
	}

	if tok.Type == Comment {
		return nil, nil
	}

	if tok.Type == LineBreak {
		return nil, nil
	}

	return value(i)
}

func value(i *iniParser) (Expression, error) {
	tok := i.current()

	if tok.Type == Identifier {
		return identifier(i)
	}

	if tok.Type == StringLiteral {
		return stringLiteral(i)
	}

	if tok.Type == NumberLiteral {
		return numberLiteral(i)
	}

	if tok.Type == FloatLiteral {
		return floatLiteral(i)
	}

	return nil, i.makeError("Unknown token")
}

func section(i *iniParser) (Expression, error) {
	i.next()

	expr, err := program(i)
	if err != nil {
		return nil, err
	}

	ident, ok := expr.(identifierExpression)
	if !ok {
		return nil, ParseError{
			Message: "Expected identifier on section header",
			Line:    i.line,
		}
	}

	if i.current().Type != CloseBrace {
		return nil, i.makeError("Expected close brace in section declaration")
	}

	currentSection = ident.text
	return SectionExpression{
		Name: ident.text,
	}, nil
}

func assignment(i *iniParser) (Expression, error) {
	name := i.previous().Text

	// Skip assign
	i.next()

	value, err := program(i)
	if err != nil {
		return nil, err
	}

	switch value.(type) {

	// Value must be an identifier, a string, a float or an int
	case identifierExpression,
		stringLiteralExpression,
		floatLiteralExpression,
		numberLiteralExpression:

		return AssignmentExpression{
			Name:    name,
			Value_:  value,
			Section: currentSection,
		}, nil
	}

	return nil, i.makeError("Expecting value on key")
}

func identifier(i *iniParser) (Expression, error) {
	// Catch random identifiers
	if i.previous().Type == i.current().Type {
		if i.peek().Type == LineBreak && i.current().Type != Comment {
			return nil, i.makeError("Unexpected identifier")
		}
	}

	i.next()

	if i.current().Type == Assign {
		return assignment(i)
	}

	if i.current().Type == CloseBrace || i.current().Type == LineBreak {
		return identifierExpression{
			text: i.previous().Text,
		}, nil
	}

	return nil, i.makeError("Invalid identifier")
}

func stringLiteral(i *iniParser) (Expression, error) {
	i.next()

	return stringLiteralExpression{
		text: i.previous().Text,
	}, nil
}

func numberLiteral(i *iniParser) (Expression, error) {
	i.next()

	number, err := strconv.Atoi(i.previous().Text)
	if err != nil {
		return nil, i.makeError("Invalid int literal")
	}

	return numberLiteralExpression{
		number: number,
		asText: i.previous().Text,
	}, nil
}

func floatLiteral(i *iniParser) (Expression, error) {
	i.next()

	float, err := strconv.ParseFloat(i.previous().Text, 64)
	if err != nil {
		return nil, i.makeError("Invalid float literal")
	}

	return floatLiteralExpression{
		number: float,
		asText: i.previous().Text,
	}, nil
}

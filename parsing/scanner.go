package parsing

import (
	"unicode"
)

type (
	// TokenType represents a type of token
	TokenType int

	// Token represents a lexical token
	Token struct {
		// See bottom of file
		Type TokenType

		Text    string
		LineNum int
	}
)

var (
	symbols = map[rune]TokenType{
		'[': OpenBrace,
		']': CloseBrace,

		'=': Assign,

		'\n': LineBreak,
	}
)

// Scan splits a rune alice of ini into a slice of tokens
func Scan(input []rune) ([]Token, []error) {
	// Output slice
	out := []Token{}
	errors := []error{}

	// Index of current rune in input slice
	pos := 0
	line := 1

	// For every character in the input
	for {
		// Exit condition
		if pos >= len(input)-1 {
			break
		}

		// Skip whitespace
		for input[pos] == ' ' || input[pos] == '\t' {
			pos++
		}

		// New token to be appended to the output slice at the end of the loop
		token := Token{}
		token.LineNum = line

		// Check for all single character tokens
		// If current character is in the smybols map...
		if _, ok := symbols[input[pos]]; ok {
			text := string(input[pos])
			if input[pos] == '\n' {
				text = "\\n"
			}

			// ...create a token from it
			token.Type = symbols[input[pos]]
			token.Text = text

		} else if unicode.IsLetter(input[pos]) || input[pos] == '_' { // Identifiers
			startPos := pos
			for unicode.IsLetter(input[pos+1]) || unicode.IsNumber(input[pos+1]) || input[pos+1] == '_' || input[pos+1] == ' ' {
				if input[pos] == '\n' {
					errors = append(errors, IniError{
						"Unexpected line-break in identifier",
						line,
					})
				}

				pos++
			}

			token.Type = Identifier
			token.Text = string(input[startPos : pos+1])

		} else if unicode.IsNumber(input[pos]) || input[pos] == '-' { // Numbers
			startPos := pos
			numType := NumberLiteral

			for unicode.IsNumber(input[pos]) || input[pos] == '.' {
				if input[pos] == '\n' {
					errors = append(errors, IniError{
						"Unexpected line-break in number",
						line,
					})
				}

				if input[pos] == '.' {
					numType = FloatLiteral
				}

				pos++
			}

			token.Type = numType
			token.Text = string(input[startPos : pos+1])

		} else if input[pos] == '"' { // Strings
			pos++

			startPos := pos
			for input[pos] != '"' {
				if input[pos] == '\n' {
					errors = append(errors, IniError{
						"Unterminated string",
						line,
					})
				}

				pos++
			}

			token.Type = StringLiteral
			token.Text = string(input[startPos:pos])

		} else if input[pos] == '#' || input[pos] == ';' { // Comments
			for input[pos+1] != '\n' {
				pos++
			}

			token.Type = Comment
			token.Text = ""

		} else {
			token.Type = Unknown
			token.Text = string(input[pos])
		}

		out = append(out, token)
		pos++
	}

	if len(errors) > 0 {
		return out, errors
	}

	return out, nil
}

// TokenType "enum"
const (
	Identifier TokenType = iota
	Assign

	StringLiteral
	FloatLiteral
	NumberLiteral

	OpenBrace
	CloseBrace

	LineBreak
	Comment

	Unknown
)

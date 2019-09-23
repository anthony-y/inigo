package parsing

import (
	"unicode"
)

type (
	// TokenType represents a type of token
	TokenType int

	// Token represents a lexical token
	Token struct {
		Type TokenType
		Text string

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

// ScanLine splits a single line of ini into a slice of tokens
func ScanLine(line []rune, num int) []Token {
	// Output slice
	out := []Token{}

	// Index of current rune in input slice
	pos := 0

	// For every character in the line
	for {
		// Exit condition
		if pos >= len(line) {
			break
		}

		// Skip whitespace
		for line[pos] == ' ' || line[pos] == '\t' {
			pos++
		}

		// New token to be appended to the output slice at the end of the loop
		token := Token{}
		token.LineNum = num

		// Check for all single character tokens
		// If current character is in the smybols map...
		if _, ok := symbols[line[pos]]; ok {
			text := string(line[pos])
			if line[pos] == '\n' {
				text = "\\n"
			}

			// ...create a token from it
			token.Type = symbols[line[pos]]
			token.Text = text

		} else if unicode.IsLetter(line[pos]) || line[pos] == '_' { // Identifiers
			startPos := pos
			for unicode.IsLetter(line[pos+1]) || unicode.IsNumber(line[pos+1]) || line[pos+1] == '_' || line[pos+1] == ' ' {
				pos++
			}

			token.Type = Identifier
			token.Text = string(line[startPos : pos+1])

		} else if unicode.IsNumber(line[pos]) || line[pos] == '-' { // Numbers
			startPos := pos
			numType := NumberLiteral

			for unicode.IsNumber(line[pos+1]) || line[pos+1] == '.' {
				if line[pos+1] == '.' {
					numType = FloatLiteral
				}

				pos++
			}

			token.Type = numType
			token.Text = string(line[startPos : pos+1])

		} else if line[pos] == '"' { // Strings
			pos++

			startPos := pos
			for line[pos] != '"' {
				if line[pos] == '\n' {
					panic("Unterminated string")
				}

				pos++
			}

			token.Type = StringLiteral
			token.Text = string(line[startPos:pos])

		} else if line[pos] == '#' || line[pos] == ';' { // Comments
			for line[pos+1] != '\n' {
				pos++
			}

			token.Type = Comment
			token.Text = ""

		} else {
			token.Type = Unknown
			token.Text = string(line[pos])
		}

		out = append(out, token)
		pos++
	}

	return out
}

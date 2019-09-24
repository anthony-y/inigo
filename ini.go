package inigo

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/anthony-y/inigo/parsing"
)

type File map[string]Section
type Section map[string]interface{}

func (f File) Section(name string) Section {
	return f[name]
}

func (s Section) Key(name string) interface{} {
	return s[name]
}

func LoadIniFile(path string) (File, []error) {
	handle, err := os.Open(path)
	if err != nil {
		return nil, []error{
			err,
		}
	}
	defer handle.Close()

	return LoadIni(handle)
}

func LoadIniFromBytes(b []byte) (File, []error) {
	return LoadIni(bytes.NewReader(b))
}

func LoadIni(reader io.Reader) (File, []error) {
	ini := File{}

	lineReader := bufio.NewScanner(reader)
	lineNum := 0

	expressions := []parsing.Expression{}
	var errors []error

	for lineReader.Scan() {
		lineNum++
		line := []rune(lineReader.Text() + "\n")

		// Ignore blank lines
		tokens := parsing.ScanLine(line, lineNum)
		if len(tokens) == 1 && tokens[0].Type == parsing.LineBreak {
			continue
		}

		// Parse
		expression, err := parsing.ParseLine(tokens)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		expressions = append(expressions, expression)
	}

	if err := lineReader.Err(); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return ini, errors
	}

	initINIFromAST(&ini, expressions)

	return ini, nil
}

func initINIFromAST(file *File, ast []parsing.Expression) {
	for _, expression := range ast {
		switch v := expression.(type) {
		case parsing.AssignmentExpression:
			if (*file)[v.Section] == nil {
				(*file)[v.Section] = Section{}
			}

			(*file)[v.Section][v.Name] = v.Value()
		}
	}
}

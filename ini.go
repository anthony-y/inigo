package inigo

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/anthony-y/inigo/parsing"
)

type (
	File    map[string]Section
	Section map[string]interface{}

	ReadError struct {
		Name string
		Err  error
	}
)

func (re ReadError) Error() string {
	return fmt.Sprintf("Failed to read %s: %s", re.Name, re.Err.Error())
}

func (f File) Section(name string) Section {
	return f[name]
}

func (s Section) Key(name string) interface{} {
	return s[name]
}

func LoadIniFile(path string) (f File, e []error) {
	handle, err := os.Open(path)
	if err != nil {
		return nil, []error{
			err,
		}
	}

	f, e = LoadIni(handle)
	handle.Close()

	return
}

func LoadIniFromBytes(b []byte) (File, []error) {
	return LoadIni(bytes.NewReader(b))
}

func LoadIni(reader io.Reader) (File, []error) {
	ini := File{}
	errors := []error{}

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return nil, []error{
			ReadError{
				Name: "io.Reader",
				Err:  err,
			},
		}
	}

	tokens, lexErrs := parsing.Scan([]rune(buf.String() + "\n"))
	if lexErrs != nil {
		errors = append(errors, lexErrs...)
	}

	expressions, parseErrs := parsing.Parse(tokens)
	if parseErrs != nil || len(parseErrs) > 0 {
		errors = append(errors, parseErrs...)
	}

	if len(errors) > 0 {
		return nil, errors
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

package inigo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
)

type (
	IniFile    map[string]IniSection  // name -> section
	IniSection map[string]interface{} // name -> variable

	ReadError struct {
		Name string
		Err  error
	}
)

func (re ReadError) Error() string {
	return fmt.Sprintf("failed to read %s: %s", re.Name, re.Err.Error())
}

func readLines(raw io.Reader) []string {
	lines := []string{}

	scanner := bufio.NewScanner(raw)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func readSectionHeader(ini IniFile, line []rune) (string, error) {
	cursor := 0
	for line[cursor] != ']' && line[cursor] != rune(0) && line[cursor] != '\n' && cursor < len(line) {
		cursor++
	}

	sectionName := string(line[0:cursor])

	if _, exists := ini[sectionName]; !exists {
		ini[sectionName] = make(map[string]interface{})
	} else {
		return "", errors.New(fmt.Sprintf("section %s already exists", sectionName))
	}

	return sectionName, nil
}

func readVariable(ini IniFile, currentSection string, line []rune) error {
	cursor := 0
	for line[cursor] != '=' && line[cursor] != rune(0) && line[cursor] != '\n' && cursor < len(line) {
		cursor++
	}
	variableName := string(line[0:cursor])

	if line[cursor] != '=' {
		return errors.New("expected '='")
	}

	value := string(line[cursor+1:])
	ini[currentSection][variableName] = value

	return nil
}

func ReadIni(raw io.Reader) (IniFile, []error) {
	var out IniFile = make(map[string]IniSection)
	errors := []error{}
	lines := readLines(raw)
	currentSection := ""

	for _, line := range lines {
		var err error = nil

		if len(line) == 0 {
			continue
		}

		if line[0] == '[' {
			currentSection, err = readSectionHeader(out, []rune(line[1:]))
		} else if line[0] == ';' || line[0] == '#' {
			continue
		} else {
			err = readVariable(out, currentSection, []rune(line))
		}

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return out, nil
	} else {
		return nil, errors
	}
}

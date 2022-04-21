package inigo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

type (
	// You can access sections as if they're JavaScript/JSON objects
	// e.g given this ini:
	//  [Section]
	//  value = 10
	//
	//  ini, _ := inigo.LoadIni(...)
	//  ini["Section"]["value"] == 10
	//
	IniFile    map[string]IniSection  // name -> section
	IniSection map[string]interface{} // name -> variable
)

// Load some ini data into memory
func LoadIni(raw io.Reader) (IniFile, []error) {
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
			continue // skip comment lines
		} else {
			err = readVariable(out, currentSection, []rune(line))
		}

		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return out, nil
	}

	return nil, errors
}

// Send formatted output of the entire ini state to an io.Writer
func (this IniFile) WriteTo(w io.Writer) (totalBytes int64, err error) {
	sectionsCount := 0

	for sectionName, fields := range this {
		sectionsNumBytes, err := w.Write([]byte(fmt.Sprintf("[%s]\n", sectionName)))
		totalBytes += int64(sectionsNumBytes)
		if err != nil {
			return int64(sectionsNumBytes), err
		}

		for fieldName, value := range fields {
			fieldsNumBytes, err := w.Write([]byte(fmt.Sprintf("%s=", fieldName)))
			totalBytes += int64(fieldsNumBytes)
			if err != nil {
				return int64(totalBytes), err
			}

			valuesNumBytes := 0

			switch value.(type) {
			case int:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("%d\n", value.(int))))
			case int32:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("%d\n", value.(int32))))
			case int64:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("%d\n", value.(int64))))
			case string:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("\"%s\"\n", value.(string))))
			case float32:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("%f\n", value.(float32))))
			case float64:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("%f\n", value.(float64))))
			case bool:
				valuesNumBytes, err = w.Write([]byte(fmt.Sprintf("%t\n", value.(bool))))
			}

			totalBytes += int64(valuesNumBytes)
			if err != nil {
				return int64(fieldsNumBytes + sectionsNumBytes + valuesNumBytes), err
			}
		}

		if sectionsCount < len(this)-1 {
			lineBreakNumBytes, err := w.Write([]byte{'\n'})
			totalBytes += int64(lineBreakNumBytes)
			if err != nil {
				return totalBytes, err
			}
		}

		sectionsCount++
	}

	return
}

// Given a correct declaration of a section, initializes one inside the given IniFile state
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

// Given a correct field declaration, initializes one with the given value
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

	err := parseAndStoreValue(ini[currentSection], variableName, []rune(value))
	if err != nil {
		return err
	}

	return nil
}

// Convert `stringValue` into an appropriate datatype, and store it under `fieldName` inside given section `ini`
func parseAndStoreValue(ini IniSection, fieldName string, stringValue []rune) error {
	first := stringValue[0]

	if unicode.IsNumber(first) || first == '-' || first == '+' {
		floating := false
		for _, ch := range stringValue {
			if ch == '.' {
				floating = true
				break
			}
		}

		if floating {
			f, err := strconv.ParseFloat(string(stringValue), 64)
			if err != nil {
				return err
			}
			ini[fieldName] = f
			return nil
		}

		number, err := strconv.ParseInt(string(stringValue), 0, 64)
		if err != nil {
			return err
		}
		ini[fieldName] = number
		return nil
	}

	if string(stringValue) == "true" {
		ini[fieldName] = true
		return nil
	}

	if string(stringValue) == "false" {
		ini[fieldName] = false
		return nil
	}

	// Ignore enclosing quotations for text values
	if first == '"' || first == '\'' {
		if stringValue[len(stringValue)-1] != '"' && stringValue[len(stringValue)-1] != '\'' {
			return errors.New("unclosed string literal, did you forget a '\"'?")
		}
		stringValue = stringValue[1 : len(stringValue)-1]
	}

	ini[fieldName] = string(stringValue)
	return nil
}

// Using a bufio.Scanner, accumulate all lines in the input into the returned slice
func readLines(raw io.Reader) []string {
	lines := []string{}

	scanner := bufio.NewScanner(raw)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

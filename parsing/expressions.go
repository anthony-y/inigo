package parsing

type (
	// Expression is any expression in the INI file
	Expression interface {
		Value() interface{}
		String() string
	}

	// AssignmentExpression represents an assignment in an INI file
	AssignmentExpression struct {
		Name    string
		Section string
		Value_  Expression
	}

	// SectionExpression represents a section in an INI file
	SectionExpression struct {
		Name string
	}

	// identifierExpression represents a name in an INI file
	identifierExpression struct {
		text string
	}

	// stringLiteralExpression represents a string value in an INI file
	stringLiteralExpression struct {
		text string
	}

	// numberLiteralExpression represents an integer value in an INI file
	numberLiteralExpression struct {
		number int
		asText string
	}

	// floatLiteralExpression represents a float value in an INI file
	floatLiteralExpression struct {
		number float64
		asText string
	}
)

// String()
func (as AssignmentExpression) String() string {
	return as.Section + "! Assignment: " + as.Name + " = " + as.Value_.String()
}

func (se SectionExpression) String() string {
	return se.Name
}

func (ie identifierExpression) String() string {
	return ie.text
}

func (sle stringLiteralExpression) String() string {
	return sle.text
}

func (nle numberLiteralExpression) String() string {
	return nle.asText
}

func (fle floatLiteralExpression) String() string {
	return fle.asText
}

// Value returns the true value of an assignment
func (as AssignmentExpression) Value() interface{} {
	return as.Value_.Value()
}

// Value returns the formatted name of a section
func (se SectionExpression) Value() interface{} {
	return "[" + se.Name + "]"
}

func (ie identifierExpression) Value() interface{} {
	return ie.text
}

func (sle stringLiteralExpression) Value() interface{} {
	return sle.text
}

func (nle numberLiteralExpression) Value() interface{} {
	return nle.number
}

func (fle floatLiteralExpression) Value() interface{} {
	return fle.number
}

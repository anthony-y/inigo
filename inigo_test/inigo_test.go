package inigo_test

import (
	"fmt"
	"testing"

	"inigo"
)

type test struct {
	section string
	key     string

	value interface{}
}

var tests = []test{
	{"", "global", "global string"},
	{"", "anotherglobal", "another one"},
	{"testsection", "local", -10},
	{"testsection", "yadig", 100.12},
}

var ini inigo.File

func TestLoad(t *testing.T) {
	var errs []error
	ini, errs = inigo.LoadIni("../testdata/test.ini")
	if errs != nil {
		for _, e := range errs {
			fmt.Println(e)
		}
		t.Fail()
	}
}

func TestFields(t *testing.T) {
	for _, test := range tests[:2] {
		expectedValue := test.value.(string)

		if ini[test.section][test.key].(string) != expectedValue {
			t.Errorf(
				"Expected value %v for ini[\"%s\"][\"%s\"], got %v",
				test.value,
				test.section,
				test.key,
				ini[test.section][test.key],
			)
		}
	}

	local := ini[tests[2].section][tests[2].key].(int)
	//panic(local)
	if local != tests[2].value.(int) {
		t.Fail()
	}

	yadig := ini[tests[3].section][tests[3].key].(float64)
	if yadig != tests[3].value.(float64) {
		t.Fail()
	}
}

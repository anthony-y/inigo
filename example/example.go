package main

import (
	"fmt"
	"github.com/anthony-y/inigo"
)

func main() {
	ini, errs := inigo.LoadIni("example.ini")
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	globals := ini[""]["globals"].(string)
	fmt.Println(globals)

	for _, line := range ini["lines"] {
		fmt.Println(line.(string))
	}

	dataSection := ini["data"]
	astring := dataSection["astring"].(string)
	aint := dataSection["aint"].(int)
	afloat := dataSection["afloat"].(float64)

	fmt.Printf("astring=\"%s\"\naint=%d\nafloat=%f\n", astring, aint, afloat)
}

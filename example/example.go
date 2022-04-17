package main

import (
	"fmt"
	"log"
	"os"

	"github.com/anthony-y/inigo"
)

func main() {
	handle, err := os.Open("example.ini")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	ini, errs := inigo.Ini(handle)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	for sectionName, fields := range ini {
		fmt.Printf("[%s]\n", sectionName)
		for fieldName, value := range fields {
			fmt.Printf("%s=%s\n", fieldName, value)
		}
	}
}

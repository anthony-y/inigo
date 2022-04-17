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

	ini, errs := inigo.LoadIni(handle)
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	ini.Write(os.Stdout)
}

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

	ini["Graphics"]["lod_level"] = 5

	writer, err := os.Create("example.ini")
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	count, err := ini.WriteTo(writer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Wrote", count, "bytes to", writer)
}

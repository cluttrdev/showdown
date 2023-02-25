package main

import (
	"flag"
	"log"

	"github.com/pkg/errors"
)

func main() {
	file := flag.String("file", "", "The file to preview")
	//port := flag.String("port", "1337", "HTTP network address")

	flag.Parse()

	if *file == "" {
		log.Fatal("No file given!")
		return
	}

	app := NewApplication(*file)

	err := app.run()
	log.Fatal(errors.Errorf("Stopped listening: %v", err))
}

package main

import (
	"log"
	"os"

	report "github.com/muravjov/goroutinereport"
)

func main() {
	file := os.Stdin
	if len(os.Args) > 1 {
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		file = f
	}
	if err := report.Report(file, os.Stdout); err != nil {
		log.Println(err)
	}
}

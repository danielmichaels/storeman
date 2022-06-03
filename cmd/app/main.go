package main

import (
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln("failed to start storeman")
	}
}

func run() error {
	println("storeman")
	return nil
}

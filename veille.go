package main

import (
	"log"
)

func main() {
	log.Println("Loading config")
	err := LoadConfig()
	if err != nil {
		panic(err)
	}

	log.Println(Config)
}

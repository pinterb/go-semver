package main

import (
	"fmt"
	"log"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	tags, err := tags()
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

	for _, tag := range tags {
		fmt.Println(tag)
	}
}

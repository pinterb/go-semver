package main

import (
	"log"

	"github.com/pinterb/go-semver/internal/cli"
)

func main() {
	if err := cli.New().Execute(); err != nil {
		log.Fatal("error during command execution: %v", err)
	}
}

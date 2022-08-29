package main

import (
	"fmt"
	"log"
	"os"

	"github.com/damienpontifex/azdo-template-docs/cmd"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Please provide file path in first argument")
		os.Exit(1)
	}

	filePath := os.Args[1]
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file at %v", filePath)
		os.Exit(1)
	}

	if parameters, err := cmd.Parse(file); err != nil {
		fmt.Printf("%v", parameters)
	}
}

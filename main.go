package main

import (
	"os"

	"github.com/damienpontifex/azdo-template-docs/cmd"
)

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatalf("Please provide file path in first argument")
	// 	os.Exit(1)
	// }

	// filePath := os.Args[1]
	file, err := os.ReadFile("/Users/ponti/wooliesx/vsts-build-templates/graphql/check-schema.yaml")
	if err != nil {
		// log.Fatalf("Failed to read file at %v", filePath)
		os.Exit(1)
	}

	parameters, err := cmd.Parse(file)
	// if err != nil {
	// fmt.Printf("%v", parameters)
	parameters.ToMarkdownTable()
	// }
}

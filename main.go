package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/damienpontifex/azdo-template-docs/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var OutputFile string
	var cli = &cobra.Command{
		Use:   "azdo-template-docs [input file]",
		Short: "Generate markdown table for input parameters of an Azure DevOps template file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cobraCommand *cobra.Command, args []string) error {
			filePath := args[0]
			if strings.HasPrefix(filePath, "~/") {
				home, _ := os.UserHomeDir()
				filePath = filepath.Join(home, filePath[2:])
			}

			file, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			parameters, err := cmd.Parse([]byte(file))
			if err != nil {
				return err
			}

			writer := os.Stdout
			if len(OutputFile) > 0 {
				writer, err = os.Create(OutputFile)
				if err != nil {
					return err
				}
			}

			parameters.ToMarkdownTable(writer)

			return nil
		},
	}

	cli.Flags().StringVarP(&OutputFile, "output-file", "o", "", "Output file, if unused stdout will be used")
	cli.Execute()
}

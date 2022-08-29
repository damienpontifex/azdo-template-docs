package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

type Parameter struct {
	Description string
	Name        string
	Type        string
	Default     *string
}

type AzDoTemplate struct {
	Parameters []Parameter
}

func Parse(s []byte) (*AzDoTemplate, error) {
	var t yaml.Node

	err := yaml.Unmarshal(s, &t)
	if err != nil {
		return nil, err
	}

	allParameters := make([]Parameter, 0)
	yamlParameters := t.Content[0].Content[1]
	for i := 0; i < len(yamlParameters.Content); i++ {
		content := yamlParameters.Content[i]
		parameter := Parameter{}
		_ = content.Decode(&parameter)
		parameter.Description = strings.ReplaceAll(content.HeadComment, "# ", "")
		allParameters = append(allParameters, parameter)
	}
	return &AzDoTemplate{Parameters: allParameters}, err
}

func (t *AzDoTemplate) ToMarkdownTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Type", "Default", "Required"})
	table.SetRowLine(true)
	for _, v := range t.Parameters {
		var defaultDescription string
		if v.Default != nil {
			defaultDescription = *v.Default
		}
		table.Append([]string{v.Name, v.Description, v.Type, defaultDescription, fmt.Sprintf("%v", v.Default == nil)})
	}
	table.Render()
}

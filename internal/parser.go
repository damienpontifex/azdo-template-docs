package internal

import (
	"fmt"
	"io"
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

	// Find the parameters node and then get the next one as the array
	// of parameter values
	var yamlParameters *yaml.Node
	for i, v := range t.Content[0].Content {
		if v.Value == "parameters" && len(t.Content[0].Content) > i+1 {
			yamlParameters = t.Content[0].Content[i+1]
		}
	}

	if yamlParameters == nil {
		// No parameters found in template file
		return &AzDoTemplate{Parameters: make([]Parameter, 0)}, nil
	}

	// Transform yaml nodes to parameter struct
	for i := 0; i < len(yamlParameters.Content); i++ {
		content := yamlParameters.Content[i]
		parameter := Parameter{}
		_ = content.Decode(&parameter)
		parameter.Description = strings.ReplaceAll(content.HeadComment, "# ", "")
		allParameters = append(allParameters, parameter)
	}
	return &AzDoTemplate{Parameters: allParameters}, err
}

func (t *AzDoTemplate) ToMarkdownTable(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
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

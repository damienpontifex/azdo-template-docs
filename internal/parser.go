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

type SimpleParameter struct {
	Parameters map[string][]string
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
	yamlMap := t.Content[0].Content
	for i, v := range yamlMap {
		if v.Value == "parameters" && len(yamlMap) > i+1 {
			yamlParameters = yamlMap[i+1]
		}
	}

	if yamlParameters == nil || len(yamlParameters.Content) == 0 {
		// No parameters found in template file
		return &AzDoTemplate{Parameters: make([]Parameter, 0)}, nil
	}

	if len(yamlParameters.Content[0].Content) == 0 {
		// Key-value parameters
		for i := 0; i < len(yamlParameters.Content); i += 2 {
			key := yamlParameters.Content[i]
			value := yamlParameters.Content[i+1]
			parameter := Parameter{
				Name:        key.Value,
				Type:        typeFromTag(key.Tag),
				Description: key.HeadComment,
				Default:     &value.Value,
			}

			allParameters = append(allParameters, parameter)
		}
	} else {
		// Array of objects as parameters
		// Transform yaml nodes to parameter struct
		for _, content := range yamlParameters.Content {
			// Array of parameters with keys for metadata
			// content := yamlParameters.Content[i]
			parameter := Parameter{}
			_ = content.Decode(&parameter)
			parameter.Description = strings.ReplaceAll(content.HeadComment, "# ", "")

			allParameters = append(allParameters, parameter)
		}
	}
	return &AzDoTemplate{Parameters: allParameters}, err
}

func (t *AzDoTemplate) ToMarkdownTable(writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Name", "Description", "Type", "Default", "Required"})
	table.SetAutoWrapText(false)
	table.SetColumnSeparator("|")
	table.SetCenterSeparator("|")
	table.SetBorders(tablewriter.Border{Top: false, Bottom: false, Left: true, Right: true})

	for _, v := range t.Parameters {
		var defaultDescription string
		if v.Default != nil {
			defaultDescription = *v.Default
		}

		table.Append([]string{v.Name, strings.ReplaceAll(v.Description, "\n", "<br/>"), v.Type, defaultDescription, fmt.Sprintf("%v", v.Default == nil)})
	}
	table.Render()
}

func typeFromTag(tag string) string {
	switch tag {
	case "!!str":
		return "string"
	case "!!bool":
		return "boolean"
	}
	return ""
}

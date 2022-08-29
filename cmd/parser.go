package cmd

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type Parameter struct {
	Description string
	Name        string
	Type        string
}

func Parse(s []byte) ([]Parameter, error) {
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
	return allParameters, err
}

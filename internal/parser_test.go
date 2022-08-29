package internal

import (
	"bytes"
	"testing"
)

var data = `
parameters:
# If grabbing the SDL from a local file. If the schemaIntrospectionUrl parameter is also set
# then this will be the file in the default working directory the schema will be
- name: localSchemaFile
  type: string
# If wanting to grab the SDL from a running server
- name: localSchemaIntrospectionUrl
  type: string
  default: ""
# If wanting to start a local server to introspect the schema
- name: commandToStartLocalServer
  type: string
  default: ""
# The graph from Apollo Studio
- name: graph
  type: string
# Apollo studio API Key
- name: apikey
  type: string
# The variant of the graph in Apollo Studio
- name: variant
  type: string
# The service name that appears in Apollo Studio
- name: serviceName
  type: string
# Endpoint that service is available on when deployed
- name: introspectionURL
  type: string
  default: ""
- name: options
  type: string
  default: ""
# Whether to use the Apollo Rover CLI instead of the npm module
- name: useRover
  type: boolean
  default: false
# Whether to fail the pipeline if the linting fails
- name: failOnLintChecks
  type: boolean
  default: false
# Whether to publish the schema file
- name: publishSchema
  type: boolean
  default: false

steps:
- ${{ if eq(parameters.useRover, true) }}:
  - template: ./install-rover-cli.yaml

  - script: |
      if [[ -n "${{ parameters.commandToStartLocalServer }}" ]]; then
        ${{ format('{0} &', coalesce(parameters.commandToStartLocalServer, 'echo "no server to run"')) }}
        sleep 2
      fi

      if [[ -n "${{ parameters.localSchemaIntrospectionUrl }}" ]]; then
        rover subgraph introspect "${{ parameters.localSchemaIntrospectionUrl }}" > "${{ parameters.localSchemaFile }}"
      fi
      rover subgraph check ${{ parameters.graph }}@${{ parameters.variant }} \
        --name ${{ parameters.serviceName }} \
        --schema ${{ parameters.localSchemaFile }}

      if [[ -n "${{ parameters.commandToStartLocalServer }}" ]]; then
        kill $!
      fi
    env:
      APOLLO_KEY: ${{ parameters.apikey }}
    displayName: Validate schema with federated graph

- ${{ else }}:
  - bash: |
     npx apollo service:check \
         --localSchemaFile="${{ parameters.localSchemaFile }}" \
         --graph="${{ parameters.graph }}" \
         --key="${{ parameters.apikey }}" \
         --variant="${{ parameters.variant }}" \
         --serviceName="${{ parameters.serviceName }}" \
         --endpoint="${{ parameters.introspectionURL }}" ${{ parameters.options }}
    displayName: 'Validate schema with federated graph'

- script: |
    cd $(mktemp -d)
    npm install graphql-schema-linter graphql
    cat <<EOF > federation.graphql
      # lint-disable types-have-descriptions
      scalar _FieldSet
      directive @external on FIELD_DEFINITION
      directive @requires(fields: _FieldSet!) on FIELD_DEFINITION
      directive @provides(fields: _FieldSet!) on FIELD_DEFINITION
      directive @key(fields: _FieldSet!) repeatable on OBJECT | INTERFACE

      # this is an optional directive discussed below
      directive @extends on OBJECT | INTERFACE
    EOF
    if [[ "${{ parameters.localSchemaFile }}" = /* ]]; then
      LOCAL_FILE="${{ parameters.localSchemaFile }}"
    else
      LOCAL_FILE="${SYSTEM_DEFAULTWORKINGDIRECTORY}/${{ parameters.localSchemaFile }}"
    fi
    ./node_modules/.bin/graphql-schema-linter "${LOCAL_FILE}" federation.graphql \
        --ignore '{"types-have-descriptions":["Query"],"fields-have-descriptions":["ProductCategoriesConnection.totalCount"]}' \
        --rules arguments-have-descriptions,defined-types-are-used,deprecations-have-a-reason,descriptions-are-capitalized,enum-values-all-caps,enum-values-have-descriptions,fields-are-camel-cased,fields-have-descriptions,input-object-values-are-camel-cased,input-object-values-have-descriptions,relay-connection-types-spec,relay-connection-arguments-spec,types-are-capitalized,types-have-descriptions
  displayName: Lint schema
  continueOnError: ${{ not(parameters.failOnLintChecks) }}


- ${{ if eq(parameters.publishSchema, true) }}:
  - publish: ${{ parameters.localSchemaFile }}
    artifact: schema
    displayName: Publish schema
`

func TestParser(t *testing.T) {
	template, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	if len(template.Parameters) != 12 {
		t.Errorf("Length of parameters should be 12, got %v", len(template.Parameters))
	}

	expectedParameters := []Parameter{
		{
			Description: "If grabbing the SDL from a local file. If the schemaIntrospectionUrl parameter is also set\nthen this will be the file in the default working directory the schema will be",
			Name:        "localSchemaFile",
			Type:        "string",
		},
		{
			Description: "If wanting to grab the SDL from a running server",
			Name:        "localSchemaIntrospectionUrl",
			Type:        "string",
		},
		{
			Description: "If wanting to start a local server to introspect the schema",
			Name:        "commandToStartLocalServer",
			Type:        "string",
		},
		{
			Description: "The graph from Apollo Studio",
			Name:        "graph",
			Type:        "string",
		},
		{
			Description: "Apollo studio API Key",
			Name:        "apikey",
			Type:        "string",
		},
		{
			Description: "The variant of the graph in Apollo Studio",
			Name:        "variant",
			Type:        "string",
		},
		{
			Description: "The service name that appears in Apollo Studio",
			Name:        "serviceName",
			Type:        "string",
		},
		{
			Description: "Endpoint that service is available on when deployed",
			Name:        "introspectionURL",
			Type:        "string",
		},
		{
			Description: "",
			Name:        "options",
			Type:        "string",
		},
		{
			Description: "Whether to use the Apollo Rover CLI instead of the npm module",
			Name:        "useRover",
			Type:        "boolean",
		},
		{
			Description: "Whether to fail the pipeline if the linting fails",
			Name:        "failOnLintChecks",
			Type:        "boolean",
		},
		{
			Description: "Whether to publish the schema file",
			Name:        "publishSchema",
			Type:        "boolean",
		},
	}

	for i, v := range template.Parameters {
		expected := expectedParameters[i]
		if v.Name != expected.Name || v.Type != expected.Type || v.Description != expected.Description {
			t.Errorf("Parameters didn't match expectation.\nGot\n%v\nExpected\n%v", v, expectedParameters[i])
		}
	}
}

func TestRdender(t *testing.T) {
	template, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Got error %v", err)
	}

	var buf bytes.Buffer
	template.ToMarkdownTable(&buf)

	want := `|            NAME             |                                                                                  DESCRIPTION                                                                                  |  TYPE   | DEFAULT | REQUIRED |
|-----------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|---------|----------|
| localSchemaFile             | If grabbing the SDL from a local file. If the schemaIntrospectionUrl parameter is also set<br/>then this will be the file in the default working directory the schema will be | string  |         | true     |
| localSchemaIntrospectionUrl | If wanting to grab the SDL from a running server                                                                                                                              | string  |         | false    |
| commandToStartLocalServer   | If wanting to start a local server to introspect the schema                                                                                                                   | string  |         | false    |
| graph                       | The graph from Apollo Studio                                                                                                                                                  | string  |         | true     |
| apikey                      | Apollo studio API Key                                                                                                                                                         | string  |         | true     |
| variant                     | The variant of the graph in Apollo Studio                                                                                                                                     | string  |         | true     |
| serviceName                 | The service name that appears in Apollo Studio                                                                                                                                | string  |         | true     |
| introspectionURL            | Endpoint that service is available on when deployed                                                                                                                           | string  |         | false    |
| options                     |                                                                                                                                                                               | string  |         | false    |
| useRover                    | Whether to use the Apollo Rover CLI instead of the npm module                                                                                                                 | boolean | false   | false    |
| failOnLintChecks            | Whether to fail the pipeline if the linting fails                                                                                                                             | boolean | false   | false    |
| publishSchema               | Whether to publish the schema file                                                                                                                                            | boolean | false   | false    |
`
	got := buf.String()
	if got != want {
		t.Errorf("got:\n[%v]\nwant:\n[%v]\n", got, want)
	}
}

package openapi3

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/oasdiff/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefsJSON(t *testing.T) {
	loader := NewLoader()

	t.Log("Marshal *T to JSON")
	data, err := json.Marshal(spec())
	require.NoError(t, err)
	require.NotEmpty(t, data)

	t.Log("Unmarshal *T from JSON")
	docA := &T{}
	err = json.Unmarshal(specJSON, &docA)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	t.Log("Resolve refs in unmarshaled *T")
	err = loader.ResolveRefsIn(docA, nil)
	require.NoError(t, err)
	t.Log("Resolve refs in marshaled *T")
	docB, err := loader.LoadFromData(data)
	require.NoError(t, err)
	require.NotEmpty(t, docB)

	t.Log("Validate *T")
	err = docA.Validate(loader.Context)
	require.NoError(t, err)
	err = docB.Validate(loader.Context)
	require.NoError(t, err)

	t.Log("Ensure representations match")
	dataA, err := json.Marshal(docA)
	require.NoError(t, err)
	dataB, err := json.Marshal(docB)
	require.NoError(t, err)
	require.JSONEq(t, string(data), string(specJSON))
	require.JSONEq(t, string(data), string(dataA))
	require.JSONEq(t, string(data), string(dataB))
}

func TestRefsYAML(t *testing.T) {
	loader := NewLoader()

	t.Log("Marshal *T to YAML")
	data, err := yaml.Marshal(spec())
	require.NoError(t, err)
	require.NotEmpty(t, data)

	t.Log("Unmarshal *T from YAML")
	docA := &T{}
	err = yaml.Unmarshal(specYAML, &docA)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	t.Log("Resolve refs in unmarshaled *T")
	err = loader.ResolveRefsIn(docA, nil)
	require.NoError(t, err)
	t.Log("Resolve refs in marshaled *T")
	docB, err := loader.LoadFromData(data)
	require.NoError(t, err)
	require.NotEmpty(t, docB)

	t.Log("Validate *T")
	err = docA.Validate(loader.Context)
	require.NoError(t, err)
	err = docB.Validate(loader.Context)
	require.NoError(t, err)

	t.Log("Ensure representations match")
	dataA, err := yaml.Marshal(docA)
	require.NoError(t, err)
	dataB, err := yaml.Marshal(docB)
	require.NoError(t, err)
	require.YAMLEq(t, string(data), string(specYAML))
	require.YAMLEq(t, string(data), string(dataA))
	require.YAMLEq(t, string(data), string(dataB))
}

var specYAML = []byte(`
openapi: '3.0'
info:
  title: MyAPI
  version: '0.1'
paths:
  "/hello":
    parameters:
    - "$ref": "#/components/parameters/someParameter"
    post:
      parameters:
      - "$ref": "#/components/parameters/someParameter"
      requestBody:
        "$ref": "#/components/requestBodies/someRequestBody"
      responses:
        '200':
          "$ref": "#/components/responses/someResponse"
components:
  parameters:
    someParameter:
      description: Some parameter
      name: example
      in: query
      schema:
        "$ref": "#/components/schemas/someSchema"
  requestBodies:
    someRequestBody:
      description: Some request body
      content: {}
  responses:
    someResponse:
      description: Some response
  schemas:
    someSchema:
      description: Some schema
  headers:
    otherHeader:
      schema: {type: string}
    someHeader:
      "$ref": "#/components/headers/otherHeader"
  examples:
    otherExample:
      value:
        name: Some example
    someExample:
      "$ref": "#/components/examples/otherExample"
  securitySchemes:
    otherSecurityScheme:
      description: Some security scheme
      type: apiKey
      in: query
      name: token
    someSecurityScheme:
      "$ref": "#/components/securitySchemes/otherSecurityScheme"
`)

var specJSON = []byte(`
{
  "openapi": "3.0",
  "info": {
    "title": "MyAPI",
    "version": "0.1"
  },
  "paths": {
    "/hello": {
      "parameters": [
        {
          "$ref": "#/components/parameters/someParameter"
        }
      ],
      "post": {
        "parameters": [
          {
            "$ref": "#/components/parameters/someParameter"
          }
        ],
        "requestBody": {
          "$ref": "#/components/requestBodies/someRequestBody"
        },
        "responses": {
          "200": {
            "$ref": "#/components/responses/someResponse"
          }
        }
      }
    }
  },
  "components": {
    "parameters": {
      "someParameter": {
        "description": "Some parameter",
        "name": "example",
        "in": "query",
        "schema": {
          "$ref": "#/components/schemas/someSchema"
        }
      }
    },
    "requestBodies": {
      "someRequestBody": {
        "description": "Some request body",
        "content": {}
      }
    },
    "responses": {
      "someResponse": {
        "description": "Some response"
      }
    },
    "schemas": {
      "someSchema": {
        "description": "Some schema"
      }
    },
    "headers": {
      "otherHeader": {
        "schema": {
          "type": "string"
      	}
      },
      "someHeader": {
        "$ref": "#/components/headers/otherHeader"
      }
    },
    "examples": {
      "otherExample": {
        "value": {
          "name": "Some example"
        }
      },
      "someExample": {
        "$ref": "#/components/examples/otherExample"
      }
    },
    "securitySchemes": {
      "otherSecurityScheme": {
        "description": "Some security scheme",
        "type": "apiKey",
        "in": "query",
        "name": "token"
      },
      "someSecurityScheme": {
        "$ref": "#/components/securitySchemes/otherSecurityScheme"
      }
    }
  }
}
`)

func spec() *T {
	parameter := &Parameter{
		Description: "Some parameter",
		Name:        "example",
		In:          "query",
		Schema: &SchemaRef{
			Ref: "#/components/schemas/someSchema",
		},
	}
	requestBody := &RequestBody{
		Description: "Some request body",
		Content:     NewContent(),
	}
	responseDescription := "Some response"
	response := &Response{
		Description: &responseDescription,
	}
	schema := &Schema{
		Description: "Some schema",
	}
	example := map[string]string{"name": "Some example"}
	return &T{
		OpenAPI: "3.0",
		Info: &Info{
			Title:   "MyAPI",
			Version: "0.1",
		},
		Paths: NewPaths(
			WithPath("/hello", &PathItem{
				Post: &Operation{
					Parameters: Parameters{
						{
							Ref:   "#/components/parameters/someParameter",
							Value: parameter,
						},
					},
					RequestBody: &RequestBodyRef{
						Ref:   "#/components/requestBodies/someRequestBody",
						Value: requestBody,
					},
					Responses: NewResponses(
						WithStatus(200, &ResponseRef{
							Ref:   "#/components/responses/someResponse",
							Value: response,
						}),
					),
				},
				Parameters: Parameters{
					{
						Ref:   "#/components/parameters/someParameter",
						Value: parameter,
					},
				},
			}),
		),
		Components: &Components{
			Parameters: ParametersMap{
				"someParameter": {Value: parameter},
			},
			RequestBodies: RequestBodies{
				"someRequestBody": {Value: requestBody},
			},
			Responses: ResponseBodies{
				"someResponse": {Value: response},
			},
			Schemas: Schemas{
				"someSchema": {Value: schema},
			},
			Headers: Headers{
				"someHeader":  {Ref: "#/components/headers/otherHeader"},
				"otherHeader": {Value: &Header{Parameter{Schema: &SchemaRef{Value: NewStringSchema()}}}},
			},
			Examples: Examples{
				"someExample":  {Ref: "#/components/examples/otherExample"},
				"otherExample": {Value: NewExample(example)},
			},
			SecuritySchemes: SecuritySchemes{
				"someSecurityScheme": {Ref: "#/components/securitySchemes/otherSecurityScheme"},
				"otherSecurityScheme": {
					Value: &SecurityScheme{
						Description: "Some security scheme",
						Type:        "apiKey",
						In:          "query",
						Name:        "token",
					},
				},
			},
		},
	}
}

func TestValidation(t *testing.T) {
	version := `
openapi: 3.0.2
`
	info := `
info:
  title: "Hello World REST APIs"
  version: "1.0"
`
	paths := `
paths:
  "/api/v2/greetings.json":
    get:
      operationId: listGreetings
      responses:
        200:
          description: "List different greetings"
  "/api/v2/greetings/{id}.json":
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
          example: "greeting"
    get:
      operationId: showGreeting
      responses:
        200:
          description: "Get a single greeting object"
`
	externalDocs := `
externalDocs:
  url: https://root-ext-docs.com
`
	tags := `
tags:
  - name: "pet"
    externalDocs:
      url: https://tags-ext-docs.com
`
	spec := version + info + paths + externalDocs + tags + `
components:
  schemas:
    GreetingObject:
      properties:
        id:
          type: string
        type:
          type: string
          default: "greeting"
        attributes:
          properties:
            description:
              type: string
`

	tests := []struct {
		name        string
		spec        string
		expectedErr string
	}{
		{
			name: "no errors",
			spec: spec,
		},
		{
			name:        "version is missing",
			spec:        strings.Replace(spec, version, "", 1),
			expectedErr: "value of openapi must be a non-empty string",
		},
		{
			name:        "version is empty string",
			spec:        strings.Replace(spec, version, "openapi: ''", 1),
			expectedErr: "value of openapi must be a non-empty string",
		},
		{
			name:        "info section is missing",
			spec:        strings.Replace(spec, info, ``, 1),
			expectedErr: "invalid info: must be an object",
		},
		{
			name:        "paths section is missing",
			spec:        strings.Replace(spec, paths, ``, 1),
			expectedErr: "invalid paths: must be an object",
		},
		{
			name: "externalDocs section is invalid",
			spec: strings.Replace(spec, externalDocs,
				strings.ReplaceAll(externalDocs, "url: https://root-ext-docs.com", "url: ''"), 1),
			expectedErr: "invalid external docs: url is required",
		},
		{
			name: "tags section is invalid",
			spec: strings.Replace(spec, tags,
				strings.ReplaceAll(tags, "url: https://tags-ext-docs.com", "url: ''"), 1),
			expectedErr: "invalid tags: invalid external docs: url is required",
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			doc := &T{}
			err := yaml.Unmarshal([]byte(tt.spec), &doc)
			require.NoError(t, err)

			err = doc.Validate(context.Background())
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddRemoveServer(t *testing.T) {
	testServerLines := []*Server{{URL: "test0.com"}, {URL: "test1.com"}, {URL: "test3.com"}}

	doc3 := &T{
		OpenAPI:    "3.0.3",
		Components: &Components{},
	}

	assert.Empty(t, doc3.Servers)

	doc3.AddServer(&Server{URL: "testserver1.com"})

	assert.NotEmpty(t, doc3.Servers)
	assert.Len(t, doc3.Servers, 1)

	doc3.Servers = Servers{}

	assert.Empty(t, doc3.Servers)

	doc3.AddServers(testServerLines[0], testServerLines[1], testServerLines[2])

	assert.NotEmpty(t, doc3.Servers)
	assert.Len(t, doc3.Servers, 3)

	doc3.Servers = Servers{}

	doc3.AddServers(testServerLines...)

	assert.NotEmpty(t, doc3.Servers)
	assert.Len(t, doc3.Servers, 3)

	doc3.Servers = Servers{}
}

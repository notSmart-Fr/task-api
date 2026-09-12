package docs

import (
	"encoding/json"
	"testing"
)

func TestNormalizeOpenAPIDocument(t *testing.T) {
	input := []byte(`{
		"components": {
			"securitySchemes": {
				"Bearer": {
					"type": "apiKey"
				}
			}
		},
		"paths": {
			"/tasks": {
				"post": {
					"requestBody": {
						"content": {
							"application/json": {
								"schema": {
									"oneOf": [
										{"type": "object"},
										{"$ref": "#/components/schemas/tasks.CreateTaskRequest"}
									]
								}
							}
						}
					}
				}
			}
		}
	}`)

	output, err := normalizeOpenAPIDocument(input)
	if err != nil {
		t.Fatalf("normalizeOpenAPIDocument() error = %v", err)
	}

	var document map[string]any
	if err := json.Unmarshal(output, &document); err != nil {
		t.Fatalf("unmarshal normalized document: %v", err)
	}

	paths := document["paths"].(map[string]any)
	operation := paths["/tasks"].(map[string]any)["post"].(map[string]any)
	requestBody := operation["requestBody"].(map[string]any)
	content := requestBody["content"].(map[string]any)
	schema := content["application/json"].(map[string]any)["schema"].(map[string]any)

	if actual := schema["$ref"]; actual != "#/components/schemas/tasks.CreateTaskRequest" {
		t.Fatalf("request schema $ref = %v, want CreateTaskRequest reference", actual)
	}
	if _, exists := schema["oneOf"]; exists {
		t.Fatal("request schema still contains oneOf")
	}

	components := document["components"].(map[string]any)
	securitySchemes := components["securitySchemes"].(map[string]any)
	bearerScheme := securitySchemes["Bearer"].(map[string]any)
	if actual := bearerScheme["type"]; actual != "http" {
		t.Fatalf("Bearer security scheme type = %v, want http", actual)
	}
	if actual := bearerScheme["scheme"]; actual != "bearer" {
		t.Fatalf("Bearer security scheme = %v, want bearer", actual)
	}
}

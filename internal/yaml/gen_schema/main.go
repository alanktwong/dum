//go:build gen_schema

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "Dum Installer Config",
		"type":    "object",
		"properties": map[string]interface{}{
			"playbook": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"book":        map[string]interface{}{"type": "string"},
					"description": map[string]interface{}{"type": "string"},
					"enabled":     map[string]interface{}{"type": "boolean"},
					"sudo":        map[string]interface{}{"type": "boolean"},
					"jetbrains": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type":                 "object",
							"additionalProperties": map[string]interface{}{"type": "string"},
						},
					},
					"plays": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"play":        map[string]interface{}{"type": "string"},
								"description": map[string]interface{}{"type": "string"},
								"enabled":     map[string]interface{}{"type": "boolean"},
								"tasks": map[string]interface{}{
									"type": "array",
									"items": map[string]interface{}{
										"type": "object",
										"properties": map[string]interface{}{
											"type": map[string]interface{}{
												"type": "string",
												"enum": []string{
													"bash",
													"dir",
													"link",
													"git",
													"brew",
													"cask",
													"cellar",
													"function",
													"mas",
													"vscode",
													"jetbrains",
												},
											},
											"task":        map[string]interface{}{"type": "string"},
											"description": map[string]interface{}{"type": "string"},
											"enabled":     map[string]interface{}{"type": "boolean"},
											"command":     map[string]interface{}{"type": "string"},
											"script":      map[string]interface{}{"type": "string"},
											"root":        map[string]interface{}{"type": "string"},
											"target":      map[string]interface{}{"type": "string"},
											"name":        map[string]interface{}{"type": "string"},
											"tap":         map[string]interface{}{"type": "string"},
											"apps": map[string]interface{}{
												"type":  "array",
												"items": map[string]interface{}{"type": "string"},
											},
										},
										"required": []string{"type", "task"},
									},
								},
							},
							"required": []string{"play", "tasks"},
						},
					},
				},
				"required": []string{"book", "plays"},
			},
		},
		"required": []string{"playbook"},
	}

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling schema: %v\n", err)
		os.Exit(1)
	}

	outputPath := "../../cfg/installer.schema.json"
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	err = os.WriteFile(outputPath, data, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated JSON schema at %s\n", outputPath)
}

// Command gen-schema generates the JSON schema for dum.yml configuration.
//
// It runs as part of 'make generate' via the go:generate directive below and
// writes the schema that is embedded into the dum binary by cmd/dum.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

//go:generate go run . ../../cmd/dum/installer.schema.json

const (
	tObject  = "object"
	tString  = "string"
	tBoolean = "boolean"
	tArray   = "array"

	kType        = "type"
	kEnabled     = "enabled"
	kProperties  = "properties"
	kDescription = "description"
	kItems       = "items"
	kRequired    = "required"
)

func main() {
	schema := map[string]any{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "Dum Installer Config",
		kType:     tObject,
		kProperties: map[string]any{
			"playbook": map[string]any{
				kType: tObject,
				kProperties: map[string]any{
					"book":       map[string]any{kType: tString},
					kDescription: map[string]any{kType: tString},
					kEnabled:     map[string]any{kType: tBoolean},
					"sudo":       map[string]any{kType: tBoolean},
					"jetbrains": map[string]any{
						kType: tArray,
						kItems: map[string]any{
							kType:                  tObject,
							"additionalProperties": map[string]any{kType: tString},
						},
					},
					"plays": map[string]any{
						kType: tArray,
						kItems: map[string]any{
							kType: tObject,
							kProperties: map[string]any{
								"play":       map[string]any{kType: tString},
								kDescription: map[string]any{kType: tString},
								kEnabled:     map[string]any{kType: tBoolean},
								"tasks": map[string]any{
									kType: tArray,
									kItems: map[string]any{
										kType: tObject,
										kProperties: map[string]any{
											kType: map[string]any{
												kType: tString,
												"enum": []string{
													"bash",
													"dir",
													"link",
													"git",
													"brew",
													"cask",
													"cellar",
													"mas",
													"vscode",
													"jetbrains",
												},
											},
											"task":       map[string]any{kType: tString},
											kDescription: map[string]any{kType: tString},
											kEnabled:     map[string]any{kType: tBoolean},
											"command":    map[string]any{kType: tString},
											"script":     map[string]any{kType: tString},
											"root":       map[string]any{kType: tString},
											"target":     map[string]any{kType: tString},
											"name":       map[string]any{kType: tString},
											"tap":        map[string]any{kType: tString},
											"apps": map[string]any{
												kType:  tArray,
												kItems: map[string]any{kType: tString},
											},
										},
										kRequired: []string{"type", "task"},
									},
								},
							},
							kRequired: []string{"play", "tasks"},
						},
					},
				},
				kRequired: []string{"book", "plays"},
			},
		},
		kRequired: []string{"playbook"},
	}

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling schema: %v\n", err)
		os.Exit(1)
	}

	outputPath := "../../cmd/dum/installer.schema.json"
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	data = append(data, '\n')
	// #nosec G703 -- outputPath is a build-time flag supplied by go:generate.
	err = os.WriteFile(outputPath, data, 0o600)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing schema: %v\n", err)
		os.Exit(1)
	}
}

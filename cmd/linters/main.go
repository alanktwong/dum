// Package main provides a custom linter for detecting single-character import aliases.
//
// This is a standalone Go analysis tool (golang.org/x/tools/go/analysis) with its own
// go.mod because:
//   - It's designed to run independently: go run ./cmd/linters
//   - It can be integrated with golangci-lint v2 as a custom linter plugin
//   - It has minimal dependencies (only golang.org/x/tools) separate from the main project
//
// When golangci-lint v2 plugin support matures, this could potentially be consolidated
// into the main module if desired.

package main

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "linters",
	Doc:  "Custom linters for dum project",
	Run:  run,
}

func main() {
	singlechecker.Main(Analyzer)
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				for _, spec := range genDecl.Specs {
					importSpec := spec.(*ast.ImportSpec)
					if importSpec.Name != nil {
						name := importSpec.Name.Name
						if len(name) == 1 && name != "_" {
							pass.Reportf(importSpec.Pos(), "import alias %q is a single character", name)
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

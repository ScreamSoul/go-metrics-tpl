package exitinmain

import (
	"go/ast"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var ExitInMainCheckAnalyzer = &analysis.Analyzer{
	Name: "exitinmain",
	Doc:  "check that excludes the use of os.Exit() in main function in main packege",
	Run:  run,
}

func isMianFunction(node ast.Node) bool {
	funcDecl, ok := node.(*ast.FuncDecl)
	return ok && funcDecl.Name.Name == "main" && len(funcDecl.Type.Params.List) == 0 && funcDecl.Type.Results == nil
}

// Helper function to determine if we're inside the main function of the main package
func isContainsMainOs(file *ast.File) bool {
	isOsImport := func() bool {
		for _, imp := range file.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue // Skip invalid imports
			}
			if strings.TrimPrefix(path, "\"") == "os" {
				return true
			}
		}
		return false
	}

	isContainsMainFunc := func() bool {
		for _, decl := range file.Decls {
			if isMianFunction(decl) {
				return true
			}
		}
		return false
	}

	isPackegeMain := func() bool {
		return file.Name.Name == "main"
	}

	return isPackegeMain() && isOsImport() && isContainsMainFunc()
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if !isContainsMainOs(file) {
			continue
		}

		for _, decl := range file.Decls {
			if !isMianFunction(decl) {
				continue
			}

			// Traversal for `main() {}` function
			ast.Inspect(decl, func(node ast.Node) bool {
				callExpr, ok := node.(*ast.CallExpr)
				if !ok {
					return true
				}

				selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
				if !ok || selExpr.Sel.Name != "Exit" {
					return true
				}

				ident, ok := selExpr.X.(*ast.Ident)
				if !ok || ident.Name != "os" {
					return true
				}

				pass.Report(analysis.Diagnostic{
					Pos:     callExpr.Pos(),
					Message: "use of os.Exit() in main function is discouraged",
					// SuggestedFixes: []analysis.SuggestedFix{{Message: "Consider returning an error instead"}},
				})

				return true
			})

		}

	}
	return nil, nil
}

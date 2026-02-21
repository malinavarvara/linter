package analazer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "linter",
	Doc:  "checks log messages for style and sensitive data",
	Run:  run,
}

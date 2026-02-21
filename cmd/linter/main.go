package main

import (
	"linter/analyzer"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() {
	unitchecker.Main(analyzer.Analyzer)
}

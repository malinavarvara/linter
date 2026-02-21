package main

import (
	"golang.org/x/tools/go/analysis"
	"linter/analyzer"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		analyzer.Analyzer,
	}, nil
}

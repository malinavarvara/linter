package main

import (
	"fmt"
	"linter/analyzer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	fmt.Println("Запуск линтера...")
	singlechecker.Main(analyzer.Analyzer)
}

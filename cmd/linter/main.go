package main

import (
	"fmt"
	"linter/analazer"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	fmt.Println("Запуск линтера...")
	singlechecker.Main(analazer.Analyzer)
}

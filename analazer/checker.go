package analazer

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		fmt.Printf("Файл: %s\n", pass.Fset.Position(file.Pos()).Filename)
		ast.Inspect(file, func(n ast.Node) bool {
			node, ok := n.(*ast.CallExpr)
			// Проверяем, что это вызов функции и что она является функцией логирования из slog или zap
			if !ok || !isLogCall(node, pass.TypesInfo) || len(node.Args) == 0 {
				return true
			}
			firstArg := node.Args[0]
			// проверка на чувствительные данные
			for _, s := range collectStringLiterals(firstArg, pass) {
				checkSensitiveData(s, node, pass)
			}
			// проверка на все остальное
			if fullMsg := extractFullString(firstArg, pass); fullMsg != "" {
				checkMessageBasic(fullMsg, node, pass)
			}

			return true
		})
	}
	return nil, nil
}

// является ли вызов функцией логирования из slog или zap
func isLogCall(node *ast.CallExpr, info *types.Info) bool {
	var funObj types.Object
	switch fun := node.Fun.(type) {
	case *ast.Ident:
		funObj = info.ObjectOf(fun)
	case *ast.SelectorExpr:
		funObj = info.ObjectOf(fun.Sel)
	default:
		return false
	}
	if funObj == nil {
		return false
	}

	// Проверка функций (например, slog.Info, zap.Info)
	if pkg := funObj.Pkg(); pkg != nil {
		switch pkg.Name() {
		case "slog":
			switch funObj.Name() {
			case "Debug", "Info", "Warn", "Error", "Log":
				return true
			}
		case "zap":
			switch funObj.Name() {
			case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal":
				return true
			}
		}
	}

	// Проверка методов (например, logger.Info)
	if sig, ok := funObj.Type().(*types.Signature); ok && sig.Recv() != nil {
		recvType := sig.Recv().Type()
		var named *types.Named
		switch t := recvType.(type) {
		case *types.Named:
			named = t
		case *types.Pointer:
			if n, ok := t.Elem().(*types.Named); ok {
				named = n
			}
		}
		if named != nil {
			if pkg := named.Obj().Pkg(); pkg != nil {
				switch pkg.Name() {
				case "slog", "zap":
					switch funObj.Name() {
					case "Debug", "Info", "Warn", "Error", "Log", "DPanic", "Panic", "Fatal":
						return true
					}
				}
			}
		}
	}
	return false
}

func collectStringLiterals(expr ast.Expr, pass *analysis.Pass) []string {
	var results []string
	ast.Inspect(expr, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.BasicLit:
			if node.Kind == token.STRING {
				s, err := strconv.Unquote(node.Value)
				if err == nil {
					results = append(results, s)
				} else {
					results = append(results, node.Value[1:len(node.Value)-1])
				}
			}
		case *ast.Ident:
			if obj, ok := pass.TypesInfo.ObjectOf(node).(*types.Const); ok {
				if basic, ok := obj.Type().(*types.Basic); ok && basic.Info()&types.IsString != 0 {
					results = append(results, constant.StringVal(obj.Val()))
				}
			}
		}
		return true
	})
	return results
}

func extractFullString(expr ast.Expr, pass *analysis.Pass) string {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			s, err := strconv.Unquote(e.Value)
			if err == nil {
				return s
			}
			return e.Value[1 : len(e.Value)-1]
		}
	case *ast.Ident:
		if obj, ok := pass.TypesInfo.ObjectOf(e).(*types.Const); ok {
			if basic, ok := obj.Type().(*types.Basic); ok && basic.Info()&types.IsString != 0 {
				return constant.StringVal(obj.Val())
			}
		}
	}
	return ""
}

func checkMessageBasic(msg string, node *ast.CallExpr, pass *analysis.Pass) {
	if msg == "" {
		return
	}

	// Правило 1: сообщение должно начинаться со строчной латинской буквы
	first := msg[0]
	if !(first >= 'a' && first <= 'z') {
		pass.Reportf(node.Pos(), "log message should start with a lowercase Latin letter")
	}

	// Правила 2 и 3: только латинские буквы, цифры и пробелы
	extraAllowed := map[rune]bool{'.': true, ':': true, '=': true, '-': true, '_': true} // разрешаем некоторые спецсимволы, часто используемые в логах
	for _, r := range msg {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || unicode.IsDigit(r) || unicode.IsSpace(r) || extraAllowed[r]) {
			pass.Reportf(node.Pos(), "log message contains invalid characters (only Latin letters, digits, spaces, and allowed punctuation are permitted)")
			break
		}
	}
}

// проверяет строку на наличие паттернов чувствительных данных
func checkSensitiveData(s string, node ast.Node, pass *analysis.Pass) {
	if s == "" {
		return
	}
	sensitivePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bpassword\s*[:=]`),
		regexp.MustCompile(`(?i)\btoken\s*[:=]`),
		regexp.MustCompile(`(?i)\bapi_?key\s*[:=]`),
		regexp.MustCompile(`(?i)\bkey\s*[:=]`),
		regexp.MustCompile(`(?i)\bsecret\s*[:=]`),
		regexp.MustCompile(`(?i)\bauth\s*[:=]`),
		regexp.MustCompile(`(?i)\bpass\s*[:=]`),
	}
	for _, re := range sensitivePatterns {
		if re.MatchString(s) {
			pass.Reportf(node.Pos(), "log message contains potentially sensitive data pattern")
			return
		}
	}
}

package foo

import (
	"fmt"
	"log/slog"
	"os"

	"go.uber.org/zap"
)

// Константы для проверки
const (
	startMsg       = "server started"
	invalidMsg     = "Сервер запущен"
	sensitiveConst = "password: "
)

// Переменные для динамических частей
var (
	token      = "abc123"
	pwd        = "secret"
	authHeader = "Bearer xyz"
	apiKey     = "12345"
	auth       = "basic"
	value      = "val"
)

func main() {
	// Инициализация логгеров
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	// ----- ПРАВИЛЬНЫЕ СООБЩЕНИЯ (не должны вызывать ошибок) -----
	slog.Info("starting server")
	slog.Debug("connection established")
	slog.Warn("disk space low")
	slog.Error("failed to connect")
	slog.Warn("warning: something went wrong...")

	logger.Info("request processed")
	logger.Debug("cache hit")

	zap.L().Info("everything is fine")
	zapLogger.Info("all systems go")

	// Правильные сообщения, где нет чувствительных паттернов
	slog.Info("user authenticated successfully")
	slog.Info("token validated")
	slog.Info("password updated")

	// Правильное сообщение из константы
	slog.Info(startMsg)

	// ----- НЕПРАВИЛЬНЫЕ: первая буква заглавная -----
	slog.Info("Starting server")      // want "lowercase"
	slog.Warn("Disk space low")       // want "lowercase"
	logger.Error("Failed to connect") // want "lowercase"
	zapLogger.Info("All systems go")  // want "lowercase"

	// ----- НЕПРАВИЛЬНЫЕ: нелатинские символы (кириллица, эмодзи) -----
	slog.Info("сервер запущен")      // want "lowercase Latin letter" "invalid characters"
	slog.Debug("ошибка подключения") // want "lowercase Latin letter" "invalid characters"
	logger.Warn("предупреждение")    // want "lowercase Latin letter" "invalid characters"
	zapLogger.Info("всё хорошо")     // want "lowercase Latin letter" "invalid characters"
	// А здесь первая буква правильная, поэтому линтер доходит до проверки символов
	slog.Info("starting сервер")              // want "invalid characters"
	slog.Info("task completed ✅")             // want "invalid characters"
	zapLogger.Debug("response time: 100ms 🚀") // want "invalid characters"

	// ----- НЕПРАВИЛЬНЫЕ: спецсимволы, не входящие в extraAllowed -----
	slog.Info("server started!")        // want "invalid characters"
	slog.Error("connection failed!!")   // want "invalid characters"
	slog.Info("are you sure?")          // want "invalid characters"
	logger.Info("user input: <script>") // want "invalid characters"

	// ----- НЕПРАВИЛЬНЫЕ: чувствительные данные -----
	// Прямые литералы с паттерном
	slog.Info("user password: 123")         // want "sensitive"
	slog.Debug("api_key=abcd")              // want "sensitive"
	logger.Info("token: xyz")               // want "sensitive"
	zapLogger.Info("authorization: Bearer") // с маленькой буквы, чтобы пройти первую проверку -> want "sensitive"
	slog.Info("password = secret")          // want "sensitive"
	slog.Info("key : value")                // want "sensitive"
	slog.Info("secret: data")               // want "sensitive"
	slog.Info("auth = basic")               // want "sensitive"

	// Конкатенация строк
	slog.Info("user password: " + "secret") // want "sensitive"
	slog.Debug("api_key=" + "12345")        // want "sensitive"
	logger.Info("token: " + token)          // want "sensitive"
	slog.Info("password = " + pwd)          // want "sensitive"
	slog.Info("secret: " + "sec")           // want "sensitive"
	slog.Info("auth = " + auth)             // want "sensitive"

	// Чувствительные данные внутри константы
	slog.Info(sensitiveConst + "123") // want "sensitive"

	// Смешанные случаи (fmt.Sprintf и другие выражения)
	slog.Info(fmt.Sprintf("token: %s", token)) // want "sensitive"
	slog.Info("prefix: " + "token=" + value)   // want "sensitive"
}

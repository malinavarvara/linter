package main

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
	slog.Info("Starting server")      // должно ругаться: lowercase Latin letter
	slog.Warn("Disk space low")       // должно ругаться
	logger.Error("Failed to connect") // должно ругаться
	zapLogger.Info("All systems go")  // должно ругаться

	// ----- НЕПРАВИЛЬНЫЕ: нелатинские символы (кириллица, эмодзи) -----
	slog.Info("сервер запущен")               // должно ругаться: invalid characters
	slog.Debug("ошибка подключения")          // должно ругаться
	logger.Warn("предупреждение")             // должно ругаться
	zapLogger.Info("всё хорошо")              // должно ругаться
	slog.Info("starting сервер")              // должно ругаться (смесь)
	slog.Info("task completed ✅")             // должно ругаться (эмодзи)
	zapLogger.Debug("response time: 100ms 🚀") // должно ругаться

	// ----- НЕПРАВИЛЬНЫЕ: спецсимволы, не входящие в extraAllowed -----
	slog.Info("server started!")        // '!' не разрешён
	slog.Error("connection failed!!")   // '!' не разрешён
	slog.Info("are you sure?")          // '?' не разрешён
	logger.Info("user input: <script>") // '<' и '>' не разрешены

	// ----- НЕПРАВИЛЬНЫЕ: чувствительные данные -----
	// Прямые литералы с паттерном
	slog.Info("user password: 123")         // должно ругаться
	slog.Debug("api_key=abcd")              // должно ругаться
	logger.Info("token: xyz")               // должно ругаться
	zapLogger.Info("Authorization: Bearer") // должно ругаться (auth)
	slog.Info("password = secret")          // должно ругаться
	slog.Info("key : value")                // должно ругаться (key с двоеточием)
	slog.Info("secret: data")               // должно ругаться
	slog.Info("auth = basic")               // должно ругаться

	// Конкатенация строк
	slog.Info("user password: " + "secret") // должно ругаться
	slog.Debug("api_key=" + "12345")        // должно ругаться
	logger.Info("token: " + token)          // должно ругаться
	slog.Info("password = " + pwd)          // должно ругаться
	slog.Info("secret: " + "sec")           // должно ругаться
	slog.Info("auth = " + auth)             // должно ругаться

	// Чувствительные данные внутри константы
	slog.Info(sensitiveConst + "123") // (sensitiveConst содержит "password: ")

	// Смешанные случаи (fmt.Sprintf и другие выражения)
	slog.Info(fmt.Sprintf("token: %s", token)) // (литерал "token: " внутри Sprintf)
	slog.Info("prefix: " + "token=" + value)   //(литерал "token=")
}

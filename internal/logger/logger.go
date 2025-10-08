package logger

import (
	"io"
	"log/slog"
	"os"
)

var Log *slog.Logger

// Init inicializa o logger global.
// Modo padrão: texto no stdout.
// Modo json: se LOG_MODE=json.
// Modo teste: silencioso.
func Init() {
	mode := os.Getenv("LOG_MODE")


	var handler slog.Handler

	switch mode{
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, nil)
	case "silent":
		handler = slog.NewTextHandler(io.Discard, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	Log = slog.New(handler)
}

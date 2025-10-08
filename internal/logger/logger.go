package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger

// Init inicializa o logger global.
// Modo padrão: texto no stdout.
// Modo json: se LOG_MODE=json.
// Modo teste: silencioso.
func Init() {
	mode := os.Getenv("LOG_MODE")

	isTest := strings.HasSuffix(os.Args[0], ".test")

	var handler slog.Handler

	switch {
	case isTest:
		handler = slog.NewTextHandler(io.Discard, nil)
	case mode == "json":
		handler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	Log = slog.New(handler)
}

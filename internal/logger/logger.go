package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init() {
	mode := os.Getenv("LOG_MODE")
	if mode == "json" {
		Log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		Log = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

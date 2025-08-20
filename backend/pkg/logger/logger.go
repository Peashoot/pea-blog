package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func New(environment string) *Logger {
	var handler slog.Handler
	
	if environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	
	logger := slog.New(handler)
	return &Logger{Logger: logger}
}

func (l *Logger) Fatal(msg string, err error) {
	l.Error(msg, "error", err)
	os.Exit(1)
}
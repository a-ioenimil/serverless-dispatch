package logger

import (
	"log/slog"
	"os"
)

// InitLogger configures a JSON logger suitable for AWS Lambda
func InitLogger() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

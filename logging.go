package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

// InitLogger sets up the application logger with JSON output
// Logs are written to the specified file path (or stdout if empty)
func InitLogger(logPath string, level slog.Level) error {
	var handler slog.Handler

	if logPath != "" {
		// Log to file
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Log to stdout
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.LogAttrs(context.Background(), slog.LevelDebug, "Gapistotle logger initialized",
		slog.String("log_path", logPath),
		slog.String("log_level", level.String()),
	)

	return nil
}

// LogDebug logs a debug message with optional key-value attributes
func LogDebug(msg string, attrs ...any) {
	slog.Debug(msg, attrs...)
}

// LogInfo logs an info message with optional key-value attributes
func LogInfo(msg string, attrs ...any) {
	slog.Info(msg, attrs...)
}

// LogWarn logs a warning message with optional key-value attributes
func LogWarn(msg string, attrs ...any) {
	slog.Warn(msg, attrs...)
}

// LogError logs an error message with optional key-value attributes
func LogError(msg string, attrs ...any) {
	slog.Error(msg, attrs...)
}

// ParseLogLevel converts a string log level to slog.Level
func ParseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug // Default to debug
	}
}

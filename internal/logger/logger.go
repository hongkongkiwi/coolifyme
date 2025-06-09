package logger

import (
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	// Create a default logger
	defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// SetLevel sets the logging level
func SetLevel(level slog.Level) {
	defaultLogger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
}

// SetJSONOutput enables JSON formatted logging
func SetJSONOutput() {
	defaultLogger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// With returns a logger with additional context
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

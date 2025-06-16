// Package logger provides logging utilities for the coolifyme CLI tool.
package logger

import (
	"log/slog"
	"os"
)

var (
	defaultLogger *slog.Logger
	colorEnabled  bool
)

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

// SetColorOutput enables or disables color output
func SetColorOutput(enabled bool) {
	colorEnabled = enabled
}

// ColorEnabled returns whether color output is enabled
func ColorEnabled() bool {
	return colorEnabled
}

// IsTerminal checks if output is going to a terminal
func IsTerminal() bool {
	stat, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
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

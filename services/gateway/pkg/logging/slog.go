package logging

import (
	"log/slog"
	"os"
)

type JSONLogger struct {
	logger slog.Logger
}

func NewLogger() Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func (l *JSONLogger) Debug(msg string, fields ...interface{}) {
	l.Debug(msg, fields...)
}

func (l *JSONLogger) Info(msg string, fields ...interface{}) {
	l.Info(msg, fields...)
}

func (l *JSONLogger) Warn(msg string, fields ...interface{}) {
	l.Warn(msg, fields...)
}

func (l *JSONLogger) Error(msg string, fields ...interface{}) {
	l.Error(msg, fields...)
}

package logger

// Logger is an interface for logging messages at different levels.
type Logger interface {
	// Debug logs a message at the debug level.
	Debug(msg string, fields ...interface{})

	// Info logs a message at the info level.
	Info(msg string, fields ...interface{})

	// Warn logs a message at the warning level.
	Warn(msg string, fields ...interface{})

	// Error logs a message at the error level.
	Error(msg string, fields ...interface{})
}

package logger

import (
	"io"
	"log/slog"
	"os"
)

// LogLevel represents the log level
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
}

// Config holds logger configuration
type Config struct {
	Level  LogLevel `json:"level"`
	Format string   `json:"format"` // "json" or "text"
	Output string   `json:"output"` // "stdout", "stderr", or file path
}

// New creates a new logger with the given configuration
func New(config Config) (*Logger, error) {
	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var writer io.Writer
	switch config.Output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// If it's not stdout/stderr, treat as file path
		if config.Output != "" {
			file, err := os.OpenFile(config.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
			writer = file
		} else {
			writer = os.Stdout
		}
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: level,
	}

	if config.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger := slog.New(handler)
	return &Logger{Logger: logger}, nil
}

// NewDefault creates a default logger
func NewDefault() *Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return &Logger{Logger: logger}
}

// WithContext returns a logger with the given context attributes
func (l *Logger) WithContext(attrs ...slog.Attr) *Logger {
	return &Logger{Logger: l.Logger.With(attrsToAny(attrs...)...)}
}

// WithError adds an error to the log context
func (l *Logger) WithError(err error) *Logger {
	return &Logger{Logger: l.Logger.With("error", err.Error())}
}

// WithRequestID adds a request ID to the log context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{Logger: l.Logger.With("request_id", requestID)}
}

// WithUserID adds a user ID to the log context
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{Logger: l.Logger.With("user_id", userID)}
}

// Helper function to convert slog.Attr to []any
func attrsToAny(attrs ...slog.Attr) []any {
	result := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		result = append(result, attr.Key, attr.Value)
	}
	return result
}

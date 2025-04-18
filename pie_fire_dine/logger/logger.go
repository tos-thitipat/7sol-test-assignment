package logger

import (
	"errors"
	"io"
	"log"
	"log/slog"
	"strings"
)

const JSONLogger = "json"

type ContextKey string

type ContextHandler struct {
	slog.Handler
}

type LogConfig struct {
	Level  string
	Format string
	Out    io.Writer
}

func ParseLogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelDebug, errors.New("unknown log level")
	}
}

func SetupLogger(config *LogConfig) {
	var loggerHandler slog.Handler
	level, err := ParseLogLevel(config.Level)
	if err != nil {
		log.Fatal(err.Error())
	}

	if config.Format == "json" {
		loggerHandler = &ContextHandler{slog.NewJSONHandler(config.Out, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})}
	} else {
		loggerHandler = &ContextHandler{slog.NewTextHandler(config.Out, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		})}
	}

	// Create logger
	logger := slog.New(loggerHandler)
	slog.SetDefault(logger)
}

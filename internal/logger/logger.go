package logger

import (
	"log/slog"
	"os"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
)

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func New(config *config.LoggerConfig) *slog.Logger {
	var logger *slog.Logger
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: *config.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	switch config.Level {
	case "trace":
		opts.Level = LevelTrace
	case "debug":
		opts.Level = slog.LevelDebug
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	case "fatal":
		opts.Level = LevelFatal
	}

	if *config.Structured {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	slog.SetDefault(logger)

	return logger
}

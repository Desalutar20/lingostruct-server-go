package logger

import (
	"log/slog"
	"os"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
)

func New(config *config.ApplicationConfig) *slog.Logger {
	var level slog.Level

	if err := level.UnmarshalText([]byte(config.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}
	var handler slog.Handler = slog.NewJSONHandler(os.Stderr, opts)

	if config.PrettyLog {
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(handler)
}

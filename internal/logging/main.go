package logging

import (
	"log/slog"
	"os"

	"github.com/panoptescloud/orca/internal/config"
)

const LogLevelNone = 10
const LogLevelNoneName = "none"

func NewSlogHandler(cfg *config.Config) (slog.Handler, error) {
	l := slog.Level(LogLevelNone)

	if cfg.Logging.Level != LogLevelNoneName {
		err := l.UnmarshalText([]byte(cfg.Logging.Level))

		if err != nil {
			return nil, err
		}
	}

	opts := &slog.HandlerOptions{
		Level: l,
	}

	var handler slog.Handler

	if cfg.Logging.Format == "text" {
		handler = slog.NewTextHandler(os.Stderr, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stderr, opts)
	}

	return handler, nil
}

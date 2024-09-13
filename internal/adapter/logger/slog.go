package logger

import (
	"golang-hexagon/internal/adapter/config"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// logger is the default logger used by the application
var logger *slog.Logger

// Set sets the logger configuration based on the environment
func Set(conf *config.App) {
	logger = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	)

	if conf.Env == config.EnvProduction {
		logRotate := &lumberjack.Logger{
			Filename:   "log/app.log",
			MaxSize:    20, // megabytes
			MaxBackups: 3,
			MaxAge:     1, // days
			Compress:   true,
		}

		logger = slog.New(
			slogmulti.Fanout(
				slog.NewJSONHandler(logRotate, nil),
				slog.NewTextHandler(os.Stderr, nil),
			),
		)
	}

	slog.SetDefault(logger)
}

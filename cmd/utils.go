package cmd

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

func checkErr(err error, msg string) {
	if err != nil {
		slog.Error(msg, "err", err)
		os.Exit(1)
	}
}

func initLogging() {
	level, validLevel := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}[strings.ToLower(logLevel)]
	if !validLevel {
		slog.Error("Invalid log level", "level", logLevel)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(
			os.Stderr,
			&tint.Options{
				TimeFormat: time.Kitchen,
				Level:      level,
			},
		),
	))
}

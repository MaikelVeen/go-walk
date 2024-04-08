package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		}),
	)

	logger.Info("Hello World!")
}

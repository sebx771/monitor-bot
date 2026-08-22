package logger

import (
	"log/slog"
	"os"
)

var handler= NewConsoleHandler(os.Stdout)
var baseLogger = slog.New(handler)

func GetBaseLogger() *slog.Logger {
	return baseLogger
}
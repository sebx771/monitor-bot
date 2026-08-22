package logger

import (
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(module string) *Logger {

    log:= GetBaseLogger()

	return &Logger{
		logger: log.With("module",module),
	}
}

func (l *Logger) Debug(message string, attrs ...any) {
    l.logger.Debug(message,attrs...)
}

func (l *Logger) Info(message string, attrs ...any){
	l.logger.Info(message,attrs...)
}
func (l *Logger) Error(message string, attrs ...any){
	l.logger.Error(message,attrs...)
}

func (l *Logger) Warn(message string , attrs ...any){
	l.logger.Warn(message,attrs...)
}

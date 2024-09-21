package mylogger

import (
	"context"
	"log/slog"
	"os"
)

type MyLogger struct {
	logger *slog.Logger
	logLvl *slog.LevelVar
}

func (ml *MyLogger) InitLogger(loggerFilePath string, level string) error {
	logFile, err := os.OpenFile(loggerFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	ml.logLvl = &slog.LevelVar{}
	if level == "debug" {
		ml.logLvl.Set(slog.LevelDebug)
	}

	opts := &slog.HandlerOptions{
		AddSource: level != "info",
		Level:     ml.logLvl,
	}
	jsonHandler := slog.NewJSONHandler(logFile, opts)
	slog.SetDefault(slog.New(jsonHandler))
	ml.logger = slog.Default()

	return nil
}

func (ml *MyLogger) setLogLevel(level slog.Level) {
	ml.logLvl.Set(level)
}

func (ml *MyLogger) WriteLog(msg string, level slog.Level, additionalAttrs *slog.Attr) {
	if additionalAttrs != nil {
		slog.LogAttrs(context.Background(), level, msg, *additionalAttrs)
	} else {
		slog.LogAttrs(context.Background(), level, msg)
	}
}

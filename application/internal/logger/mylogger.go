package mylogger

import (
	"context"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strings"
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
		// Он запоминает строку, где вызвалась функция WriteLog,
		// а строку, где она выполнилась (то есть в файле mylogger.go
		// То есть пока что данная опция бесполезна (я не могу узнать, в какой строке произошла ошибка)
		// AddSource: level != "info",
		Level: ml.logLvl,
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

func LogCallerInfo() *slog.Attr {
	pc, file, line, _ := runtime.Caller(2)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	info := slog.Group(
		"call place info",
		slog.String("package name", packageName),
		slog.String("file name", fileName),
		slog.String("function name", funcName),
		slog.Int("line", line),
	)

	return &info
}

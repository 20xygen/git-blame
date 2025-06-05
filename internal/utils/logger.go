package utils

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logDir     = "/var/log/blame"
	logFile    = "blame.log"
	maxSizeMB  = 10 // MB
	maxBackups = 5
	maxAgeDays = 30
)

type silentWriter struct{}

func (w silentWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var defaultWriter = silentWriter{}

func SetupLogger() *slog.Logger {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return slog.New(slog.NewTextHandler(defaultWriter, nil))
	}

	logPath := filepath.Join(logDir, logFile)
	logRotator := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAgeDays,
		Compress:   true,
	}

	if _, err := logRotator.Write([]byte("")); err != nil {
		return slog.New(slog.NewTextHandler(defaultWriter, nil))
	}

	return slog.New(slog.NewJSONHandler(logRotator, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

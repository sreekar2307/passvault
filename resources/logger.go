package resources

import (
	"context"
	"log/slog"
	"os"
	"passVault/dtos"
)

var (
	fileHandler slog.Handler
	logger      *slog.Logger
	level       slog.Level
)

func initLogger() error {
	logFile := config.GetString(dtos.ConfigKeys.Logger.File)
	err := level.UnmarshalText([]byte(config.GetString(dtos.ConfigKeys.Logger.Level)))
	if err != nil {
		return err
	}
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	fileHandler = slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: level,
	})
	logger = slog.New(fileHandler)
	return nil
}

func Logger(_ context.Context) *slog.Logger {
	return logger
}

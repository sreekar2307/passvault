package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"passVault/cmd/backup"
	"passVault/cmd/migrate"
	"passVault/cmd/server"
	"passVault/resources"
	"syscall"
)

func main() {
	var (
		ctx    = context.Background()
		logger = resources.Logger(ctx)
	)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic occurred in main go-routine ", "error", err)
		}
	}()

	if len(os.Args) < 2 {
		panic("command not provided")
	}

	baseCommand := os.Args[1]

	var (
		signalCancel context.CancelFunc
	)

	ctx, signalCancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer signalCancel()

	switch baseCommand {
	case "server":
		server.Server(ctx, os.Args[2:]...)
	case "migrate":
		migrate.Migrate(ctx, os.Args[2:]...)
	case "backup":
		backup.Backup(ctx, os.Args[2:]...)
	default:
		panic(fmt.Sprintf("%s command not supported", baseCommand))
	}

	logger.Info("closing main go routine...")
}

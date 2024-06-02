package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"passVault/cmd/migrate"
	"passVault/cmd/server"
	"syscall"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("panic occurred in main go-routine ", "error", err)
		}
	}()

	if len(os.Args) < 2 {
		panic("command not provided")
	}

	baseCommand := os.Args[1]

	ctx := context.Background()

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
	default:
		panic(fmt.Sprintf("%s command not supported", baseCommand))
	}

	slog.Info("closing main go routine...")
}

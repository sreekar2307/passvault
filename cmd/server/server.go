package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"passVault/app"
	"passVault/dtos"
	"passVault/resources"
	"time"
)

var serverFlagSet *flag.FlagSet

func init() {
	serverFlagSet = flag.NewFlagSet("server", flag.ExitOnError)
}
func Server(ctx context.Context, args ...string) {
	var (
		config = resources.Config()
		logger = resources.Logger(ctx)
	)
	if err := serverFlagSet.Parse(args); err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s",
			config.GetString(dtos.ConfigKeys.Server.Host),
			config.GetString(dtos.ConfigKeys.Server.Port)),
		Handler: app.SetupRouter(ctx),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Error("failed to start sever...", "error", err.Error())
			} else {
				logger.Info("server was closed...", "error", err.Error())
			}
		}
		logger.Info("closing server go routine...")
	}()

	<-ctx.Done()

	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: ", "error", err.Error())
	}
}

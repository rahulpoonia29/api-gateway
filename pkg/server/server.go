package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rahul/api-gateway/utils"
)

// StartServer initializes and starts the HTTP server for the API gateway
func StartServer(port int, app *utils.App) error {
	// Add context information to logs
	logger := app.Logger.With("component", "server", "port", port)

	addr := fmt.Sprintf(":%d", port)
	hHTTP := &HTTPHandler{
		app: app,
	}

	server := &http.Server{
		Addr:    addr,
		Handler: hHTTP,
	}

	logger.Info("starting api gateway server", "port", port)

	signalChan := make(chan os.Signal, 1)
	// Set up signal notification
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("server closed gracefully")
			return
		} else if err != nil {
			logger.Error("error starting server", "error", err.Error())
		}
	}()

	logger.Info("api gateway server started", "address", addr)

	// Wait for either an OS signal
	<-signalChan

	// Graceful shutdown of the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("error shutting down server", "error", err.Error())
	}

	return nil
}

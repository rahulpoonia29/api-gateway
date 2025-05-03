package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rahul/api-gateway/utils"
)

// StartServer initializes and starts the HTTP server for the API gateway
func StartServer(port int, app *utils.App) error {
	app.Logger.Info("starting API gateway server", "port", port)

	hHTTP := &HTTPHandler{
		app: app,
	}

	http.Handle("/", hHTTP)

	serverAddr := fmt.Sprintf(":%d", port)
	app.Logger.Info("server listening", "address", serverAddr)

	err := http.ListenAndServe(serverAddr, nil)

	// Check for specific errors
	if errors.Is(err, http.ErrServerClosed) {
		app.Logger.Info("server closed")
		return nil
	} else if err != nil {
		app.Logger.Error("error starting server", "error", err)
		return err
	}

	return nil
}

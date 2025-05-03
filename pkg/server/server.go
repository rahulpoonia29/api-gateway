package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/rahul/api-gateway/utils"
)

func StartServer(port int, app *utils.App) {
	app.Logger.Info("starting API gateway server", "port", port)

	hHTTP := &HTTPHandler{
		app: app,
	}

	http.Handle("/", hHTTP)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server is closed")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}

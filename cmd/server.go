package cmd

import (
	"log/slog"
	"os"

	"github.com/armon/go-radix"
	"github.com/rahul/api-gateway/pkg/config"
	"github.com/rahul/api-gateway/pkg/server"
	"github.com/rahul/api-gateway/utils"
	"github.com/spf13/cobra"
)

// NewServerStartCMD creates a new command to start a new http server
func NewServerStartCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the API Gateway server",
		Long:  "This command initializes and starts the API Gateway server, configuring all routes from the configuration file and handling incoming HTTP requests.",
		Run: func(cmd *cobra.Command, args []string) {
			NewServerStart(cmd, args)
		},
	}

	return cmd
}

func NewServerStart(cmd *cobra.Command, args []string) {
	// Initialize structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Get configuration file path from flags
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		logger.Error("error getting config flag", "error", err)
	}

	var defaultConfigPath = "./config.json"
	if configPath == "" {
		configPath = defaultConfigPath
		logger.Warn("config file not specified, using default path", "path", configPath)
	}

	// Validate and load configuration
	gatewayConfig, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		logger.Error("error loading configuration", "error", err)
		os.Exit(1)
	}

	r := radix.New()

	// Insert all the services into the radix tree
	for _, service := range gatewayConfig.Services {
		if !service.Active {
			continue
		}
		r.Insert(service.Proxy.ListenPath, service)
	}

	app := &utils.App{}

	app.RouteTree = r
	app.Logger = logger

	server.StartServer(gatewayConfig.Gateway.Port, app)
}

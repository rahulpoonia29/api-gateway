package cmd

import (
	"log/slog"

	"github.com/armon/go-radix"
	"github.com/rahul/api-gateway/pkg/config"
	"github.com/rahul/api-gateway/pkg/prettylog"

	"github.com/rahul/api-gateway/pkg/server"
	"github.com/rahul/api-gateway/utils"
	"github.com/spf13/cobra"
)

// NewServerStartCMD creates a new command to start a new http server
func NewServerStartCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the API Gateway server",
		Long:  "This command initializes and starts the API Gateway server, configuring all routes from the configuration file and handling incoming HTTP requests.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startServer(cmd)
		},
	}
}

func startServer(cmd *cobra.Command) error {
	logger := prettylog.ConfigureLogger(slog.LevelDebug)

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		logger.Error("Error getting config flag", "error", err)
		return err
	}

	gatewayConfig, err := config.LoadConfig(configPath, logger)
	if err != nil {
		logger.Error("Error loading configuration", "error", err)
		return err
	}

	if gatewayConfig.Gateway.LogLevel != "" {
		if err = prettylog.UpdateLogLevel(logger, gatewayConfig.Gateway.LogLevel); err != nil {
			logger.Warn("Error updating log level, falling back to INFO", "error", err)
			prettylog.UpdateLogLevel(logger, config.Info)
		}
	}

	logger.Debug("Configuration loaded", "services_count", len(gatewayConfig.Services))

	routeTree := loadServices(gatewayConfig.Services)

	app := &utils.App{
		RouteTree: routeTree,
		Logger:    logger,
	}

	return server.StartServer(gatewayConfig.Gateway.Port, app)
}

func loadServices(services []config.ServiceConfig) *radix.Tree {
	r := radix.New()
	for _, service := range services {
		if !service.Enabled {
			continue
		}
		r.Insert(service.Proxy.ListenPath, service)
	}
	return r
}

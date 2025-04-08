package cmd

import (
	"fmt"
	"os"

	"github.com/armon/go-radix"
	"github.com/rahul/api-gateway/pkg/config"
	"github.com/rahul/api-gateway/pkg/server"
	"github.com/rahul/api-gateway/utils"
	"github.com/spf13/cobra"
)

// NewServerStartCMD creates a new command to start a new http server
func NewServerStartCMD(app *utils.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the API Gateway server",
		Long:  "This command starts the API Gateway server and listens for incoming requests.",
		Run: func(cmd *cobra.Command, args []string) {
			NewServerStart(app)
		},
	}

	return cmd
}

func NewServerStart(app *utils.App) {
	var defaultConfigPath = "./config.json"

	gatewayConfig, err := config.LoadConfig(defaultConfigPath)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	r := radix.New()
	for _, service := range gatewayConfig.Services {
		r.Insert(service.Proxy.ListenPath, service)
	}

	app.RouteTree = r

	// if value != nil {
	// 	serviceConfig, ok := value.(config.ServiceConfig)
	// 	if ok {
	// 		fmt.Printf("service name: %s\n", serviceConfig.Proxy.Upstream.Targets[0])
	// 	} else {
	// 		fmt.Println("Could not cast value to ServiceConfig")
	// 	}
	// }

	server.StartServer(gatewayConfig.Gateway.Port, app)
}

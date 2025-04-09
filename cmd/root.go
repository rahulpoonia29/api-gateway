package cmd

import (
	"fmt"
	"os"

	"github.com/armon/go-radix"
	"github.com/rahul/api-gateway/pkg/config"
	"github.com/rahul/api-gateway/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command, called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "api-gateway",
	Short: "A simple API gateway",
	Long:  `A simple API gateway written in Go.`,
}

// Global application instance
var app *utils.App

// Configuration file path
var cfgFile string

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("failed to execute command:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		config.InitConfig(cfgFile)
	})

	// Define flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default is ./config.json)")

	// Initialize app
	app = &utils.App{
		RouteTree: radix.New(),
	}

	// Add commands
	rootCmd.AddCommand(NewServerStartCMD(app))
}

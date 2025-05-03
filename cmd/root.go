package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command, called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "api-gateway",
	Short: "A simple API gateway",
	Long:  `A simple API gateway written in Go. Provides routing, authentication, rate limiting, etc.`,
}

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
	// Define flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default is ./config.json)")

	// Add commands
	rootCmd.AddCommand(NewServerStartCMD())
}

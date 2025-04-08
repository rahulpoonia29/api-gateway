package cmd

import (
	"fmt"
	"os"

	"github.com/armon/go-radix"
	"github.com/rahul/api-gateway/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "api-gateway",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Could not run command")
		os.Exit(1)

	}
}

func init() {

	app := &utils.App{
		RouteTree: radix.New(),
	}

	// Add start command
	rootCmd.AddCommand(NewServerStartCMD(app))
}

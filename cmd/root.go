package cmd

import (
	"os"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lumogon",
	Short: "Lumogon",
	Long:  `Lumogon is a tool for inspecting, reporting on, and analyzing your container applications.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logging.Stderr("Error initialising command: %s", err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&logging.Debug, "debug", "d", false, "Print debug logging")
}

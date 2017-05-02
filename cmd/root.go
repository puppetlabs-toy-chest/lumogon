package cmd

import (
	"os"

	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

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
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cli.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&logging.Debug, "debug", "d", false, "Print debug logging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".lumogon") // name of config file (without extension)
	viper.AddConfigPath("/config")  // adding /config directory as first search path
	viper.BindEnv("LUMOGON_")       // Match all ENV vars starting with LUMOGON_
	viper.AutomaticEnv()            // read in environment variables that match BindEnv

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logging.Stderr("Using config file: %s", viper.ConfigFileUsed())
	}
}

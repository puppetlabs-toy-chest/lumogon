package cmd

import (
	"os"

	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/scheduler"
	"github.com/puppetlabs/lumogon/types"
	"github.com/spf13/cobra"
)

var opts = types.ClientOptions{}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan one or more containers and print the collected information",
	Long:  `Creates and attaches a container to the specified containers, inspect the container and then output that information as JSON to STDOUT`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduler := scheduler.New(args, opts)
		scheduler.Run(registry.Registry)
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Scan one or more containers and send the collected information to the Lumogon service",
	Long:  `Creates and attaches a container to the specified containers, inspect the container and then output that information as JSON over HTTP `,
	Run: func(cmd *cobra.Command, args []string) {
		if opts.ConsumerURL == "" {
			opts.ConsumerURL = os.Getenv("LUMOGON_ENDPOINT")
		}
		scheduler := scheduler.New(args, opts)
		scheduler.Run(registry.Registry)
	},
}

func init() {
	RootCmd.AddCommand(scanCmd)
	RootCmd.AddCommand(reportCmd)
	RootCmd.PersistentFlags().BoolVarP(&opts.KeepHarvesters, "keep-harvesters", "k", false, "Keeps harvester containers instead of automatically deleting")
	reportFlags := reportCmd.Flags()
	scanFlags := scanCmd.Flags()

	reportFlags.StringVar(&opts.ConsumerURL, "endpoint", "", "Use a custom HTTP endpoint for sending the results of the scan")
	reportFlags.IntVar(&opts.Timeout, "timeout", 60, "Time in seconds to wait for results, defaults to 60 seconds")
	scanFlags.IntVar(&opts.Timeout, "timeout", 60, "Time in seconds to wait for results, defaults to 60 seconds")
}

package cmd

import (
	"github.com/puppetlabs/lumogon/analytics"
	"github.com/spf13/cobra"
)

// CapabilityCmd is the capability root command, if invoked without a
// subcommand it will return the cli help output.
var CapabilityCmd = &cobra.Command{
	Use:    "capability",
	Short:  "Capability parent command",
	Hidden: true,
	Long:   `Long Capability Parent command`,
	PreRun: func(cmd *cobra.Command, args []string) {
		analytics.MeasureUsage("capability")
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(CapabilityCmd)
}

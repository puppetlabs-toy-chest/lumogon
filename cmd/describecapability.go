package cmd

import (
	"fmt"

	"os"

	"github.com/puppetlabs/lumogon/analytics"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/spf13/cobra"
)

// harvestCmd captures capability data on the attached container
var describeCapabilityCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe capability",
	Long:  `Long Describe capability`,
	PreRun: func(cmd *cobra.Command, args []string) {
		analytics.MeasureUsage("describe")
	},
	Run: func(cmd *cobra.Command, args []string) {
		description, err := registry.Registry.DescribeCapability(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(description)
	},
}

func init() {
	CapabilityCmd.AddCommand(describeCapabilityCmd)
}

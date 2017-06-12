package cmd

import (
	"os"

	"github.com/puppetlabs/lumogon/analytics"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/harvester"
	"github.com/puppetlabs/lumogon/harvester/rpcreceiver"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"

	"github.com/spf13/cobra"
)

// harvestCmd captures capability data on the attached container
var harvestCmd = &cobra.Command{
	Use:    "harvest",
	Short:  "Harvests metadata from the specified container",
	Long:   `Harvests metadata from the specified container`,
	Hidden: true,
	PreRun: func(cmd *cobra.Command, args []string) {
		analytics.ScreenView("harvest")
	},
	Run: func(cmd *cobra.Command, args []string) {
		target := types.TargetContainer{
			ID:   args[0],
			Name: args[1],
		}
		logging.Debug("Harvesting from container: %s [%s]", target.Name, target.ID)

		containerReport := harvester.GenerateContainerReport(target, registry.Harvest(nil, target.ID))
		harvesterHostname, _ := os.Hostname()
		containerReport.HarvesterContainerID = harvesterHostname
		_, err := rpcreceiver.SendResult(*containerReport, harvesterHostname)
		if err != nil {
			logging.Debug("Error Submitting result to remote server: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(harvestCmd)
}

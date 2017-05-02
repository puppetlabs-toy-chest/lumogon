package cmd

import (
	"os"

	"github.com/puppetlabs/transparent-containers/cli/capabilities/registry"
	"github.com/puppetlabs/transparent-containers/cli/harvester"
	"github.com/puppetlabs/transparent-containers/cli/harvester/rpcreceiver"
	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/types"

	"github.com/spf13/cobra"
)

// harvestCmd captures capability data on the attached container
var harvestCmd = &cobra.Command{
	Use:    "harvest",
	Short:  "Harvests metadata from the specified container",
	Long:   `Harvests metadata from the specified container`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		target := types.TargetContainer{
			ID:   args[0],
			Name: args[1],
		}
		logging.Stderr("Harvesting from container: %s [%s]", target.Name, target.ID)

		containerReport := harvester.GenerateContainerReport(target, registry.Harvest(nil, target.ID))
		harvesterHostname, _ := os.Hostname()
		containerReport.HarvesterContainerID = harvesterHostname
		_, err := rpcreceiver.SendResult(*containerReport, harvesterHostname)
		if err != nil {
			logging.Stderr("Error Submitting result to remote server: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(harvestCmd)
}

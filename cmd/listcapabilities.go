package cmd

import (
	"fmt"

	"github.com/puppetlabs/lumogon/analytics"
	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/spf13/cobra"
)

// harvestCmd captures capability data on the attached container
var listCapabilityCmd = &cobra.Command{
	Use:   "list",
	Short: "List available capabilities",
	Long:  `Long List available capabilities`,
	PreRun: func(cmd *cobra.Command, args []string) {
		analytics.MeasureUsage("list")
	},
	Run: func(cmd *cobra.Command, args []string) {
		attachedCapabilities := registry.Registry.AttachedCapabilities()
		dockerAPICapabilities := registry.Registry.DockerAPICapabilities()
		if len(attachedCapabilities) == 0 && len(dockerAPICapabilities) == 0 {
			fmt.Println("No capabilities found")
		}
		if len(attachedCapabilities) > 0 {
			fmt.Println("Attached Capabilities")
			for _, capability := range attachedCapabilities {
				fmt.Printf(" - %s (%s)\n", capability.Name, capability.Title)
			}
		}
		if len(dockerAPICapabilities) > 0 {
			fmt.Println("Docker API Capabilities Capabilities")
			for _, capability := range dockerAPICapabilities {
				fmt.Printf(" - %s (%s)\n", capability.Name, capability.Title)
			}
		}
	},
}

func init() {
	CapabilityCmd.AddCommand(listCapabilityCmd)
}

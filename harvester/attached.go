package harvester

import (
	"context"
	"fmt"

	"sync"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/harvester/rpcreceiver"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

// RunAttachedHarvester handles gathering AttachedCapabilities from the target containers.
// It creates an RPC Receiver server and RPC results channel which it will use to send
// results back to this method.
// Runs createAndRunHarvester in a goroutine for each of the requested target containers.
// Blocks on receiving a result from each of the target containers via the RPC results
// channel, when a result is received it will attempt to remove that associated
// attached container which performed the harvest before sending the result to the
// collector via the main results channel, resultsCh.
func RunAttachedHarvester(ctx context.Context, wg *sync.WaitGroup, targets []*types.TargetContainer, capabilites []types.AttachedCapability, resultsCh chan types.ContainerReport, opts types.ClientOptions, client dockeradapter.Client) error {
	defer logging.Stderr("[Attached Harvester] Exiting")
	defer wg.Done()

	validTargets := []*types.TargetContainer{}
	for _, target := range targets {
		for _, capability := range capabilites {
			if _, ok := capability.SupportedOS["all"]; ok {
				validTargets = append(validTargets, target)
				break
			}
			if _, ok := capability.SupportedOS[target.OSID]; ok {
				validTargets = append(validTargets, target)
				break
			}
		}
	}

	if len(validTargets) == 0 {
		errorMsg := fmt.Errorf("[Scheduler] No targets found with supported capabilities")
		return errorMsg
	}
	logging.Stderr("[Attached Harvester] Running")
	if len(capabilites) == 0 {
		logging.Stderr("[Attached Harvester] No Attached Capabilities found")
		return nil
	}

	rpcReceiverResultsCh := make(chan types.ContainerReport)

	logging.Stderr("[Attached Harvester] Starting RPC Receiver")
	go rpcreceiver.Run("attachedharvester", 42586, rpcReceiverResultsCh) // TODO Port is still hardcoded client side (probably not an issue?)

	logging.Stderr("[Attached Harvester] Creating [%d] harvesting containers", len(validTargets))
	for _, target := range validTargets {
		go createAndRunHarvester(ctx, client, *target, opts, rpcReceiverResultsCh)
	}

	for _ = range validTargets {
		result := <-rpcReceiverResultsCh
		logging.Stderr("[Attached Harvester] RPC result received from name: %s, ID: %s", result.ContainerName, result.ContainerID)
		removeContainer(ctx, client, result.HarvesterContainerID, opts)
		logging.Stderr("[Attached Harvester] Sending to collector via resultsCh")
		resultsCh <- result
	}
	return nil
}

// createAndRunHarvester creates and runs a container attached to the namespace of the target
// container which will run the harvest command to run the harvest functions from any registered
// AttachedCapabilities.
func createAndRunHarvester(ctx context.Context, client dockeradapter.Client, target types.TargetContainer, opts types.ClientOptions, rpcReceiverResultsCh chan types.ContainerReport) {
	logging.Stderr("[Attached Harvester] Creating attached container for target %s", target)
	harvester := NewAttachedContainer(client, types.ClientOptions{KeepHarvesters: opts.KeepHarvesters})
	// TODO get image name from the current container or set alternate default for non-container use
	harvester.GetImage("puppet/lumogon")
	harvester.Attach(target)
	harvester.Run()
}

package rpcreceiver

import (
	"log"
	"net"
	"net/rpc"

	"fmt"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
)

// Ack contains the reponse from the RPC Receiver server
type Ack bool

// RemoteMethods is a dummy struct used to hang RPC remote methods off.
type RemoteMethods struct {
	resultsCh chan types.ContainerReport
}

// Run creates an RPC Receiver server and registers the remoteMethods which are
// used by attached harvesting containers to submit results back to the scheduler.
func Run(from string, listenPort int, resultsCh chan types.ContainerReport) {
	bindAddress := fmt.Sprintf("0.0.0.0:%d", listenPort)
	logging.Stderr("[RPC Receiver] Starting listener: %s", from)
	address, err := net.ResolveTCPAddr("tcp", bindAddress)
	if err != nil {
		log.Fatal(err)
	}
	logging.Stderr("[RPC Receiver] Bind address resolved: %s", bindAddress)

	inbound, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	logging.Stderr("[RPC Receiver] Listening on: %s", bindAddress)

	logging.Stderr("[RPC Receiver] Registering RemoteMethods")
	remoteMethods := &RemoteMethods{resultsCh: resultsCh}
	rpc.Register(remoteMethods)
	logging.Stderr("[RPC Receiver] Accepting requests")
	rpc.Accept(inbound)
}

// SubmitCapabilities used to submit data to the server
func (r *RemoteMethods) SubmitCapabilities(data *types.ContainerReport, ack *Ack) error {
	logging.Stderr("[RPC Receiver] Results for %d capabilities for container %s received from attached container TODO", len(data.Capabilities), data.ContainerID) // TODO would be useful to be able to get some info about the harvester here.
	logging.Stderr("[RPC Receiver] Sending result to the Attached Harverster result channel")
	r.resultsCh <- *data
	*ack = true
	return nil
}

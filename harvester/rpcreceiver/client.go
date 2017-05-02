package rpcreceiver

import (
	"log"
	"net/rpc"

	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/types"
)

// SendResult submits the ContainerReport populated with the output from each
// AttachedCapability run on the attached container.
// This sends to the `scheduler` host which will have been aliased at container
// creation to the node where the the scheduler is running.
func SendResult(result types.ContainerReport, harvesterHostname string) (bool, error) {
	client, err := rpc.Dial("tcp", "scheduler:42586")
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	var ack bool
	err = client.Call("RemoteMethods.SubmitCapabilities", result, &ack)
	if err != nil {
		logging.Stderr("[RPC Client] Ack received: %t", ack)
		log.Fatal(err)
		return false, err
	}
	return true, nil
}

package collector

import (
	"context"

	"sync"

	"encoding/json"

	"github.com/puppetlabs/transparent-containers/cli/capabilities/registry"
	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/storage"
	"github.com/puppetlabs/transparent-containers/cli/types"
	"github.com/puppetlabs/transparent-containers/cli/utils"
	"github.com/puppetlabs/transparent-containers/cli/version"
)

var mu sync.Mutex
var results map[string]types.ContainerReport

// RunCollector starts the collector which will block on reading all
// expected ContainerReports from the results channel, before creating
// and storing a report.
func RunCollector(ctx context.Context, wg *sync.WaitGroup, targets []types.TargetContainer, resultsCh chan types.ContainerReport, consumerURL string) {
	logging.Stderr("[Collector] Running")
	defer logging.Stderr("[Collector] Exiting")
	defer wg.Done()

	results = make(map[string]types.ContainerReport)
	// Expecting a result per type for each target container
	expectedResults := len(targets) * registry.Registry.TypesCount()

	logging.Stderr("[Collector] Waiting for %d results", expectedResults)
	for i := 1; i <= expectedResults; i++ {
		logging.Stderr("[Collector] Received result [%d]", i)
		result := <-resultsCh
		cacheResult(result)
		logging.Stderr("[Collector] Result received from name: %s, ID: %s", result.ContainerName, result.ContainerID)
	}
	logging.Stderr("[Collector] Creating report")

	report, err := createReport(results)
	if err != nil {
		return
	}
	storeReport(report, consumerURL)

}

// cacheResult caches the supplied types.ContainerReport.
// It consists of a map of container IDs to ContainerReports either adding
// a new key or appending the capabilities to an existing ContainerReport.
func cacheResult(result types.ContainerReport) {
	logging.Stderr("[Collector] Caching result")
	defer logging.Stderr("[Collector] Caching result complete")
	mu.Lock()
	defer mu.Unlock()
	if _, ok := results[result.ContainerID]; ok {
		for capabilityID, capabilityData := range result.Capabilities {
			results[result.ContainerID].Capabilities[capabilityID] = capabilityData
		}
		return
	}
	results[result.ContainerID] = result
}

// createReport returns a pointer to a types.Report built from the supplied
// map of container IDs to types.ContainerReport.
func createReport(results map[string]types.ContainerReport) (types.Report, error) {
	logging.Stderr("[Collector] Creating report")
	marshalledResult, _ := json.Marshal(results)
	logging.Stderr("[Collector] %s", string(marshalledResult))
	report := NewReport(utils.GenerateUUID4(), version.Version)
	report.Containers = results
	logging.Stderr("[Collector] Report created")
	return *report, nil //TODO do we really want a pointer here?
}

// storeReport marshalls the supplied types.Report and sends it to the
// storage package for persistance to the specified consumerURL.
func storeReport(report types.Report, consumerURL string) error {
	logging.Stderr("[Collector] Storing report")
	marshalledReport, err := json.Marshal(report)
	if err != nil {
		logging.Stderr("[Collector] Error marshalling report: %s ", err)
		return err
	}
	err = storage.StoreResult(string(marshalledReport), consumerURL)
	if err != nil {
		logging.Stderr("[Collector] Error storing report: %s ", err)
		return err
	}
	logging.Stderr("[Collector] Report stored")
	return nil
}

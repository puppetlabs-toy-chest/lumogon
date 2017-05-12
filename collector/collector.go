package collector

import (
	"context"
	"sync"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/storage"
	"github.com/puppetlabs/lumogon/types"
)

var mu sync.Mutex
var results map[string]types.ContainerReport

// RunCollector starts the collector which will block on reading all
// expected ContainerReports from the results channel, before sending
// them to the ReportStorage backend.
func RunCollector(ctx context.Context, wg *sync.WaitGroup, expectedResults int, resultsCh chan types.ContainerReport, backend storage.ReportStorage) error {
	defer logging.Stderr("[Collector] Exiting")
	defer wg.Done()

	doneChannel := make(chan int)

	results = make(map[string]types.ContainerReport)

	go func() {
		logging.Stderr("[Collector] Waiting for %d results", expectedResults)
		for i := 1; i <= expectedResults; i++ {
			result := <-resultsCh
			logging.Stderr("[Collector] Received result [%d]", i)
			cacheResult(result)
			logging.Stderr("[Collector] Result received from name: %s, ID: %s", result.ContainerName, result.ContainerID)
		}
		doneChannel <- 0
	}()

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	var err error
	select {
	case <-doneChannel:
		logging.Stderr("[Collector] All expected results received")
		resultsWg.Done()
	case <-ctx.Done():
		logging.Stderr("[Collector] Context timed out waiting for results, continuing...")
		resultsWg.Done()
		err = ctx.Err()
	}
	resultsWg.Wait()

	logging.Stderr("[Collector] Generating report")
	err = backend.Store(results)
	return err
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

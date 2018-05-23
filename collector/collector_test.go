package collector

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/puppetlabs/lumogon/types"
)

// Using a global to capture the report sent to MockStorage.Read
// unable to mutate the receiver as its not a pointer
// TODO is there a better way of doing this? Using a pointer receiver
// just to facilitate testing seemed wrong if we're not altering state
var storedReport map[string]types.ContainerReport

type MockStorage struct{}

func (m MockStorage) Store(containerReports map[string]types.ContainerReport, reportID string) error {
	fmt.Println("entering store")
	fmt.Printf("storing reportID: %s\n", reportID)
	fmt.Printf("attempting to store: %v\n", containerReports)
	defer fmt.Println("exiting store")
	storedReport = containerReports
	return nil
}

var collectorTests = []struct {
	title           string
	ctxTimeoutSec   int
	testTimeoutSec  int
	expectedResults int
	receivedResults int
	expectError     bool
}{
	{
		title:           "All expected results received",
		ctxTimeoutSec:   2,
		testTimeoutSec:  5,
		expectedResults: 2,
		receivedResults: 2,
		expectError:     false,
	},
	{
		title:           "Some expected results timeout",
		ctxTimeoutSec:   2,
		testTimeoutSec:  5,
		expectedResults: 2,
		receivedResults: 1,
		expectError:     false,
	},
	{
		title:           "(Test should time out)",
		ctxTimeoutSec:   5,
		testTimeoutSec:  1,
		expectedResults: 2,
		receivedResults: 1,
		expectError:     true,
	},
}

func Test_collector(t *testing.T) {
	for _, test := range collectorTests {
		// Do not run these in parallel, currently using
		// a global
		t.Run(test.title, func(t *testing.T) {
			if testing.Short() {
				t.Skipf("skipping in short mode: %s", t.Name())
			}

			t.Logf("starting: %s, test timeout: %d", t.Name(), test.testTimeoutSec)
			defer t.Logf("exiting: %s", t.Name())

			c := make(chan error)
			r := make(chan types.ContainerReport)
			var wg sync.WaitGroup
			wg.Add(1)

			t.Logf("creating context with timeout [%d]", test.ctxTimeoutSec)
			testCtx, cancel := context.WithTimeout(context.Background(), time.Duration(test.ctxTimeoutSec)*time.Second)
			defer cancel()

			go func() {
				t.Logf("starting Collector, expectedResults [%d]", test.expectedResults)
				defer t.Logf("exiting Collector")
				RunCollector(testCtx, &wg, test.expectedResults, r, MockStorage{}, "dummy_report_id")
				c <- nil
			}()

			t.Logf("sending results to resultsChannel")
			for i := 1; i <= test.receivedResults; i++ {
				r <- mockContainerReport(i)
			}

			select {
			case <-c:
				if len(storedReport) != test.receivedResults {
					t.Errorf("stored number of results [%d] (map[string]types.ContainerReport), does not match expected [%d]",
						len(storedReport),
						test.receivedResults)
				}
				if !reflect.DeepEqual(storedReport, mockContainerReports(test.receivedResults)) {
					t.Logf("Expected: %v", mockContainerReports(test.receivedResults))
					t.Logf("Actual: %v", storedReport)
					t.Errorf("Stored containerReports does not match expected value")
				}
			case <-time.After(time.Duration(test.testTimeoutSec) * time.Second):
				t.Logf("Test timed out")
				if !test.expectError {
					t.Errorf("Test timeout unexpected")
				}
			}
		})
	}
}

func mockContainerReport(i int) types.ContainerReport {
	return types.ContainerReport{
		Schema:               fmt.Sprintf("testSchema_%d", i),
		Generated:            fmt.Sprintf("testGenerated_%d", i),
		ContainerReportID:    fmt.Sprintf("testContainerReportID_%d", i),
		ContainerID:          fmt.Sprintf("testContainerID_%d", i),
		HarvesterContainerID: fmt.Sprintf("testHarvesterContainerID_%d", i),
		ContainerName:        fmt.Sprintf("testContainerID_%d", i),
		Capabilities: map[string]types.Capability{
			fmt.Sprintf("testCapability_%d", i): types.Capability{
				Name: fmt.Sprintf("testCapabilityName_%d", i),
			},
		},
	}
}

func mockContainerReports(num int) map[string]types.ContainerReport {
	containerReports := map[string]types.ContainerReport{}
	for i := 1; i <= num; i++ {
		containerReports[fmt.Sprintf("testContainerID_%d", i)] = mockContainerReport(i)
	}
	return containerReports
}

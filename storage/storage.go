package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/puppetlabs/lumogon/analytics"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
	"github.com/puppetlabs/lumogon/version"
)

// Storage submits the captured data an appropriate destination
type Storage struct {
	ConsumerURL string
}

// ReportStorage handles persistence of generated container reports
type ReportStorage interface {
	Store(map[string]types.ContainerReport) error
}

// Store marshalls the supplied types.Report before storing it
func (s Storage) Store(results map[string]types.ContainerReport) error {

	report, err := createReport(results)
	if err != nil {
		return err
	}

	logging.Stderr("[Storage] Storing report")
	marshalledReport, err := json.Marshal(report)
	if err != nil {
		logging.Stderr("[Storage] Error marshalling report: %s ", err)
		return err
	}
	err = storeResult(string(marshalledReport), s.ConsumerURL)
	if err != nil {
		logging.Stderr("[Storage] Error storing report: %s ", err)
		return err
	}
	logging.Stderr("[Storage] Report stored")
	return nil
}

// storeResult stores the harvested result, currently this just
// involves printing it to stdout where its manually passed to the
// Lambda consumer
func storeResult(result string, consumerURL string) error {
	var postResponse struct {
		Token string
		URL   string
	}

	if result == "" {
		errorMsg := fmt.Sprintf("[Storage] No harvesting result found")
		logging.Stderr(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	// When not communicating with a remote Consumer, simply output JSON on Stdout
	if consumerURL == "" {
		os.Stdout.WriteString(result + "\n")
		return nil
	}

	jsonStr := []byte(result)
	// TODO Move HTTP Post Helper method elsewhere
	logging.Stderr("[Storage] Posting result to: %s", consumerURL)
	resp, err := http.Post(consumerURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		logging.Stderr("[Storage] Error posting result, [%s], exiting..", err)
		os.Exit(1)
	}

	analytics.Event("upload", "UX")

	err = json.NewDecoder(resp.Body).Decode(&postResponse)
	if err != nil {
		errorMsg := fmt.Sprintf("[Storage] Unable to decode JSON response from server [%s], exiting..", err)
		logging.Stderr(errorMsg)
		os.Exit(1)
	}

	// When dealing with a remote Consumer, output an appropriate report URL
	// TODO Move user interface outside of StorageFunction (return values, handle in sched)
	finalURL := utils.FormatReportURL(postResponse.URL, postResponse.Token)
	fmt.Fprintf(os.Stdout, "\n%s\n", finalURL)

	return nil
}

// createReport returns a pointer to a types.Report built from the supplied
// map of container IDs to types.ContainerReport.
func createReport(results map[string]types.ContainerReport) (types.Report, error) {
	logging.Stderr("[Storage] Marshalling JSON")
	marshalledResult, err := json.Marshal(results)
	if err != nil {
		return types.Report{}, err
	}
	logging.Stderr("[Storage] Marshalling successful %s", string(marshalledResult))

	hostname, _ := os.Hostname()
	report := types.Report{
		Schema:        "http://puppet.com/lumogon/core/draft-01/schema#1",
		Generated:     time.Now().String(),
		Host:          hostname,
		Owner:         "default",
		Group:         []string{"default"},
		ClientVersion: version.Version,
		ReportID:      utils.GenerateUUID4(),
		Containers:    results,
	}
	logging.Stderr("[Storage] Report created")
	return report, nil //TODO do we really want a pointer here?
}

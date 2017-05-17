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

	if s.ConsumerURL == "" {
		return outputResult(report)
	}

	err = storeResult(report, s.ConsumerURL)
	if err != nil {
		logging.Stderr("[Storage] Error storing report: %s ", err)
		return err
	}

	return nil
}

// outputResult pretty-prints a JSON marshalled version of the harvested
// result to STDOUT
func outputResult(report types.Report) error {
	result, err := json.MarshalIndent(report, "", "  ")

	if string(result[:]) == "" {
		errorMsg := fmt.Sprintf("[Storage] No harvesting result found")
		logging.Stderr(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	if err != nil {
		logging.Stderr("[Storage] error marshalling report: %s", err)
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("%s\n", result))
	return nil
}

// storeResult stores the harvested result, posting a JSON-marshalled
// version of the report to the consumerURL.
func storeResult(report types.Report, consumerURL string) error {
	var postResponse struct {
		Token string
		URL   string
	}

	logging.Stderr("[Storage] Storing report")
	result, err := json.Marshal(report)

	if string(result[:]) == "" {
		errorMsg := fmt.Sprintf("[Storage] No harvesting result found")
		logging.Stderr(errorMsg)
		return fmt.Errorf(errorMsg)
	}

	if err != nil {
		logging.Stderr("[Storage] error marshalling report: %s", err)
		return err
	}

	jsonStr := []byte(result)

	// TODO Move HTTP Post Helper method elsewhere
	logging.Stderr("[Storage] Posting result to: %s", consumerURL)
	resp, err := http.Post(consumerURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		logging.Stderr("[Storage] Error posting result, [%s], exiting.", err)
		os.Exit(1)
	}

	analytics.Event("upload", "UX")

	err = json.NewDecoder(resp.Body).Decode(&postResponse)
	if err != nil {
		errorMsg := fmt.Sprintf("[Storage] Unable to decode JSON response from server [%s], exiting.", err)
		logging.Stderr(errorMsg)
		os.Exit(1)
	}

	// output an appropriate report URL
	// TODO Move user interface outside of StorageFunction (return values, handle in sched)
	finalURL := utils.FormatReportURL(postResponse.URL, postResponse.Token)
	fmt.Fprintf(os.Stdout, "\n%s\n", finalURL)

	logging.Stderr("[Storage] Report stored")
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

	report := types.Report{
		Schema:        "http://puppet.com/lumogon/core/draft-01/schema#1",
		Generated:     time.Now().String(),
		Owner:         "default",
		Group:         []string{"default"},
		ClientVersion: version.Version,
		ReportID:      utils.GenerateUUID4(),
		Containers:    results,
	}
	logging.Stderr("[Storage] Report created")
	return report, nil //TODO do we really want a pointer here?
}

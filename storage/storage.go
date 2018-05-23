package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
	"github.com/puppetlabs/lumogon/version"
)

// Storage submits the captured data an appropriate destination
type Storage struct {
	ConsumerURL string
}

// ReportStorage handles persistence of generated container reports, taking
// a map with a report for each container and the unique reportID for the overall
// scan or report.
type ReportStorage interface {
	Store(map[string]types.ContainerReport, string) error
}

// Store marshalls the supplied types.Report before storing it
func (s Storage) Store(results map[string]types.ContainerReport, reportID string) error {
	report, err := createReport(results, reportID)
	if err != nil {
		return err
	}

	if s.ConsumerURL == "" {
		return outputResult(report)
	}

	err = storeResult(report, s.ConsumerURL)
	if err != nil {
		logging.Debug("[Storage] Error storing report id=%s: %s ", reportID, err.Error())
		return err
	}

	return nil
}

// formatReport returns a properly formatted byte array of the JSON
// marshalled version of the given Report, indented or unindented
func formatReport(report types.Report, indent bool) ([]byte, error) {
	var result []byte
	var err error

	if indent {
		result, err = json.MarshalIndent(report, "", "  ")
	} else {
		result, err = json.Marshal(report)
	}

	if err != nil {
		logging.Debug("[Storage] error marshalling report: %s", err)
		return nil, err
	}

	resultString := string(result[:])
	if resultString == "" {
		errorMsg := fmt.Sprintf("[Storage] No harvesting result found")
		logging.Debug(errorMsg)
		return nil, fmt.Errorf(errorMsg)
	}

	return result, nil
}

// outputResult pretty-prints a JSON marshalled version of the harvested
// report to STDOUT
func outputResult(report types.Report) error {
	result, err := formatReport(report, true) // indented report
	if err != nil {
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

	logging.Debug("[Storage] Storing report")
	jsonStr, err := formatReport(report, false) // unindented report
	if err != nil {
		return err
	}

	// TODO Move HTTP Post Helper method elsewhere
	logging.Debug("[Storage] Posting result to: %s", consumerURL)
	resp, err := http.Post(consumerURL, "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		logging.Debug("[Storage] Error posting result, [%s], exiting.", err)
		os.Stderr.WriteString(fmt.Sprintf("Unable to connect to the HTTP endpoint at %s\nSending scan report to fallback: STDOUT\n", consumerURL))
		outputResult(report)
		os.Exit(1)
	}

	err = json.NewDecoder(resp.Body).Decode(&postResponse)
	if err != nil {
		errorMsg := fmt.Sprintf("[Storage] Unable to decode JSON response from server [%s], exiting.", err)
		logging.Debug(errorMsg)
		os.Exit(1)
	}

	// output an appropriate report URL
	// TODO Move user interface outside of StorageFunction (return values, handle in sched)
	finalURL := utils.FormatReportURL(postResponse.URL, postResponse.Token)
	fmt.Fprintf(os.Stdout, "\n%s\n", finalURL)

	logging.Debug("[Storage] Report stored")
	return nil
}

// createReport returns a pointer to a types.Report built from the supplied
// map of container IDs to types.ContainerReport.
func createReport(results map[string]types.ContainerReport, reportID string) (types.Report, error) {
	logging.Debug("[Storage] Marshalling JSON")
	marshalledResult, err := json.Marshal(results)
	if err != nil {
		return types.Report{}, err
	}
	logging.Debug("[Storage] Marshalling successful %s", string(marshalledResult))

	report := types.Report{
		Schema:        "http://puppet.com/lumogon/core/draft-01/schema#1",
		Generated:     time.Now().String(),
		Owner:         "default",
		Group:         []string{"default"},
		ClientVersion: version.Version,
		ReportID:      reportID,
		Containers:    results,
	}
	logging.Debug("[Storage] Report created")
	return report, nil //TODO do we really want a pointer here?
}

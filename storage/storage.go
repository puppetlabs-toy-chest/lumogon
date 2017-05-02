package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/utils"
)

// Storage submits the captured data an appropriate destination
type Storage interface {
	StoreResult(result string, consumerURL string) error
}

// StoreResult stores the harvested result, currently this just
// involves printing it to stdout where its manually passed to the
// Lambda consumer
func StoreResult(result string, consumerURL string) error {
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
	resp, err := http.Post(consumerURL, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		errorMsg := fmt.Sprintf("[Storage] Unable to post message, [%s], exiting..", err)
		logging.Stderr(errorMsg)
		os.Exit(1)
	}

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

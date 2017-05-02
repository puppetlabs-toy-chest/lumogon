package utils

import (
	"strings"
	"testing"
)

func Test_FormatReportURL_trailingSlash(t *testing.T) {
	consumerURL := "https://url-with-trailing-slash/"
	token := "beep-beep-imma-token"
	expectedReportURL := "https://url-with-trailing-slash/beep-beep-imma-token"
	reportURL := FormatReportURL(consumerURL, token)
	if strings.Compare(reportURL, expectedReportURL) != 0 {
		t.Errorf("Formatted reportURL: %s, does not match expected result: %s", reportURL, expectedReportURL)
	}
}

func Test_FormatReportURL_noTrailingSlash(t *testing.T) {
	consumerURL := "https://url-with-trailing-slash"
	token := "beep-beep-imma-token"
	expectedReportURL := "https://url-with-trailing-slash/beep-beep-imma-token"
	reportURL := FormatReportURL(consumerURL, token)
	if strings.Compare(reportURL, expectedReportURL) != 0 {
		t.Errorf("Formatted reportURL: %s, does not match expected result: %s", reportURL, expectedReportURL)
	}
}

func Test_FormatReportURL_missingURL(t *testing.T) {
	consumerURL := ""
	token := "beep-beep-imma-token"
	expectedReportURL := "http://localhost:3000/beep-beep-imma-token"
	reportURL := FormatReportURL(consumerURL, token)
	if strings.Compare(reportURL, expectedReportURL) != 0 {
		t.Errorf("Formatted reportURL: %s, does not match expected result: %s", reportURL, expectedReportURL)
	}
}

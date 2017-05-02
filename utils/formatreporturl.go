package utils

import "strings"

// FormatReportURL parses and returns the correct reporting URL
func FormatReportURL(consumerURL string, token string) string {
	if consumerURL == "" {
		consumerURL = "http://localhost:3000"
	}

	url := []string{strings.Trim(consumerURL, "/"), token}
	return strings.Join(url, "/")
}

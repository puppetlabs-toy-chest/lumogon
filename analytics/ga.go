package analytics

import (
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

// MeasureUsage decides whether to proceed with telemetry, and runs a goroutine to post
func MeasureUsage(command string) {
	if viper.GetBool("disable-analytics") != true {
		PostMeasurement(command)
	}
}

// PostMeasurement sends telemetry data to Google Analytics
func PostMeasurement(cd string) {
	vals := make(url.Values, 0)
	vals.Add("v", "1")               // GA Measurement Protocol Version
	vals.Add("tid", "UA-54263865-7") // Tracking ID
	vals.Add("cid", "1")             // Unique ID
	vals.Add("an", "lumogen")        // Application Name
	vals.Add("av", "0.0.0")          // Application Version
	vals.Add("t", "event")           // Type: Event
	vals.Add("ec", "UX")             // Event Category
	vals.Add("ea", cd)               // Event Action

	uri := "https://www.google-analytics.com/collect?" + vals.Encode()

	_, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
}

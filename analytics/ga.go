package analytics

import (
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/puppetlabs/lumogon/version"
	"github.com/spf13/viper"
)

// GaHost is the Hostname for Google Analytics Measurement Endpoint
const GaHost = "www.google-analytics.com"

// ScreenView contains information about a given Application
// being accessed
type ScreenView struct {
	ProtocolVersion    int    `url:"v"`
	Type               string `url:"t"`
	ScreenName         string `url:"cd"`
	ApplicationName    string `url:"an"`
	ApplicationVersion string `url:"av"`
	TrackingID         string `url:"tid"`
	UniqueID           string `url:"cid"`
}

// NewScreenView returns a new ScreenView struct pre-populated with
// everything but the ScreenName being viewed
func NewScreenView() *ScreenView {
	v := version.Version
	uid := "1"

	return &ScreenView{
		ProtocolVersion:    1,
		TrackingID:         "UA-54263865-7",
		Type:               "screenview",
		ApplicationName:    "lumogon",
		ApplicationVersion: "0.0.0", // v.VersionString(),
		UniqueID:           uid,
	}
}

// PostMeasurement sends telemetry data to Google Analytics
func (s *ScreenView) PostMeasurement() bool {
	v, _ := query.Values(s)
	client := http.DefaultClient
	req := &http.Request{
		Method: "GET",
		Host:   GaHost,
		URL: &url.URL{
			Host:     GaHost,
			Scheme:   "https",
			Path:     "/collect",
			RawQuery: v.Encode(),
		},
	}

	client.Do(req)
	return true
}

// MeasureUse checks to see if it's authorized to post analytics events,
// and proceeds accordingly
func (s *ScreenView) MeasureUse() bool {
	if viper.GetBool("disable-analytics") != true {
		return s.PostMeasurement()
	}
	return true
}

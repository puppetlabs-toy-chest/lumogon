package analytics

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/version"
	"github.com/spf13/viper"
)

// GaHost is the Hostname for Google Analytics Measurement Endpoint
const GaHost = "www.google-analytics.com"

// ScreenViewMessage returns a struct representing a ScreenView on GA
type ScreenViewMessage struct {
	ScreenName string `url:"cd,omitempty"`
}

// EventMessage returns a struct representing an Event on GA
type EventMessage struct {
	Category string `url:"ec,omitempty"`
	Action   string `url:"ea,omitempty"`
}

// UserSession is a struct containing connection details and shared
// context among GA hit types
type UserSession struct {
	ProtocolVersion    int          `url:"v"`
	Type               string       `url:"t"`
	UniqueID           string       `url:"cid"`
	TrackingID         string       `url:"tid"`
	ApplicationName    string       `url:"an"`
	ApplicationVersion string       `url:"av"`
	DisableTransmit    bool         `url:"-"`
	HTTPClient         *http.Client `url:"-"`
	ScreenViewMessage
	EventMessage
}

// NewUserSession provides a setup UserSession struct with sane defaults
func NewUserSession() *UserSession {
	v := version.Version
	ctx := context.Background()
	c, _ := dockeradapter.New()

	return &UserSession{
		ProtocolVersion:    1,
		TrackingID:         "UA-54263865-7",
		ApplicationName:    "lumogon",
		ApplicationVersion: v.VersionString(),
		UniqueID:           c.HostID(ctx),
		DisableTransmit:    viper.GetBool("disable-analytics"),
		HTTPClient:         http.DefaultClient,
	}
}

// ScreenView is the public function to post a ScreenView analytics message to GA
func ScreenView(screen string) {
	logging.Debug("[Analytics] Initializing Google Analytics: %s", screen)
	u := *NewUserSession()
	u.Type = "screenview"
	u.ScreenName = screen
	go u.PostMeasurement()
}

// Event is the public function to post a ScreenView event analytics message to GA
func Event(action string, category string) {
	logging.Debug("[Analytics] Gathering additional Google Analytics for event: %s", category)
	u := *NewUserSession()
	u.Type = "event"
	u.Action = action
	u.Category = category
	go u.PostMeasurement()
}

// PostMeasurement is the internal function responsible for calling upstream
// Makes a determiniation whether to post based on User Input via Viper
func (u UserSession) PostMeasurement() (*http.Response, error) {
	v, _ := query.Values(u)
	req := &http.Request{
		Method: "POST",
		Host:   GaHost,
		URL: &url.URL{
			Host:     GaHost,
			Scheme:   "https",
			Path:     "/collect",
			RawQuery: v.Encode(),
		},
	}

	if u.DisableTransmit == true {
		logging.Debug("[Analytics] Skipping submission of Google Analytics event")
		return nil, nil
	}

	logging.Debug("[Analytics] Submitting event to Google Analytics")
	return u.HTTPClient.Do(req)
}

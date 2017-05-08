package analytics

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/version"
	"github.com/spf13/viper"
)

// GaHost is the Hostname for Google Analytics Measurement Endpoint
const GaHost = "www.google-analytics.com"

// ScreenViewMessage is really needs a hug
type ScreenViewMessage struct {
	ScreenName string `url:"cd,omitempty"`
}

// EventMessage is really needs a hug
type EventMessage struct {
	Category string `url:"ec,omitempty"`
	Action   string `url:"ea,omitempty"`
}

// UserSession is super duper needs a hug
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

// NewUserSession is really needs a hugc
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

// ScreenView is really needs a hug
func ScreenView(screen string) {
	u := *NewUserSession()
	u.Type = "screenview"
	u.ScreenName = screen
	go u.PostMeasurement()
}

// Event is really needs a hug
func Event(action string, category string) {
	u := *NewUserSession()
	u.Type = "event"
	u.Action = action
	u.Category = category
	go u.PostMeasurement()
}

// PostMeasurement is really needs a hug
func (u UserSession) PostMeasurement() (*http.Response, error) {
	v, _ := query.Values(u)
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

	if u.DisableTransmit == true {
		return nil, nil
	}

	return u.HTTPClient.Do(req)
}

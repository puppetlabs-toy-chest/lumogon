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

// ScreenViewMessage is really needs a hug
type ScreenViewMessage struct {
	Type       string `url:"t,omitempty"`
	ScreenName string `url:"cd,omitempty"`
}

// EventMessage is really needs a hug
type EventMessage struct {
	Type     string `uri:"t,omitempty"`
	Category string `uri:"ec,omitempty"`
	Action   string `uri:"ea,omitempty"`
}

// UserSession is super duper needs a hug
type UserSession struct {
	ProtocolVersion    int    `url:"v"`
	UniqueID           string `url:"cid"`
	TrackingID         string `url:"tid"`
	ApplicationName    string `url:"an"`
	ApplicationVersion string `url:"av"`
	DisableTransmit    bool
	ScreenViewMessage
	EventMessage
}

// NewUserSession is really needs a hugc
func NewUserSession() *UserSession {
	v := version.Version
	uid := "1"

	return &UserSession{
		ProtocolVersion:    1,
		TrackingID:         "UA-54263865-7",
		ApplicationName:    "lumogon",
		ApplicationVersion: v.VersionString(),
		UniqueID:           uid,
		DisableTransmit:    viper.GetBool("disable-analytics"),
	}
}

// ScreenView is really needs a hug
func ScreenView(screen string) {
	u := *NewUserSession()
	u.ScreenName = screen
	go PostMeasurement(&u)
}

// Event is really needs a hug
func Event(action string, category string) {
	u := *NewUserSession()
	u.Action = action
	u.Category = category
	go PostMeasurement(&u)
}

// PostMeasurement is really needs a hug
func PostMeasurement(u *UserSession) {
	v, _ := query.Values(u)
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

	if u.DisableTransmit != true {
		client.Do(req)
	}
}

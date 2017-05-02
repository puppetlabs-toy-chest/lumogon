package types

import "github.com/puppetlabs/transparent-containers/cli/version"

// Report Lumogon metadata report includes details of the client performing
// the capture and a map of all harvested containers and their capability data
type Report struct {
	Schema        string                     `json:"$schema"`
	Generated     string                     `json:"generated"`
	Owner         string                     `json:"owner"`
	Group         []string                   `json:"group"`
	ClientVersion version.ClientVersion      `json:"client_version"`
	ReportID      string                     `json:"reportid"`
	Containers    map[string]ContainerReport `json:"containers"`
}

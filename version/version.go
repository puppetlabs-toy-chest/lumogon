package version

var (
	// BuildVersion set at build time
	BuildVersion string
	// BuildTime set at build time
	BuildTime string
	// BuildSHA set at build time
	BuildSHA string
)

// ClientVersion contains information about the current client
type ClientVersion struct {
	BuildVersion string
	BuildTime    string
	BuildSHA     string
}

// Version constructed at build time
var Version = ClientVersion{BuildVersion, BuildTime, BuildSHA}

// VersionString returns the absolute version (AppVer-SHA)
func (c *ClientVersion) VersionString() string {
	return c.BuildVersion + "-" + c.BuildSHA
}

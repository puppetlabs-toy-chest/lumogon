package cmd

import (
	"bytes"
	"text/template"

	"fmt"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/version"
	"github.com/spf13/cobra"
)

var versionTemplate = `Client:
 Version:      {{.BuildVersion}}
 Git commit:   {{.BuildSHA}}
 Built:        {{.BuildTime}}`

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the Lumogon version information",
	Long:  `Returns the build version, time and SHA of the current Lumogon client.`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion() {
	// Dump the rendered version template to stdout
	renderedVersionOutput, err := renderVersionTemplate()
	if err != nil {
		logging.Stderr("Unable to render version template [%s]", err)
	}
	fmt.Println(renderedVersionOutput)
}

func renderVersionTemplate() (string, error) {
	versionTemplate, err := template.New("version").Parse(versionTemplate)
	if err != nil {
		return "", err
	}
	var versionOutputBuffer bytes.Buffer
	err = versionTemplate.Execute(&versionOutputBuffer, version.Version)
	if err != nil {
		return "", err
	}
	return versionOutputBuffer.String(), nil
}

package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/pkg/archive"
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/plugin"
	"github.com/puppetlabs/lumogon/types"
)

func main() {
	fmt.Println("main function in pluginapiimpl")
}

// CustomPlugin TODO
type CustomPlugin struct{}

// Print TODO
func (c CustomPlugin) Print(client dockeradapter.Harvester, ID string, target types.TargetContainer) {
	fmt.Printf("I AM IN A PLUGING YEOOOO FFS: %s\n", ID)
}

// Metadata TODO
func (c CustomPlugin) Metadata() *plugin.Metadata {
	return &plugin.Metadata{
		Schema:      "http://puppet.com/lumogon/capability/diff/draft-01/schema#1",
		ID:          "diff_plugin",
		Name:        "Changed Files",
		Description: `The diff capability returns files changed from the initial image as a map["changed file"]"change type"`,
		Type:        plugin.DockerAPI,
		Version:     "0.0.1",
		GitSHA:      "yoIHeardYouLikeGitSHAs",
		SupportedOS: map[string]int{"all": 1},
	}
}

// Harvest TODO
func (c CustomPlugin) Harvest(client dockeradapter.Harvester, ID string, target types.TargetContainer) (map[string]interface{}, error) {
	ctx := context.Background()

	changedFiles, err := getChangedFiles(ctx, client, ID, target)
	if err != nil {
		logging.Debug("[Plugin Diff] Error getting changed files: %v", err)
		return nil, err
	}

	filtered, err := payloadfilter.Filter(changedFiles)
	if err != nil {
		logging.Debug("[Plugin Diff] Error filtering changedFiles output: %v", changedFiles)
		return nil, err
	}
	return filtered, nil
}

// LumogonPluginImpl TODO
var LumogonPluginImpl plugin.Minimal = CustomPlugin{}

func getChangedFiles(ctx context.Context, client dockeradapter.Diff, id string, target types.TargetContainer) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	diffs, err := client.ContainerDiff(ctx, target.ID)
	if err != nil {
		errorMsg := fmt.Sprintf("[Plugin Diff] Error getting diff from targetContainer: %s, error: %s", target.Name, err)
		logging.Debug(errorMsg)
		return nil, err
	}

	for _, diff := range diffs {
		logging.Debug("[Plugin Diff]   Path: %s, Kind %d", diff.Path, diff.Kind)
		var kind string
		switch diff.Kind {
		case archive.ChangeModify:
			kind = "Modified"
		case archive.ChangeAdd:
			kind = "Added"
		case archive.ChangeDelete:
			kind = "Deleted"
		}
		result[diff.Path] = kind
		logging.Debug("[Plugin Diff] EXITING PLUGIN AFTER ONE FILE")
		break // Only get one file for the demo
	}
	return result, nil
}

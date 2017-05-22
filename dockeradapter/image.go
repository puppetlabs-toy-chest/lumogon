package dockeradapter

import (
	"context"
	"regexp"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/utils"
)

// LocalImageID returns the imageID sha256 hex for the local container
func LocalImageID(ctx context.Context, client Inspector) (*string, error) {
	localContainerID, err := utils.GetLocalContainerID("/proc/self/cgroup")
	if err != nil {
		logging.Stderr("[LocalImageID] unable to determine local ContainerID: %v", err)
		return nil, err
	}
	var imageIDRegex = regexp.MustCompile(`^sha256:([a-z0-9]+)`)

	json, err := client.ContainerInspect(ctx, localContainerID)
	if err != nil {
		logging.Stderr("[LocalImageID] error inspecting local container [ID: %s]: %v", localContainerID)
		return nil, err
	}

	imageID := imageIDRegex.FindStringSubmatch(json.Image)[1]
	logging.Stderr("[LocalImageID] found local imageID: %s", imageID)
	return &imageID, nil
}

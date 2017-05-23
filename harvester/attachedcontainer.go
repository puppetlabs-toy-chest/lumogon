package harvester

import (
	"context"
	"fmt"

	"os"

	dockercontainer "github.com/docker/docker/api/types/container"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// AttachedContainer is a container attached to a running target container
// used to gather capability information
type AttachedContainer struct {
	imageName     string
	id            string
	name          string
	start         string
	end           string
	schedulerID   string
	ctx           context.Context
	result        *types.ContainerReport
	target        types.TargetContainer
	containerBody *dockercontainer.ContainerCreateCreatedBody
	client        dockeradapter.Client
	keepHarvester bool
	err           error
}

// AttachedContainerInterface interface exposes Harvester lifecycle functions
type AttachedContainerInterface interface {
	Attach(target types.TargetContainer)
	Run()
	Result() (*types.ContainerReport, error)
}

// NewAttachedContainer returns a pointer to a harvester AttachedContainer
func NewAttachedContainer(client dockeradapter.Client, opts types.ClientOptions) *AttachedContainer {

	hostname, _ := os.Hostname()

	attachedContainer := &AttachedContainer{
		start:         utils.GetTimestamp(),
		name:          utils.GetRandomName("lumogon_"),
		client:        client,
		schedulerID:   hostname,
		keepHarvester: opts.KeepHarvesters,
		ctx:           context.Background(),
	}

	return attachedContainer
}

// Attach creates a harvest container attached to the target container
func (a *AttachedContainer) Attach(target types.TargetContainer) {
	// Exit immediately if a harvester error has already been thrown
	if a.err != nil {
		return
	}
	a.target = target
	a.createContainer()
}

// GetImage ensures that the requested image exists, pulling it if its not already
// present on the system
func (a *AttachedContainer) GetImage(imageName string) {
	if !dockeradapter.ImageExists(a.ctx, a.client, imageName) {
		logging.Stderr("[AttachedContainer] Pulling image: %s", imageName)
		err := a.client.ImagePull(a.ctx, imageName)
		if err != nil {
			a.err = err
		}
	}
	a.imageName = imageName // TODO where to put this?
}

// Run starts the attached harvester container and captures the Lumogon ContainerReport result
func (a *AttachedContainer) Run() {
	// Exit immediately if a harvester error has already been thrown
	if a.err != nil {
		return
	}

	logging.Stderr("[AttachedContainer] Starting harvester ID: %s, attached to: %s [%s]", a.id, a.target.ID, a.target.Name)
	err := a.client.ContainerStart(a.ctx, a.id)
	if err != nil {
		a.err = err
	}
}

// createContainer creates the AttachedContainer attaching it to the target container
func (a *AttachedContainer) createContainer() {
	logging.Stderr("[AttachedContainer] Creating container: %s", a.name)
	logging.Stderr("[AttachedContainer] Attaching container to ID: %s, Name: %s", a.target.ID, a.target.Name)

	command := []string{"harvest", a.target.ID, a.target.Name, "-d"}
	// Envvars used by gopsutil to query attached container
	envvars := []string{
		"HOST_PROC=/proc/1/root/proc",
		"HOST_ETC=/proc/1/root/etc",
		"HOST_SYS=/proc/1/root/sys",
	}

	binds := []string{"/var/run/docker.sock:/var/run/docker.sock"}
	kernelCapabilities := []string{"sys_admin"} // TODO - Need to investigate making the harvester immutable? minimise risk of altering attached namespace
	pidMode := fmt.Sprintf("container:%s", a.target.ID)
	schedulerAliasHostname := "scheduler"
	// Add an aliass for the scheduler to each harvester
	links := []string{fmt.Sprintf("%s:%s", a.schedulerID, schedulerAliasHostname)}

	container, err := a.client.ContainerCreate(a.ctx, command, envvars, a.imageName, binds, links, kernelCapabilities, pidMode, a.name, !a.keepHarvester)
	if err != nil {
		logging.Stderr("[AttachedContainer] Error creating container: %s", err)
		a.err = err
		return
	}
	a.containerBody = &container
	a.id = container.ID
}

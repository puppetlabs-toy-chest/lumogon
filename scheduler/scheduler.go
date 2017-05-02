package scheduler

import (
	"context"
	"fmt"

	"sync"

	"github.com/puppetlabs/transparent-containers/cli/capabilities/registry"
	"github.com/puppetlabs/transparent-containers/cli/collector"
	"github.com/puppetlabs/transparent-containers/cli/dockeradapter"
	"github.com/puppetlabs/transparent-containers/cli/harvester"
	"github.com/puppetlabs/transparent-containers/cli/logging"
	"github.com/puppetlabs/transparent-containers/cli/types"
	"github.com/puppetlabs/transparent-containers/cli/utils"
)

// Scheduler runs in the Lumogon client container and handles gathering
// capabilities, building report data and submitting to the consumer
// endpoint.
type Scheduler struct {
	harvesters   []harvester.AttachedContainer
	capabilities registry.CapabilitiesRegistry
	targets      []types.TargetContainer
	client       dockeradapter.Client
	report       types.Report
	args         *[]string
	opts         *types.ClientOptions
	start        string
	end          string
	err          error
}

var wg sync.WaitGroup

// New returns a pointer to a Scheduler
func New(args []string, opts types.ClientOptions) *Scheduler {
	logging.Stderr("[Scheduler] Creating scheduler")
	scheduler := Scheduler{
		start: utils.GetTimestamp(),
		args:  &args,
		opts:  &opts,
	}

	client, err := dockeradapter.New()
	if err != nil {
		scheduler.err = err
	}
	scheduler.client = client

	return &scheduler
}

// Run starts the scheduler
func (s *Scheduler) Run() {
	defer logging.Stderr("[Scheduler] Exiting")
	logging.Stderr("[Scheduler] Running")
	// Exit immediately if a harvester error has already been thrown
	if s.err != nil {
		return
	}

	ctx := context.Background()
	resultsChannel := make(chan types.ContainerReport)
	s.getTargetContainers()
	wg.Add(1)
	go collector.RunCollector(ctx, &wg, s.targets, resultsChannel, s.opts.ConsumerURL)

	wg.Add(1)
	err := harvester.RunAttachedHarvester(ctx, &wg, s.targets, registry.Registry.AttachedCapabilities(), resultsChannel, *s.opts, s.client)
	if err != nil {
		logging.Stderr("[Scheduler] Error running Attached harvesters")
	}

	wg.Add(1)
	err = harvester.RunDockerAPIHarvester(ctx, &wg, s.targets, registry.Registry.DockerAPICapabilities(), resultsChannel, s.client)
	if err != nil {
		logging.Stderr("[Scheduler] Error running Docker API harvesters")
	}

	wg.Wait()
}

func (s *Scheduler) getTargetContainers() {
	ctx := context.Background()
	if len(*s.args) > 0 {
		s.targets = stringsToTargetContainers(ctx, *s.args, s.client)
	} else {
		targetContainerIDs, err := s.client.ContainerList(ctx)
		if err != nil {
			errorMsg := fmt.Errorf("[Scheduler] Unable to list containers, error: %s", err)
			s.err = errorMsg
		}

		localContainerID, err := utils.GetLocalContainerID("/proc/self/cgroup")
		if err == nil {
			logging.Stderr("[Scheduler] Excluding scheduler container from harvested containers, ID: %s", localContainerID)
			targetContainerIDs = utils.RemoveStringFromSlice(targetContainerIDs, localContainerID)
		}
		s.targets = stringsToTargetContainers(ctx, targetContainerIDs, s.client)
	}
}

// stringToTargetContainer converts a container ID or Name string into types.TargetContainer
func stringToTargetContainer(ctx context.Context, containerIDOrName string, client dockeradapter.Inspector) (types.TargetContainer, error) {
	containerJSON, err := client.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		error := fmt.Sprintf("[Scheduler] Unable to find target container: %s, error: %s", containerIDOrName, err)
		logging.Stderr(error)
		return types.TargetContainer{}, err
	}
	targetContainer := types.TargetContainer{
		ID:   containerJSON.ContainerJSONBase.ID,
		Name: containerJSON.ContainerJSONBase.Name,
	}
	return targetContainer, nil
}

// stringsToTargetContainers converts a slice of container ID or Name strings into a slice of types.TargetContainer
func stringsToTargetContainers(ctx context.Context, containerIDsOrNames []string, client dockeradapter.Inspector) []types.TargetContainer {
	targetContainers := []types.TargetContainer{}
	for _, containerIDOrName := range containerIDsOrNames {
		targetContainer, err := stringToTargetContainer(ctx, containerIDOrName, client)
		if err != nil {
			continue
		}
		targetContainers = append(targetContainers, targetContainer)
	}
	return targetContainers
}

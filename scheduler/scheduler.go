package scheduler

import (
	"context"
	"os"

	"sync"

	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/collector"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/harvester"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/storage"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// Scheduler runs in the Lumogon client container and handles gathering
// capabilities, building report data and submitting to the consumer
// endpoint.
type Scheduler struct {
	harvesters   []harvester.AttachedContainer
	capabilities registry.CapabilitiesRegistry
	targets      []*types.TargetContainer
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
func (s *Scheduler) Run(r registry.IRegistry) {
	defer logging.Stderr("[Scheduler] Exiting")
	logging.Stderr("[Scheduler] Running")
	// Exit immediately if a harvester error has already been thrown
	if s.err != nil {
		return
	}

	ctx := context.Background()
	resultsChannel := make(chan types.ContainerReport)
	targets, err := dockeradapter.NormaliseTargets(ctx, s.args, s.client)
	if err != nil {
		logging.Stderr("[Scheduler] Unable to normalise targets: %s. Exiting...", err)
		os.Exit(1)
	}
	s.targets = targets

	expectedResultCount := getExpectedResultCount(s.targets, r)

	storageBackend := storage.Storage{ConsumerURL: s.opts.ConsumerURL}
	wg.Add(1)
	go collector.RunCollector(ctx, &wg, expectedResultCount, resultsChannel, storageBackend)

	wg.Add(1)
	err = harvester.RunAttachedHarvester(ctx, &wg, s.targets, r.AttachedCapabilities(), resultsChannel, *s.opts, s.client)
	if err != nil {
		logging.Stderr("[Scheduler] Error running Attached harvesters: %s", err)
	}

	wg.Add(1)
	err = harvester.RunDockerAPIHarvester(ctx, &wg, s.targets, r.DockerAPICapabilities(), resultsChannel, s.client)
	if err != nil {
		logging.Stderr("[Scheduler] Error running Docker API harvesters")
	}

	wg.Wait()
}

func getExpectedResultCount(targets []*types.TargetContainer, registry registry.IRegistry) int {
	expectedResults := 0
	for _, target := range targets {
		for _, capability := range registry.AttachedCapabilities() {
			if utils.KeyInMap("all", capability.SupportedOS) || utils.KeyInMap(target.OSID, capability.SupportedOS) {
				expectedResults++
				break
			}
		}
		for _, capability := range registry.DockerAPICapabilities() {
			if utils.KeyInMap("all", capability.SupportedOS) || utils.KeyInMap(target.OSID, capability.SupportedOS) {
				expectedResults++
				break
			}
		}
	}
	return expectedResults
}

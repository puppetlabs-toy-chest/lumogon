package scheduler

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/puppetlabs/lumogon/capabilities/registry"
	"github.com/puppetlabs/lumogon/collector"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/dockeradapter/versions"
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
	reportID     string
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
	logging.Debug("[Scheduler] Creating scheduler")
	scheduler := Scheduler{
		start:    utils.GetTimestamp(),
		args:     &args,
		opts:     &opts,
		reportID: utils.GenerateUUID4(),
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
	defer logging.Debug("[Scheduler] Exiting")
	logging.Debug("[Scheduler] Running")
	// Exit immediately if a harvester error has already been thrown
	if s.err != nil {
		return
	}

	timeout := s.opts.Timeout
	logging.Debug("[Scheduler] Creating context with timeout [%d]", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	resultsChannel := make(chan types.ContainerReport)

	targets, err := dockeradapter.NormaliseTargets(ctx, s.args, s.client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to normalise target containers: %s.\nExiting...", err)
		os.Exit(1)
	}
	if len(targets) == 0 {
		fmt.Fprintln(os.Stderr, "No valid or running target containers found. Exiting...")
		os.Exit(1)
	}
	s.targets = targets

	expectedResultCount := getExpectedResultCount(s.targets, r)

	storageBackend := storage.Storage{ConsumerURL: s.opts.ConsumerURL}
	wg.Add(1)
	go collector.RunCollector(ctx, &wg, expectedResultCount, resultsChannel, storageBackend, s.reportID)

	wg.Add(1)
	go harvester.RunAttachedHarvester(ctx, &wg, s.targets, r.AttachedCapabilities(), resultsChannel, *s.opts, s.client, s.reportID)

	wg.Add(1)
	go harvester.RunDockerAPIHarvester(ctx, &wg, s.targets, r.DockerAPICapabilities(), resultsChannel, s.client)

	logging.Debug("[Scheduler] Waiting")
	wg.Wait()

	if !s.opts.KeepHarvesters {
		autoRemoveSupportedAPIVersion := "1.25"
		if versions.LessThan(s.client.ServerAPIVersion(), autoRemoveSupportedAPIVersion) {
			logging.Debug("[Scheduler] Cleaning up harvester containers explicitly as Server API version %s < %s", s.client.ServerAPIVersion(), autoRemoveSupportedAPIVersion)
			ctx := context.Background()
			harvesterLabel := fmt.Sprintf("lumogon_report_id=%s", s.reportID)
			if err := s.client.CleanupHarvesters(ctx, harvesterLabel); err != nil {
				fmt.Fprintf(os.Stderr, "Error returned when deleting containers: %s", err.Error())
			}
		}
	}
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

package harvester

import (
	"context"
	"io/ioutil"
	"plugin"

	"sync"

	"fmt"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	lumogonplugin "github.com/puppetlabs/lumogon/plugin"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

func getPlugins(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logging.Debug("Unable to read plugins")
		return nil, err
	}

	logging.Debug("Found Plugins:")
	plugins := []string{}
	for _, f := range files {
		plugins = append(plugins, path+f.Name())
		logging.Debug(" - ", f.Name())
	}
	return plugins, nil
}

// RunPluginHarvesters TODO
func RunPluginHarvesters(ctx context.Context, wg *sync.WaitGroup, targets []*types.TargetContainer, pluginDir string, resultsCh chan types.ContainerReport, client dockeradapter.Harvester) error {
	defer logging.Debug("[Plugin Harvester] Exiting")
	defer wg.Done()

	logging.Debug("[Plugin Harvester] Running")

	plugins, err := getPlugins(pluginDir)
	if err != nil {
		logging.Debug("[Plugin Harvester] Error looking for plugins in directory: %s", pluginDir)
		return nil
	}

	if len(plugins) == 0 {
		logging.Debug("[Plugin Harvester] No Docker API Capabilities found")
		return nil
	}

	pluginResultsCh := make(chan *types.ContainerReport)

	for _, target := range targets {
		go harvestPlugins(*target, client, plugins, pluginResultsCh)
	}

	for i := range targets {
		result := <-pluginResultsCh
		logging.Debug("[Plugin Harvester] Result [%d] received from name: %s, ID: %s", i, result.ContainerName, result.ContainerID)
		logging.Debug("[Plugin Harvester] Sending to collector via resultsCh")
		resultsCh <- *result
	}

	return nil
}

func harvestPlugins(target types.TargetContainer, client dockeradapter.Harvester, plugins []string, dockerAPIResultsCh chan *types.ContainerReport) {
	harvestedData := map[string]types.Capability{}

	logging.Debug("[Plugin Harvester] Harvesting %d dockerAPI capabilities", len(plugins))
	for _, plugin := range plugins {
		p, err := getPlugin(plugin)
		if err != nil {
			logging.Debug("[Plugin Harvester] Error getting plugin: %v", err)
			return
		}

		metadata := (*p).Metadata()

		ID := utils.GenerateUUID4()

		result, err := (*p).Harvest(client, ID, target)
		if err != nil {
			logging.Debug("[Plugin Harvester] error invoking Harvest on plugin: %v", err)
		}

		harvestedData[metadata.Name] = types.Capability{
			Schema:      metadata.Schema,
			Title:       metadata.ID,
			Name:        metadata.Name,
			Description: metadata.Description,
			Type:        "TODOplugin",
			HarvestID:   ID,
			Payload:     result,
			SupportedOS: metadata.SupportedOS,
		}
	}

	dockerAPIResultsCh <- GenerateContainerReport(target, harvestedData)
}

func getPlugin(path string) (*lumogonplugin.Plugin, error) {
	lib, err := plugin.Open(path)
	if err != nil {
		logging.Debug("[Plugin Harvester] Error opening plugin: %s", path)
		return nil, err
	}

	p, err := lib.Lookup("Impl")
	if err != nil {
		logging.Debug("[Plugin Harvester] Error loading plugin")
		return nil, err
	}

	fn, ok := p.(*lumogonplugin.Plugin)
	if ok != true {
		err = fmt.Errorf("[Plugin Harvester] Unable to get LumogonPlugin")
		logging.Debug("%s", err)
		return nil, err
	}

	return fn, nil
}

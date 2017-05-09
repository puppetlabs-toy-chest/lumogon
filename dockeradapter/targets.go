package dockeradapter

import (
	"archive/tar"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/types"
	"github.com/puppetlabs/lumogon/utils"
)

// NormaliseTargets TODO
func NormaliseTargets(ctx context.Context, args *[]string, client Client) ([]*types.TargetContainer, error) {
	targets := []*types.TargetContainer{}
	if len(*args) > 0 {
		targets = stringsToTargetContainers(ctx, *args, client)
	} else {
		targetContainerIDs, err := client.ContainerList(ctx)
		if err != nil {
			logging.Stderr("[Scheduler] Unable to list containers, error: %s", err)
			return nil, err
		}

		localContainerID, err := utils.GetLocalContainerID("/proc/self/cgroup")
		if err == nil {
			logging.Stderr("[Scheduler] Excluding scheduler container from harvested containers, ID: %s", localContainerID)
			targetContainerIDs = utils.RemoveStringFromSlice(targetContainerIDs, localContainerID)
		}
		targets = stringsToTargetContainers(ctx, targetContainerIDs, client)
	}
	for _, target := range targets {
		targetOS, err := getContainerOS(ctx, target.ID, client)
		if err != nil {
			logging.Stderr("[Scheduler] Error getting OS for target, ID: %s, removing from list of targets", target.ID)
			target.OSID = "unknown"
			continue
		}
		target.OSID = targetOS
	}
	return targets, nil
}

func getContainerOS(ctx context.Context, containerID string, client CopyFrom) (string, error) {
	osReleaseFile := "/etc/os-release"
	fileNotFound := regexp.MustCompile(`no such file or directory`)

	reader, _, err := client.CopyFromContainer(ctx, containerID, osReleaseFile, true)
	if err != nil {
		if fileNotFound.MatchString(err.Error()) {
			logging.Stderr("[Scheduler] File not found - assuming scratch")
			return "scratch", nil
		}
		fmt.Printf("Error reading file: %s\n", err)
		return "", err
	}
	defer reader.Close()

	var osRelease = make(map[string]string)
	tr := tar.NewReader(reader)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		logging.Stderr("reading from file: %s", hdr.Name)
		scanner := bufio.NewScanner(tr)
		for scanner.Scan() {
			line := scanner.Text()
			splitString := strings.SplitN(line, "=", 2)
			if len(splitString) != 2 {
				logging.Stderr("unable to extract key and value from line: %s", line)
				continue
			}
			key := splitString[0]
			value := splitString[1]
			if key == "" {
				logging.Stderr("ignoring empty key from line: %s", line)
				continue
			}
			osRelease[key] = value
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	if val, ok := osRelease["ID"]; ok {
		if val != "" {
			logging.Stderr("[Scheduler] detected OS Release ID: %s", val)
			return val, nil
		}
	}

	logging.Stderr("[Scheduler] unable to determine OS Release ID from file, setting to unknown")
	return "unknown", nil
}

// stringToTargetContainer converts a container ID or Name string into types.TargetContainer
func stringToTargetContainer(ctx context.Context, containerIDOrName string, client Inspector) (*types.TargetContainer, error) {
	containerJSON, err := client.ContainerInspect(ctx, containerIDOrName)
	if err != nil {
		error := fmt.Sprintf("[Scheduler] Unable to find target container: %s, error: %s", containerIDOrName, err)
		logging.Stderr(error)
		return &types.TargetContainer{}, err
	}
	targetContainer := types.TargetContainer{
		ID:   containerJSON.ContainerJSONBase.ID,
		Name: containerJSON.ContainerJSONBase.Name,
	}
	return &targetContainer, nil
}

// stringsToTargetContainers converts a slice of container ID or Name strings into a slice of types.TargetContainer
func stringsToTargetContainers(ctx context.Context, containerIDsOrNames []string, client Inspector) []*types.TargetContainer {
	targetContainers := []*types.TargetContainer{}
	for _, containerIDOrName := range containerIDsOrNames {
		targetContainer, err := stringToTargetContainer(ctx, containerIDOrName, client)
		if err != nil {
			continue
		}
		targetContainers = append(targetContainers, targetContainer)
	}
	return targetContainers
}

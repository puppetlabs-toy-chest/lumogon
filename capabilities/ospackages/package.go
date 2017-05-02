package ospackages

import (
	"context"
	"fmt"

	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/utils"
)

func runPackageCmd(client dockeradapter.Harvester, containerID string, cmd []string) (map[string]string, error) {
	attachStdout := true
	attachStderr := true

	ctx := context.Background()
	execInstance, err := client.ContainerExecCreate(ctx, containerID, cmd, attachStdout, attachStderr)
	if err != nil {
		err = fmt.Errorf("[Package] Unable to create exec: %s", err)
		logging.Stderr(err.Error())
		return nil, err
	}
	logging.Stderr("[Package] Container Exec Created, ID: %s", execInstance.ID)

	att, err := client.ContainerExecAttach(ctx, execInstance.ID, cmd, attachStdout, attachStderr)
	if err != nil {
		err = fmt.Errorf("[Package] Unable to attach exec: %s", err)
		logging.Stderr(err.Error())
		return nil, err
	}
	logging.Stderr("[Package] Container Exec Attached, ID: %s", execInstance.ID)

	defer att.Close()

	logging.Stderr("[Package] About to do exec start for package capability")
	err = client.ContainerExecStart(ctx, execInstance.ID)
	if err != nil {
		err = fmt.Errorf("[Package] Error starting exec: %s", err)
		logging.Stderr(err.Error())
		return nil, err
	}

	execInfo, err := client.ContainerExecInspect(ctx, execInstance.ID)
	if err != nil {
		err = fmt.Errorf("[Package] Error inspecting exec: %s", err)
		logging.Stderr(err.Error())
		return nil, err
	}

	if execInfo.ExitCode != 0 {
		logging.Stderr("[Package] package query returned a non-zero exit code: %d", execInfo.ExitCode)
		return nil, err
	}
	logging.Stderr("[Package] package query returned exit code: %d", execInfo.ExitCode)

	execStdout, err := dockeradapter.FilterDockerStream(att.Reader, 1)
	if err != nil {
		return nil, err
	}

	packages, err := utils.CsvToMap(execStdout)
	if err != nil {
		logging.Stderr("[Package] Error converting package csv to map")
		return nil, err
	}

	return packages, nil
}

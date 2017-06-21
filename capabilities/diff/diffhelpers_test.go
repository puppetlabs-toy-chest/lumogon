package diff

import (
	"context"

	"github.com/puppetlabs/lumogon/types"
)

func createContainerDiffFn(files []types.ChangedFile, err error) func(context.Context, string) ([]types.ChangedFile, error) {
	var containerDiffFn = func(ctx context.Context, containerID string) ([]types.ChangedFile, error) {
		if err != nil {
			return nil, err
		}
		return files, nil
	}
	return containerDiffFn
}

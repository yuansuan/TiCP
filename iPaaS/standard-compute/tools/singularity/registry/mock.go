package registry

import (
	"context"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/tools/singularity/registry/image"
)

type mockClient struct{}

func (c *mockClient) Pull(ctx context.Context, locator image.Locator, options ...PullOption) (*image.Locally, error) {
	return &image.Locally{}, nil
}

func (c *mockClient) Push(locator image.Locator, options ...PushOption) (*image.Remotely, error) {
	return &image.Remotely{}, nil
}

func (c *mockClient) Search(pattern string, options ...SearchOption) ([]*image.DefaultedLocator, error) {
	return []*image.DefaultedLocator{}, nil
}

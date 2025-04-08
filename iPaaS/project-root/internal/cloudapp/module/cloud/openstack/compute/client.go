package compute

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/startstop"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/errors"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/consts"
	_servers "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/openstack/compute/servers"
)

type Client struct {
	sdkClient *gophercloud.ServiceClient
}

func NewClient(provider *gophercloud.ProviderClient, endpointOpts gophercloud.EndpointOpts, cfg *config.OpenStack) (*Client, error) {
	sdkClient, err := openstack.NewComputeV2(provider, endpointOpts)
	if err != nil {
		return nil, fmt.Errorf("new openstack compute client failed, %w", err)
	}
	ensureEndpoint(sdkClient, cfg)
	sdkClient.Microversion = cfg.Compute.MicroVersion
	sdkClient.HTTPClient = http.Client{Timeout: 5 * time.Second}

	return &Client{
		sdkClient: sdkClient,
	}, nil
}

func ensureEndpoint(sdkClient *gophercloud.ServiceClient, cfg *config.OpenStack) {
	sdkClient.Endpoint = cfg.Compute.NovaEndpoint
}

func (c *Client) Base() *gophercloud.ServiceClient {
	return c.sdkClient
}

func (c *Client) SetMicroVersion(microVersion string) {
	c.sdkClient.Microversion = microVersion
}

func (c *Client) ListServers(listOpts servers.ListOpts) ([]_servers.Server, error) {
	pager := servers.List(c.sdkClient, listOpts)
	if pager.Err != nil {
		return nil, fmt.Errorf("list servers failed, %w", pager.Err)
	}

	pages, err := pager.AllPages()
	if err != nil {
		return nil, fmt.Errorf("list servers response to all pages failed, %w", err)
	}

	srvs := make([]_servers.Server, 0)
	err = servers.ExtractServersInto(pages, &srvs)
	if err != nil {
		return nil, fmt.Errorf("convert list servers response to []servers.Server failed, %w", err)
	}

	return srvs, nil
}

func (c *Client) GetServer(serverID string) (*_servers.Server, error) {
	getResult := servers.Get(c.sdkClient, serverID)
	if getResult.Err != nil {
		detailedErr := errors.ErrorDetails{}
		e := errors.ExtractErrorInto(getResult.Err, &detailedErr)
		if e == nil {
			_, exist := detailedErr["itemNotFound"]
			if exist {
				logging.Default().Warnf("itemNotFound: %+v", detailedErr["ItemNotFound"])
				return nil, consts.ErrInstanceNotFound
			}
		}
		return nil, fmt.Errorf("get server by ID: [%s] failed, %w", serverID, getResult.Err)
	}

	server := new(_servers.Server)
	if err := getResult.ExtractInto(server); err != nil {
		return nil, fmt.Errorf("get server by ID: [%s], result extract to *servers.Server failed, %w", serverID, err)
	}

	return server, nil
}

func (c *Client) CreateServer(createOpts servers.CreateOpts, createWithBootVolume bool) (*_servers.Server, error) {
	flavorResult := flavors.Get(c.sdkClient, createOpts.FlavorRef)
	if flavorResult.Err != nil {
		return nil, fmt.Errorf("get flavor id [%s] failed, %w", createOpts.FlavorRef, flavorResult.Err)
	}
	logging.Default().Infow("openstack createServer", "req", createOpts)

	flavor, err := flavorResult.Extract()
	if err != nil {
		return nil, fmt.Errorf("extract flavorResult to flavor failed, %w", err)
	}

	var createResult servers.CreateResult
	if !createWithBootVolume {
		createResult = servers.Create(c.sdkClient, createOpts)
		if createResult.Err != nil {
			return nil, fmt.Errorf("create server failed, %w", createResult.Err)
		}
	} else {
		blockDevices := []bootfromvolume.BlockDevice{
			{
				BootIndex:       0,
				DestinationType: bootfromvolume.DestinationVolume,
				SourceType:      bootfromvolume.SourceImage,
				UUID:            createOpts.ImageRef,
				VolumeSize:      flavor.Disk,
			},
		}

		createResult = bootfromvolume.Create(c.sdkClient, bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			BlockDevice:       blockDevices,
		})
		if createResult.Err != nil {
			return nil, fmt.Errorf("create server with boot volume failed, %w", createResult.Err)
		}
	}

	server := new(_servers.Server)
	if err := createResult.ExtractInto(server); err != nil {
		return nil, fmt.Errorf("create server result extract to *servers.Server failed, %w", err)
	}

	return server, nil
}

func (c *Client) RestoreServer(createOpts servers.CreateOpts, bootVolumeId string) (string, error) {
	flavorResult := flavors.Get(c.sdkClient, createOpts.FlavorRef)
	if flavorResult.Err != nil {
		return "", fmt.Errorf("get flavor id [%s] failed, %w", createOpts.FlavorRef, flavorResult.Err)
	}

	blockDevices := []bootfromvolume.BlockDevice{
		{
			BootIndex:       0,
			DestinationType: bootfromvolume.DestinationVolume,
			SourceType:      bootfromvolume.SourceVolume,
			UUID:            bootVolumeId,
		},
	}

	createResult := bootfromvolume.Create(c.sdkClient, bootfromvolume.CreateOptsExt{
		CreateOptsBuilder: createOpts,
		BlockDevice:       blockDevices,
	})
	if createResult.Err != nil {
		return "", fmt.Errorf("boot server from volume [%s] failed, %w", bootVolumeId, createResult.Err)
	}

	srv := new(_servers.Server)
	if err := createResult.ExtractInto(srv); err != nil {
		return "", fmt.Errorf("create server result extract to *servers.Server failed, %w", err)
	}

	return srv.ID, nil
}

func (c *Client) DeleteServer(serverID string) error {
	resp, err := c.sdkClient.Delete(c.sdkClient.ServiceURL("servers", serverID), nil)
	if err != nil {
		// 404 认为实例已经被删除
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}

		return fmt.Errorf("delete server by ID: [%s] failed, %w", serverID, err)
	}

	return nil
}

func (c *Client) StartServer(serverID string) error {
	return startstop.Start(c.sdkClient, serverID).Err
}

func (c *Client) StopServer(serverID string) error {
	return startstop.Stop(c.sdkClient, serverID).Err
}

func (c *Client) RestartServer(serverID string) error {
	return servers.Reboot(c.sdkClient, serverID, servers.RebootOpts{
		Type: servers.HardReboot,
	}).Err
}

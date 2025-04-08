package openstack

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	jsoniter "github.com/json-iterator/go"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/domain"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/openstack/compute"
	_servers "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/openstack/compute/servers"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/utils"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/util"
)

const (
	defaultCreateSessionTimeoutForAlarm = 10 * time.Minute
)

type Cloud struct {
	computeClient *compute.Client

	cfg     *config.OpenStack
	zones   *config.Zones
	metrics *utils.Metrics
}

func Initializer() func(cfg interface{}) (cloud.Cloud, error) {
	return func(cfgI interface{}) (cloud.Cloud, error) {
		cfg, ok := cfgI.(*config.OpenStack)
		if !ok {
			return nil, fmt.Errorf("config is not *config.OpenStack")
		}

		providerClient, err := newProvider(gophercloud.AuthOptions{
			AllowReauth:                 true,
			IdentityEndpoint:            cfg.Auth.IdentityEndpoint,
			ApplicationCredentialID:     cfg.Auth.CredentialID,
			ApplicationCredentialSecret: cfg.Auth.CredentialSecret,
		})
		if err != nil {
			return nil, fmt.Errorf("new provider failed, %w", err)
		}

		computeClient, err := compute.NewClient(providerClient, gophercloud.EndpointOpts{}, cfg)
		if err != nil {
			return nil, fmt.Errorf("new compute client failed, %w", err)
		}

		c := &Cloud{
			computeClient: computeClient,
			cfg:           cfg,
			metrics:       utils.NewMetrics("openstack"),
		}

		return c, nil
	}
}

func (c *Cloud) RunInstance(software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error) {
	c.metrics.CallAPI("RunInstance")

	createServerOpts := servers.CreateOpts{
		Name:        fmt.Sprintf("CloudApp_%s", session.Id.String()),
		ImageRef:    software.ImageId,
		FlavorRef:   hardware.InstanceType,
		UserData:    []byte(base64.StdEncoding.EncodeToString([]byte(instance.UserScript))),
		AdminPass:   instance.SshPassword,
		ConfigDrive: utils.BoolPtr(true),
		Networks: []servers.Network{
			{
				UUID: c.cfg.Network.Uuid,
			},
		},
		Tags: c.cfg.Tags,
	}

	server, err := c.computeClient.CreateServer(createServerOpts, c.cfg.CreateWithBootVolume)
	if err != nil {
		return "", fmt.Errorf("create server failed, %w", err)
	}

	return server.ID, nil
}

func (c *Cloud) TerminateInstance(instanceID string) error {
	return c.computeClient.DeleteServer(instanceID)
}

func (c *Cloud) DescribeInstance(instanceID string) (domain.Instance, error) {
	server, err := c.computeClient.GetServer(instanceID)
	if err != nil {
		return nil, fmt.Errorf("get server failed, %w", err)
	}
	c.checkIfErrorInstanceExistsAndAlarm([]_servers.Server{*server})

	return newInstance(server), nil
}

func (c *Cloud) DescribeInstances() ([]domain.Instance, error) {
	srvs, err := c.computeClient.ListServers(servers.ListOpts{
		Tags: stringSliceToString(c.cfg.Tags),
	})
	if err != nil {
		return nil, fmt.Errorf("list servers failed, %w", err)
	}
	c.checkIfErrorInstanceExistsAndAlarm(srvs)

	instances := make([]domain.Instance, 0, len(srvs))
	for i := range srvs {
		instances = append(instances, newInstance(&srvs[i]))
	}

	return instances, nil
}

type port struct {
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
	Type    string `json:"OS-EXT-IPS:type"`
	Addr    string `json:"addr"`
	Version int    `json:"version"`
}

// ParseRawInstanceIP 解析raw中服务器ip
func (c *Cloud) ParseRawInstanceIP(raw string) (string, error) {
	rawData := new(_servers.Server)

	if err := jsoniter.UnmarshalFromString(raw, rawData); err != nil {
		return "", fmt.Errorf("unmarshal from string failed, %w", err)
	}

	portsI, ok := rawData.Addresses[c.cfg.Network.Name]
	if !ok {
		return "", fmt.Errorf("cannot find ports which name = %s", c.cfg.Network.Name)
	}

	data, err := jsoniter.Marshal(portsI)
	if err != nil {
		return "", fmt.Errorf("ports json marshal failed, %w", err)
	}

	ports := make([]port, 0)
	if err = jsoniter.Unmarshal(data, &ports); err != nil {
		return "", fmt.Errorf("unmarshal to ports failed, %w", err)
	}

	if len(ports) == 0 {
		return "", fmt.Errorf("port not found")
	}

	return ports[0].Addr, nil
}

func (c *Cloud) StartInstance(instanceId string) error {
	return c.computeClient.StartServer(instanceId)
}

func (c *Cloud) StopInstance(instanceId string) error {
	return c.computeClient.StopServer(instanceId)
}

func (c *Cloud) RestartInstance(instanceId string) error {
	return c.computeClient.RestartServer(instanceId)
}

func (c *Cloud) ParseBootVolumeId(raw string) (string, error) {
	rawData := new(_servers.Server)

	if err := jsoniter.UnmarshalFromString(raw, rawData); err != nil {
		return "", fmt.Errorf("unmarshal from string failed, %w", err)
	}

	if len(rawData.AttachedVolumes) == 0 {
		// no attached volume, not return error here
		return "", nil
	}

	// only attach one volume for now so just return the first volume id.
	return rawData.AttachedVolumes[0].ID, nil
}

func (c *Cloud) RestoreInstance(software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error) {
	c.metrics.CallAPI("RestoreInstance")

	createServerOpts := servers.CreateOpts{
		Name:        fmt.Sprintf("CloudApp_%s", session.Id.String()),
		ImageRef:    software.ImageId,
		FlavorRef:   hardware.InstanceType,
		UserData:    []byte(base64.StdEncoding.EncodeToString([]byte(instance.UserScript))),
		ConfigDrive: utils.BoolPtr(true),
		Networks: []servers.Network{
			{
				UUID: c.cfg.Network.Uuid,
			},
		},
		Tags: c.cfg.Tags,
	}

	return c.computeClient.RestoreServer(createServerOpts, instance.BootVolumeId)
}

func (c *Cloud) checkIfErrorInstanceExistsAndAlarm(srvs []_servers.Server) {
	for _, srv := range srvs {
		if strings.ToUpper(srv.Status) == "ERROR" {
			logging.Default().Errorf("instance status is ERROR, id = [%s]", srv.ID)
		} else {
			if shouldAlarmIfCreateTimeout(srv.Created, srv.VmState) {
				logging.Default().Errorf("instance status is not [ACTIVE] or [SHUTOFF] or [STOPPED] or [DELETED], id = [%s], status = [%s], created = [%s], currentTime: [%s], timeoutDuration: [%s]",
					srv.ID, srv.VmState, srv.Created, time.Now(), config.GetConfig().CreateSessionTimeoutForAlarm)
			}
		}
	}
}

// 与当前时间做比较，超过阈值时间，vmState为非ACTIVE/SHUTOFF/STOPPED/DELETED的，需要报警
func shouldAlarmIfCreateTimeout(created time.Time, vmState string) bool {
	createSessionTimeoutForAlarm := config.GetConfig().CreateSessionTimeoutForAlarm
	if createSessionTimeoutForAlarm == 0 {
		createSessionTimeoutForAlarm = defaultCreateSessionTimeoutForAlarm
	}

	return !util.EmptyTime(created) &&
		created.Add(config.GetConfig().CreateSessionTimeoutForAlarm).Before(time.Now()) &&
		!isVmStateACTIVEorSHUTOFForSTOPPEDorDELETED(vmState)
}

func isVmStateACTIVEorSHUTOFForSTOPPEDorDELETED(vmState string) bool {
	vmState = strings.ToUpper(vmState)
	// ACTIVE: VM is running with the specified image.
	// STOPPED: VM is not running, and the image is on disk.
	// ref: https://wiki.openstack.org/wiki/VMState
	return vmState == "ACTIVE" || vmState == "SHUTOFF" || vmState == "STOPPED" || vmState == "DELETED"
}

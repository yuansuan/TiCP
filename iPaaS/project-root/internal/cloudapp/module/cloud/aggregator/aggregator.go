package aggregator

import (
	"errors"
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/domain"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/openstack"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
)

// ErrBadAvailableZone 无效的可用区
var ErrBadAvailableZone = errors.New("bad available zone")

type Aggregator struct {
	zones             *config.Zones
	m                 map[zone.Zone]cloud.Cloud
	log               *logging.Logger
	cloudInitializers map[config.CloudType]func(cloudCfg interface{}) (cloud.Cloud, error)
}

func New() (cloud.Aggregator, error) {
	zones := config.GetConfig().CloudApp.Zones
	if zones == nil {
		return nil, fmt.Errorf("config.Zones cannot be nil")
	}

	a := &Aggregator{
		zones: &zones,
		m:     make(map[zone.Zone]cloud.Cloud),
		log:   logging.Default(),
	}
	a.initCloudInitializers()

	for zoneName, opt := range zones {
		c, err := a.determineCloud(opt)
		if err != nil {
			return nil, fmt.Errorf("determine cloud failed, %w", err)
		}

		a.Register(zoneName, c)
	}

	return a, nil
}

// RunInstance 根据会话信息创建相应的虚拟机实例
func (a *Aggregator) RunInstance(zone zone.Zone, software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error) {
	c, err := a.Select(zone)
	if err != nil {
		return "", err
	}

	return c.RunInstance(software, instance, session, hardware)
}

// StopInstance 停止并销毁实例Id对应的实例
func (a *Aggregator) TerminateInstance(zone zone.Zone, instanceID string) error {
	c, err := a.Select(zone)
	if err != nil {
		return err
	}

	return c.TerminateInstance(instanceID)
}

// DescribeInstance 查询指定实例Id的实例信息
func (a *Aggregator) DescribeInstance(zone zone.Zone, instanceID string) (domain.Instance, error) {
	c, err := a.Select(zone)
	if err != nil {
		return nil, err
	}

	return c.DescribeInstance(instanceID)
}

func (a *Aggregator) ParseRawInstanceIP(zone zone.Zone, raw string) (string, error) {
	c, err := a.Select(zone)
	if err != nil {
		return "", err
	}
	return c.ParseRawInstanceIP(raw)
}

// DescribeInstances 列出所有当前集群正常运行的实例列表
func (a *Aggregator) DescribeInstances() ([]domain.Instance, error) {
	var instances []domain.Instance
	for name, c := range a.m {
		resp, err := c.DescribeInstances()
		if err != nil {
			a.log.Warnw("describe cloud instances failed", "name", name, "error", err)

			// 如果是到 openstack 的网络有问题，先不处理这个集群里的实例
			if strings.Contains(err.Error(), "i/o timeout") ||
				strings.Contains(err.Error(), "connection reset by peer") {
				continue
			}
		}

		instances = append(instances, resp...)
	}

	logging.Default().Debugw("cloud instances", "instances", instances)
	return instances, nil
}

func (a *Aggregator) GetZoneOpts(zone zone.Zone) (*config.ZoneOption, error) {
	zoneOpts, exist := (*a.zones)[zone]
	if !exist {
		return nil, fmt.Errorf("zone [%s] not support", zone)
	}

	return zoneOpts, nil
}

func (a *Aggregator) StartInstance(zone zone.Zone, instanceId string) error {
	c, err := a.Select(zone)
	if err != nil {
		return err
	}

	return c.StartInstance(instanceId)
}

func (a *Aggregator) StopInstance(zone zone.Zone, instanceId string) error {
	c, err := a.Select(zone)
	if err != nil {
		return err
	}

	return c.StopInstance(instanceId)
}

func (a *Aggregator) RestartInstance(zone zone.Zone, instanceId string) error {
	c, err := a.Select(zone)
	if err != nil {
		return err
	}

	return c.RestartInstance(instanceId)
}

// ParseBootVolumeId only openstack support for now
func (a *Aggregator) ParseBootVolumeId(zone zone.Zone, raw string) (string, error) {
	c, err := a.Select(zone)
	if err != nil {
		return "", err
	}

	return c.ParseBootVolumeId(raw)
}

// RestoreInstance only openstack support for now
func (a *Aggregator) RestoreInstance(zone zone.Zone, software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error) {
	c, err := a.Select(zone)
	if err != nil {
		return "", err
	}

	return c.RestoreInstance(software, instance, session, hardware)
}

func (a *Aggregator) Register(z zone.Zone, c cloud.Cloud) {
	a.m[z] = c
}

func (a *Aggregator) Select(z zone.Zone) (cloud.Cloud, error) {
	if !z.IsValid() || z.IsEmpty() {
		return nil, ErrBadAvailableZone
	}

	c, ok := a.m[z]
	if !ok {
		return nil, ErrBadAvailableZone
	}

	return c, nil
}

func (a *Aggregator) initCloudInitializers() {
	a.cloudInitializers = make(map[config.CloudType]func(cloudCfg interface{}) (cloud.Cloud, error))
	a.cloudInitializers[config.OpenstackCloudType] = openstack.Initializer()
}

func (a *Aggregator) determineCloud(zoneOpts *config.ZoneOption) (cloud.Cloud, error) {
	c, exist := a.cloudInitializers[zoneOpts.Cloud]
	if !exist {
		return nil, fmt.Errorf("third cloud name: %s not implement yet", zoneOpts.Cloud)
	}

	cloudCfg, err := a.determineCloudCfg(zoneOpts)
	if err != nil {
		return nil, fmt.Errorf("determine cloud config failed, %w", err)
	}

	return c(cloudCfg)
}

func (a *Aggregator) determineCloudCfg(zoneOpts *config.ZoneOption) (interface{}, error) {
	switch zoneOpts.Cloud {
	case config.TencentCloudType:
		return zoneOpts.Tencent, nil
	case config.ShanheCloudType:
		return zoneOpts.Shanhe, nil
	case config.OpenstackCloudType:
		return zoneOpts.OpenStack, nil
	default:
		return nil, fmt.Errorf("unknown cloud type: %s", zoneOpts.Cloud)
	}
}

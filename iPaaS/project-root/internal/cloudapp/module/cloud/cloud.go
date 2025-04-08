package cloud

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	zone "github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/cloud/domain"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/cloudapp/module/dao/models"
)

type Aggregator interface {
	RunInstance(zone zone.Zone, software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error)

	TerminateInstance(zone zone.Zone, instanceID string) error

	DescribeInstance(zone zone.Zone, instanceID string) (domain.Instance, error)

	DescribeInstances() ([]domain.Instance, error)

	ParseRawInstanceIP(zone zone.Zone, raw string) (string, error)

	GetZoneOpts(zone zone.Zone) (*config.ZoneOption, error)

	StartInstance(zone zone.Zone, instanceId string) error

	StopInstance(zone zone.Zone, instanceId string) error

	RestartInstance(zone zone.Zone, instanceId string) error

	ParseBootVolumeId(zone zone.Zone, raw string) (string, error)

	RestoreInstance(zone zone.Zone, software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error)
}

// Cloud 第三方云服务接口
type Cloud interface {
	RunInstance(software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error)

	TerminateInstance(instanceID string) error

	DescribeInstance(instanceID string) (domain.Instance, error)

	DescribeInstances() ([]domain.Instance, error)

	ParseRawInstanceIP(raw string) (string, error)

	StartInstance(instanceId string) error

	StopInstance(instanceId string) error

	RestartInstance(instanceId string) error

	ParseBootVolumeId(raw string) (string, error)

	RestoreInstance(software *models.Software, instance *models.Instance, session *models.Session, hardware *models.Hardware) (string, error)
}

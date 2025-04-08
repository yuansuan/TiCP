package delete

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"

type Option func(req *hardware.AdminDeleteRequest) error

func (api API) Id(id string) Option {
	return func(req *hardware.AdminDeleteRequest) error {
		req.HardwareId = &id
		return nil
	}
}

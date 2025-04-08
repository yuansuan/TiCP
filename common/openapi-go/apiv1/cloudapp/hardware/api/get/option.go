package get

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
)

type Option func(req *hardware.APIGetRequest) error

func (api API) Id(id string) Option {
	return func(req *hardware.APIGetRequest) error {
		req.HardwareId = &id
		return nil
	}
}

package delete

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"

type Option func(req *software.AdminDeleteRequest) error

func (api API) Id(id string) Option {
	return func(req *software.AdminDeleteRequest) error {
		req.SoftwareId = &id
		return nil
	}
}

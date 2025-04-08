package add

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
)

type Option func(req *software.AdminPostRequest) error

func (api API) Name(name string) Option {
	return func(req *software.AdminPostRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *software.AdminPostRequest) error {
		req.Desc = &desc
		return nil
	}
}

func (api API) Icon(icon string) Option {
	return func(req *software.AdminPostRequest) error {
		req.Icon = &icon
		return nil
	}
}

func (api API) Platform(platform string) Option {
	return func(req *software.AdminPostRequest) error {
		req.Platform = &platform
		return nil
	}
}

func (api API) ImageId(imageId string) Option {
	return func(req *software.AdminPostRequest) error {
		req.ImageId = &imageId
		return nil
	}
}

func (api API) InitScript(initScript string) Option {
	return func(req *software.AdminPostRequest) error {
		req.InitScript = &initScript
		return nil
	}
}

func (api API) GpuDesired(gpuDesired bool) Option {
	return func(req *software.AdminPostRequest) error {
		req.GpuDesired = &gpuDesired
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *software.AdminPostRequest) error {
		req.Zone = &zone
		return nil
	}
}

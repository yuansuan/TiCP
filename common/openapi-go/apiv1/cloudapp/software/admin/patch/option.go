package put

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
)

type Option func(req *software.AdminPatchRequest) error

func (api API) Id(id string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.SoftwareId = &id
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Desc(desc string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.Desc = &desc
		return nil
	}
}

func (api API) Icon(icon string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.Icon = &icon
		return nil
	}
}

func (api API) Platform(platform string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.Platform = &platform
		return nil
	}
}

func (api API) ImageId(imageId string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.ImageId = &imageId
		return nil
	}
}

func (api API) InitScript(initScript string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.InitScript = &initScript
		return nil
	}
}

func (api API) GpuDesired(gpuDesired bool) Option {
	return func(req *software.AdminPatchRequest) error {
		req.GpuDesired = &gpuDesired
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *software.AdminPatchRequest) error {
		req.Zone = &zone
		return nil
	}
}

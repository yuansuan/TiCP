package list

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
)

type Option func(req *software.AdminListRequest) error

func (api API) UserId(userId string) Option {
	return func(req *software.AdminListRequest) error {
		req.UserId = &userId
		return nil
	}
}

func (api API) Name(name string) Option {
	return func(req *software.AdminListRequest) error {
		req.Name = &name
		return nil
	}
}

func (api API) Platform(platform string) Option {
	return func(req *software.AdminListRequest) error {
		req.Platform = &platform
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *software.AdminListRequest) error {
		req.Zone = &zone
		return nil
	}
}

func (api API) PageOffset(pageOffset int) Option {
	return func(req *software.AdminListRequest) error {
		req.PageOffset = &pageOffset
		return nil
	}
}

func (api API) PageSize(pageSize int) Option {
	return func(req *software.AdminListRequest) error {
		req.PageSize = &pageSize
		return nil
	}
}

package list

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
)

type Option func(req *session.AdminListRequest) error

func (api API) SessionIds(sessionIds string) Option {
	return func(req *session.AdminListRequest) error {
		req.SessionIds = &sessionIds
		return nil
	}
}

func (api API) Status(status string) Option {
	return func(req *session.AdminListRequest) error {
		req.Status = &status
		return nil
	}
}

func (api API) Zone(zone string) Option {
	return func(req *session.AdminListRequest) error {
		req.Zone = &zone
		return nil
	}
}

func (api API) PageOffset(pageOffset int) Option {
	return func(req *session.AdminListRequest) error {
		req.PageOffset = &pageOffset
		return nil
	}
}

func (api API) PageSize(pageSize int) Option {
	return func(req *session.AdminListRequest) error {
		req.PageSize = &pageSize
		return nil
	}
}

func (api API) UserIds(userIds string) Option {
	return func(req *session.AdminListRequest) error {
		req.UserIds = &userIds
		return nil
	}
}

func (api API) WithDeleted(withDeleted bool) Option {
	return func(req *session.AdminListRequest) error {
		req.WithDeleted = withDeleted
		return nil
	}
}

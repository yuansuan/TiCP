package mount

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.UmountRequest) error

func (api API) SessionId(sessionId string) Option {
	return func(req *session.UmountRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

func (api API) MountPoint(mountPoint string) Option {
	return func(req *session.UmountRequest) error {
		req.MountPoint = &mountPoint
		return nil
	}
}

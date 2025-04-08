package mount

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.MountRequest) error

func (api API) SessionId(sessionId string) Option {
	return func(req *session.MountRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

func (api API) ShareDirectory(shareDirectory string) Option {
	return func(req *session.MountRequest) error {
		req.ShareDirectory = &shareDirectory
		return nil
	}
}

func (api API) MountPoint(mountPoint string) Option {
	return func(req *session.MountRequest) error {
		req.MountPoint = &mountPoint
		return nil
	}
}

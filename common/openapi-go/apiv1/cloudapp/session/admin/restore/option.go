package restore

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.AdminRestoreRequest)

func (api API) SessionId(sessionId string) Option {
	return func(req *session.AdminRestoreRequest) {
		req.SessionId = &sessionId
	}
}

func (api API) UserId(userId string) Option {
	return func(req *session.AdminRestoreRequest) {
		req.UserId = &userId
	}
}

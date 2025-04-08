package restore

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.ApiRestoreRequest)

func (api API) SessionId(sessionId string) Option {
	return func(req *session.ApiRestoreRequest) {
		req.SessionId = &sessionId
	}
}

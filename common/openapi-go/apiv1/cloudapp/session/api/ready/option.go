package ready

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.ApiReadyRequest) error

func (api API) Id(sessionId string) Option {
	return func(req *session.ApiReadyRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

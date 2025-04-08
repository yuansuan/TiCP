package close

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.ApiCloseRequest) error

func (api API) Id(sessionId string) Option {
	return func(req *session.ApiCloseRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

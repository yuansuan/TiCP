package close

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"

type Option func(req *session.AdminCloseRequest) error

func (api API) Id(sessionId string) Option {
	return func(req *session.AdminCloseRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

func (api API) Reason(reason string) Option {
	return func(req *session.AdminCloseRequest) error {
		req.Reason = &reason
		return nil
	}
}

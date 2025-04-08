package poweroff

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
)

type Option func(req *session.PowerOffRequest)

func (api API) Id(sessionId string) Option {
	return func(req *session.PowerOffRequest) {
		req.SessionId = &sessionId
	}
}

package reboot

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/session"
)

type Option func(req *session.RebootRequest)

func (api API) Id(sessionId string) Option {
	return func(req *session.RebootRequest) {
		req.SessionId = &sessionId
	}
}

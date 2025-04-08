package get_url

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"

type Option func(req *remoteapp.ApiGetRequest) error

func (api API) SessionId(sessionId string) Option {
	return func(req *remoteapp.ApiGetRequest) error {
		req.SessionId = &sessionId
		return nil
	}
}

func (api API) RemoteAppName(remoteAppName string) Option {
	return func(req *remoteapp.ApiGetRequest) error {
		req.RemoteAppName = &remoteAppName
		return nil
	}
}

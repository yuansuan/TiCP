package delete

import "github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/remoteapp"

type Option func(req *remoteapp.AdminDeleteRequest) error

func (api API) Id(remoteAppId string) Option {
	return func(req *remoteapp.AdminDeleteRequest) error {
		req.RemoteAppId = &remoteAppId
		return nil
	}
}

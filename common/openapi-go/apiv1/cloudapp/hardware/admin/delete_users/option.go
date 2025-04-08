package delete_users

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/hardware"
)

type Option func(req *hardware.AdminDeleteUsersRequest)

func (api API) Users(users []string) Option {
	return func(req *hardware.AdminDeleteUsersRequest) {
		req.Users = users
	}
}

func (api API) Hardwares(hardwares []string) Option {
	return func(req *hardware.AdminDeleteUsersRequest) {
		req.Hardwares = hardwares
	}
}

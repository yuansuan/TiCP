package add_users

import (
	"github.com/yuansuan/ticp/common/project-root-api/cloud_app/v1/software"
)

type Option func(req *software.AdminPostUsersRequest)

func (api API) Users(users []string) Option {
	return func(req *software.AdminPostUsersRequest) {
		req.Users = users
	}
}

func (api API) Softwares(softwares []string) Option {
	return func(req *software.AdminPostUsersRequest) {
		req.Softwares = softwares
	}
}

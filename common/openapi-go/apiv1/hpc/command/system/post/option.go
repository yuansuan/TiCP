package post

import (
	"github.com/yuansuan/ticp/common/project-root-api/hpc/v1/command"
)

type Option func(req *command.SystemPostRequest)

func (api API) Command(cmd string) Option {
	return func(req *command.SystemPostRequest) {
		req.Command = cmd
	}
}

// Timeout 单位：秒
func (api API) Timeout(timeout int) Option {
	return func(req *command.SystemPostRequest) {
		req.Timeout = timeout
	}
}

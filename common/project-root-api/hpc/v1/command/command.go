package command

import (
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type SystemPostRequest struct {
	Command string `json:"Command"`
	// Timeout 单位：秒
	Timeout int `json:"Timeout"`
}

type SystemPostResponse struct {
	v20230530.Response

	Data *SystemPostResponseData `json:"Data,omitempty"`
}

type SystemPostResponseData struct {
	Stdout    string `json:"Stdout,omitempty"`
	Stderr    string `json:"Stderr,omitempty"`
	ExitCode  int    `json:"ExitCode,omitempty"`
	IsTimeout bool   `json:"IsTimeout,omitempty"`
}

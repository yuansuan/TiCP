package execscript

import (
	"github.com/yuansuan/ticp/common/project-root-api/rdpgo"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	rdpgo.BaseRequest
	ScriptRunner         string `json:"ScriptRunner"`
	ScriptContentEncoded string `json:"ScriptContentEncoded"`
	WaitTillEnd          bool   `json:"WaitTillEnd"`
}

type Response struct {
	v20230530.Response

	Data *ResponseData `json:"Data,omitempty"`
}

type ResponseData struct {
	ExitCode int    `json:"ExitCode"`
	Stdout   string `json:"Stdout"`
	Stderr   string `json:"Stderr"`
}

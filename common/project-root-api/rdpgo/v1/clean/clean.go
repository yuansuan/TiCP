package clean

import (
	"github.com/yuansuan/ticp/common/project-root-api/rdpgo"
	v20230530 "github.com/yuansuan/ticp/common/project-root-api/schema/v20230530"
)

type Request struct {
	rdpgo.BaseRequest
}

type Response struct {
	v20230530.Response
}

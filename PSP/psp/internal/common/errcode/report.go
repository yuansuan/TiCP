package errcode

import (
	"google.golang.org/grpc/codes"
)

// Report service error codes from 19001 ~ 20000
const (
	ErrCommonReportGetFailed       codes.Code = 19001
	ErrCommonReportExportFailed    codes.Code = 19002
	ErrCommonReportTypeUnsupported codes.Code = 19003
)
const (
	MsgCommonReportFailGet         = "获取报表信息失败"
	MsgCommonReportFailExport      = "导出报表信息失败"
	MsgCommonReportTypeUnsupported = "报表类型不支持"
)

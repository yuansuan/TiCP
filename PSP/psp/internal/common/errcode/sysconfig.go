package errcode

import (
	"google.golang.org/grpc/codes"
)

// SysConfig 错误码范围: 20001 ~ 21000
// 命名格式：Err + 服务名 + 具体错误
const (
	ErrSysConfigGetGlobalFailed   codes.Code = 20001
	ErrSysConfigGetJobBurstFailed codes.Code = 20002
	ErrSysConfigSetJobBurstFailed codes.Code = 20003
	ErrSysConfigGetJobFailed      codes.Code = 20004
	ErrSysConfigSetJobFailed      codes.Code = 20005

	ErrSysConfigSetEmailFailed codes.Code = 20006
	ErrSysConfigGetEmailFailed codes.Code = 20007

	ErrNoticeTestSendEmailFailed codes.Code = 20008
	ErrNoticeGetEmailFailed      codes.Code = 20009
	ErrNoticeSetEmailFailed      codes.Code = 20010

	ErrSysConfigGetThreePersonFailed codes.Code = 20011
	ErrSysConfigSetThreePersonFailed codes.Code = 20012
)

// SysConfigCodeMsg ...
var SysConfigCodeMsg = map[codes.Code]string{
	ErrSysConfigGetGlobalFailed:   "获取全局配置失败",
	ErrSysConfigGetJobBurstFailed: "获取作业爆发配置失败",
	ErrSysConfigSetJobBurstFailed: "设置作业爆发配置失败",
	ErrSysConfigGetJobFailed:      "获取作业配置失败",
	ErrSysConfigSetJobFailed:      "设置作业配置失败",
	ErrSysConfigSetEmailFailed:    "设置邮件通知配置失败",
	ErrSysConfigGetEmailFailed:    "获取邮件通知配置失败",

	ErrNoticeTestSendEmailFailed:     "发送电子邮件测试失败",
	ErrNoticeGetEmailFailed:          "获取电子邮件信息失败",
	ErrNoticeSetEmailFailed:          "设置电子邮件信息失败",
	ErrSysConfigGetThreePersonFailed: "获取三员管理配置失败",
	ErrSysConfigSetThreePersonFailed: "设置三员管理配置失败",
}

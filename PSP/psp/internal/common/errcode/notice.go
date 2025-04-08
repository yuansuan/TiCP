package errcode

import (
	"google.golang.org/grpc/codes"
)

// Notice service error codes from 13001 to 14000
const (
	ErrNoticeFailSend                  codes.Code = 13001
	ErrNoticeFailList                  codes.Code = 13002
	ErrNoticeFailRead                  codes.Code = 13003
	ErrNoticeSendEmailFailed           codes.Code = 13006
	ErrNoticeSendEmailNotEnable        codes.Code = 13007
	ErrNoticeNotSettingEmailTemplate   codes.Code = 13008
	ErrNoticeEmailTemplateTypeNotMatch codes.Code = 13009
)

// Notice service error message
const (
	MsgNoticeFailSend = "发送消息失败"
	MsgNoticeFailList = "获取消息列表失败"
	MsgNoticeFailRead = "设置消息已读失败"
)

// NoticeCodeMsg ...
var NoticeCodeMsg = map[codes.Code]string{
	ErrNoticeFailSend:                  "发送消息失败",
	ErrNoticeFailList:                  "获取消息列表失败",
	ErrNoticeFailRead:                  "设置消息已读失败",
	ErrNoticeSendEmailFailed:           "发送电子邮件失败",
	ErrNoticeSendEmailNotEnable:        "未开启发送电子邮件功能",
	ErrNoticeNotSettingEmailTemplate:   "没有配置电子邮件模版",
	ErrNoticeEmailTemplateTypeNotMatch: "邮件模版类型不匹配",
}

package dto

import "github.com/yuansuan/ticp/PSP/psp/pkg/xtype"

type OnlineUserListResponse struct {
	List []*OnlineUserResponse `json:"list"`
	Page *xtype.PageResp       `json:"page"`
}

type OnlineUserResponse struct {
	Name  string `json:"name"`  // 用户名
	Count int64  `json:"count"` // 会话数
}

type OnlineUserInfoListResponse struct {
	List []*OnlineUserInfoResponse `json:"list"`
	Page *xtype.PageResp           `json:"page"`
}

type OnlineUserInfoResponse struct {
	Jti        string `json:"jti"`         // 会话id
	ExpireTime string `json:"expire_time"` // 过期时间
	IP         string `json:"ip"`          // IP地址
}

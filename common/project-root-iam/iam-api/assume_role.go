package iam_api

import "time"

type AssumeRoleRequest struct {
	// Yrn 代表 yuansuan resource name
	// yrn:partition:service:region:account-id:resource-id
	// 例如 yrn:ys:iam::123456:role/CloudComputeRole
	// 例如 yrn:ys:iam::YSCSP:role/VipUserRole
	RoleYrn string `json:"RoleYrn"`
	// RoleSessionName 可以帮助你区分或识别创建会话的不同实体或应用程序。
	// 例如，你可能有多个独立的服务或应用程序，它们都以同一个 IAM 角色的身份进行操作，但每个服务或应用程序在调用 AssumeRole 时会使用不同的 RoleSessionName。
	// 这样，在查看日志时，你就可以知道是哪个服务或应用程序执行的特定操作。
	RoleSessionName string `json:"RoleSessionName"`
	DurationSeconds int    `json:"DurationSeconds"`
}

type Credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SessionToken    string `json:"SessionToken"`
}

type AssumeRoleResponse struct {
	Credentials *Credentials `json:"Credentials"`
	ExpireTime  time.Time    `json:"ExpireTime"`
}

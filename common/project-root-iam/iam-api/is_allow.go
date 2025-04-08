package iam_api

type IsAllowRequest struct {
	Action   string `json:"Action"`
	Resource string `json:"Resource"`
	// 填请求者的AccessKeyId
	Subject string `json:"Subject"`
}

type IsAllowResponse struct {
	Allow   bool   `json:"Allow"`
	Message string `json:"Message"`
}

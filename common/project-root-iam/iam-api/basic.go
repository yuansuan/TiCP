package iam_api

type BasicResponse struct {
	RequestId    string      `json:"RequestId"`
	ErrorCode    string      `json:"ErrorCode"`
	ErrorMessage string      `json:"ErrorMessage"`
	Data         interface{} `json:"Data"`
}

package iam_api

import "time"

type AddSecretResponse GetSecretResponse
type DeleteSecretRequest struct {
	AccessKeyId string `json:"AccessKeyId"`
}

type DeleteSecretResponse struct {
}

type ListSecretRequest struct {
	UserID string `json:"UserID"`
}

type ListAllSecretResponse struct {
	Secrets []*CacheSecret `json:"Secrets"`
}

type ListSecretResponse []*GetSecretResponse

type CacheSecret struct {
	AccessKeyId     string    `json:"AccessKeyId"`
	AccessKeySecret string    `json:"AccessKeySecret"`
	YSId            string    `json:"YSId"`
	Expire          time.Time `json:"Expire"`
	Tag             string    `json:"Tag"`
	SessionToken    string    `json:"SessionToken"`
}

type GetSecretRequest struct {
	AccessKeyId string `json:"AccessKeyId"`
}

type GetSecretResponse struct {
	AccessKeyId     string    `json:"AccessKeyId"`
	AccessKeySecret string    `json:"AccessKeySecret"`
	YSId            string    `json:"YSId"`
	Expire          time.Time `json:"Expire"`
}

type IsYSProductAccountRequest struct {
	UserId string `json:"UserId"`
}

type IsYSProductAccountResponse struct {
	IsYSProductAccount bool `json:"IsYSProductAccount"`
}

type AdminAddSecretRequest struct {
	UserId string `json:"UserId"`
	Tag    string `json:"Tag"`
}
type AdminGetSecretRequest GetSecretRequest
type AdminAddSecretResponse GetSecretResponse

type AdminGetSecretResponse struct {
	GetSecretResponse
	Tag string `json:"Tag"`
}

type AdminListSecretRequest struct {
	UserId string `json:"UserId"`
}

type AdminListSecretsRequest struct {
	PageOffset int64 `json:"PageOffset"`
	PageSize   int64 `json:"PageSize"`
}

type AdminListSecretResponse []*AdminGetSecretResponse

type AdminDeleteSecretRequest struct {
	AccessKeyId string `json:"AccessKeyId"`
}

type AdminDeleteSecretResponse struct {
}

type AdminUpdateTagRequest struct {
	AccessKeyId string `json:"AccessKeyId"`
	Tag         string `json:"Tag"`
	UserId      string `json:"UserId"`
}

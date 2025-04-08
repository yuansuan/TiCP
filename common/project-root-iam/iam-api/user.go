package iam_api

import "time"

type AdminGetUserRequest struct {
	UserId string `json:"UserId"`
}

type AdminGetUserResponse struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}

type AdminAddUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type AdminAddUserResponse struct {
	UserId string `json:"userId"`
}

type AdminUpdateUserRequest struct {
	Name   string `json:"name"`
	UserId string `json:"userId"`
}

type AdminListUserByNameRequest struct {
	PageOffset int64  `json:"PageOffset"`
	PageSize   int64  `json:"PageSize"`
	Name       string `json:"Name"`
}

type AdminListUserByNameResponse struct {
	Users []*User `json:"users"`
	Total int64   `json:"total"`
}

type User struct {
	Ysid            string    `json:"ysid"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	RealName        string    `json:"real_name"`
	UserName        string    `json:"user_name"`
	DisplayUserName string    `json:"display_user_name"`
	UserChannel     string    `json:"user_channel"`
	UserSource      string    `json:"user_source"`
	UserReferer     string    `json:"user_referer"`
	CreateTime      time.Time `json:"create_time"`
}

type AddAccountRequest struct {
	Phone                   string `json:"Phone"`
	Name                    string `json:"Name"`
	CompanyName             string `json:"CompanyName"`
	Password                string `json:"Password"`
	UnifiedSocialCreditCode string `json:"UnifiedSocialCreditCode"`
	UserChannel             string `json:"UserChannel"`
	Email                   string `json:"Email"`
}

type AddAccountResponse struct {
	Phone                   string `json:"Phone"`
	Name                    string `json:"Name"`
	CompanyName             string `json:"CompanyName"`
	UserChannel             string `json:"UserChannel"`
	Password                string `json:"Password"`
	YsId                    string `json:"YsId"`
	AccessKeyId             string `json:"AccessKeyId"`
	AccessKeySecret         string `json:"AccessKeySecret"`
	UnifiedSocialCreditCode string `json:"UnifiedSocialCreditCode"`
	Email                   string `json:"Email"`
}

type ExchangeCredentialsRequest struct {
	YsId     string `json:"YsId"`
	Phone    string `json:"Phone"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type ExchangeCredentialsResponse AddAccountResponse

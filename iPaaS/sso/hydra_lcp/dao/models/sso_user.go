package models

import (
	"time"
)

type SsoUser struct {
	Ysid           int64     `xorm:"not null pk BIGINT(20)"`
	Name           string    `xorm:"not null VARCHAR(128)"`
	RealName       string    `xorm:"not null VARCHAR(128)"`
	PwdHash        string    `xorm:"VARCHAR(128)"`
	Email          string    `xorm:"unique VARCHAR(128)"`
	Company        string    `xorm:"not null default '' VARCHAR(255)"`
	Phone          string    `xorm:"unique VARCHAR(16)"`
	WechatUnionId  string    `xorm:"unique VARCHAR(28)"`
	WechatOpenId   string    `xorm:"not null VARCHAR(128)"`
	WechatNickName string    `xorm:"not null VARCHAR(128)"`
	HeadimgUrl     string    `xorm:"TEXT"`
	UserChannel    string    `xorm:"comment('用户渠道，eg:官网注册/微信/微博/抖音/知乎...') VARCHAR(255)"`
	UserSource     string    `xorm:"comment('用户来源，eg:xxx活动注册/xxx用户推荐') VARCHAR(255)"`
	UserReferer    int64     `xorm:"comment('用户推荐人，保留字段，判断用户推荐关系及营销价值') BIGINT(20)"`
	CreateTime     time.Time `xorm:"created not null default 'CURRENT_TIMESTAMP' DATETIME"`
	ModifyTime     time.Time `xorm:"updated not null default 'CURRENT_TIMESTAMP' DATETIME"`
	IsActivated    bool      `xorm:"default 'true' bool"`
	IdmId          string    `xorm:"unique VARCHAR(28)"`
}

// ReferedUser ReferedUser
type ReferedUser struct {
	Ysid           int64  `json:"user_id"`
	Name           string `json:"name"`
	RealName       string `json:"real_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	WechatUnionId  string `json:"unionid"`
	WechatOpenId   string `json:"openid"`
	WechatNickName string `json:"wechat_nick_name"`
	HeadimgUrl     string `json:"headimg_url"`
	UserChannel    string `json:"user_channel"`
	UserSource     string `json:"user_source"`
	UserReferer    int64  `json:"user_referer"`
}

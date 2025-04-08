package models

import (
	"time"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
)

const (
	// OffiaccountBindingPlatformCloud Yuansuan Cloud
	OffiaccountBindingPlatformCloud = "cloud"
	// OffiaccountBindingPlatformOMS Yuansuan OMS
	OffiaccountBindingPlatformOMS = "oms"
	// OffiaccountBindingNotificationTypeJob Job notification
	OffiaccountBindingNotificationTypeJob = "job"
	// OffiaccountBindingNotificationTypeBalance Balance notification
	OffiaccountBindingNotificationTypeBalance = "balance"

	OffiaccountAutoReplyModeSubscribe = "reply_subscribe"
	OffiaccountAutoReplyModeGeneral   = "reply_general"
)

// OffiaccountBinding 公众号绑定数据模型
type OffiaccountBinding struct {
	Id                    snowflake.ID `xorm:"not null pk BIGINT(20)"`
	Platform              string       `xorm:"not null VARCHAR(45)"`
	UserId                snowflake.ID `xorm:"BIGINT(20)"`
	CreateBy              snowflake.ID `xorm:"BIGINT(20)"`
	UpdateBy              snowflake.ID `xorm:"BIGINT(20)"`
	CompanyId             snowflake.ID `json:"company_id" xorm:"BIGINT(20)"`
	CompanyIds            string       `json:"company_ids" xorm:"not null TEXT"`
	WechatOpenid          string       `xorm:"not null VARCHAR(128)"`
	WechatUnionid         string       `xorm:"VARCHAR(28)"`
	WechatNickname        string       `xorm:"not null VARCHAR(128)"`
	WechatHeadimgurl      string       `xorm:"TEXT"`
	WechatLanguage        string       `xorm:"VARCHAR(45)"`
	UserGender            int32        `xorm:"not null default 0 TINYINT(4)"`
	UserCity              string       `xorm:"VARCHAR(128)"`
	NotificationType      string       `xorm:"VARCHAR(45)"`
	NotificationActivated int32        `xorm:"not null default 0 TINYINT(4)"`
	IsSubscribed          int32        `xorm:"not null default 0 TINYINT(4)"`
	SubscribeScene        string       `xorm:"VARCHAR(45)"`
	SubscribeTime         time.Time    `json:"subscribe_time" xorm:"not null default CURRENT_TIMESTAMP DATETIME"`
	UnsubscribeTime       time.Time    `json:"unsubscribe_time" xorm:"default null DATETIME"`
	ActivateTime          time.Time    `json:"activate_time" xorm:"default null DATETIME"`
	DeactivateTime        time.Time    `json:"deactivate_time" xorm:"default null DATETIME"`
}

// OffiaccountReplyRule 微信公众号消息自动回复规则
type OffiaccountReplyRule struct {
	Id         snowflake.ID `json:"id" xorm:"not null pk BIGINT(20)"`
	RuleName   string       `json:"rule_name" xorm:"not null VARCHAR(256)"`
	Keywords   string       `json:"keywords" xorm:"not null TEXT"`
	ReplyList  string       `josn:"reply_content" xorm:"TEXT"`
	ReplyMode  string       `json:"reply_mode" xorm:"not null VARCHAR(45)"`
	IsActive   string       `json:"is_active" xorm:"not null VARCHAR(20)"`
	CreatorId  snowflake.ID `json:"creator_id" xorm:"not null BIGINT(20)"`
	CreateTime time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP DATETIME"`
	UpdateTime time.Time    `json:"update_time" xorm:"default null DATETIME"`
}

// 待(公众号)通知的Job
type JobToNotify struct {
	Id         snowflake.ID `json:"id" xorm:"not null pk BIGINT(20)"`
	UserId     snowflake.ID `xorm:"BIGINT(20)"`
	JobId      snowflake.ID `xorm:"BIGINT(20)"`
	Status     int32        `xorm:"not null default 0 TINYINT(4)"`
	CreateTime time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP DATETIME"`
	UpdateTime time.Time    `json:"update_time" xorm:"default null DATETIME"`
}

type OffiaccountMenu struct {
	Id         snowflake.ID `json:"id" xorm:"not null pk BIGINT(20)"`
	AppId      string       `json:"app_id" xorm:"not null VARCHAR(45)"`
	Button     string       `json:"button" xorm:"text not null"`
	CreatorId  snowflake.ID `json:"creator_id" xorm:"BIGINT(20)"`
	CreateTime time.Time    `json:"create_time" xorm:"not null default CURRENT_TIMESTAMP DATETIME"`
	UpdateTime time.Time    `json:"update_time" xorm:"default null DATETIME"`
}

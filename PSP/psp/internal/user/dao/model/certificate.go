package model

import (
	"time"

	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type OpenapiUserCertificate struct {
	Id          snowflake.ID `json:"id" xorm:"pk BIGINT(20)"`
	UserId      snowflake.ID `json:"user_id" xorm:"BIGINT(20)"`
	Certificate string       `json:"certificate" xorm:"VARCHAR(20)"`
	CreatedAt   time.Time    `json:"created_at" xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
	UpdatedAt   time.Time    `json:"updated_at" xorm:"not null default 'CURRENT_TIMESTAMP' DATETIME"`
}

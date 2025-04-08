package dao

import (
	"encoding/json"
	"fmt"

	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"gorm.io/gorm"
)

// AuthzPolicy defines iam policy type.
type AuthzPolicy struct {
	Policy ladon.DefaultPolicy
}

type Policy struct {
	// 目前只支持主账号ID，即YSId
	// 如果userid为空，则代表公共基本策略
	ID              snowflake.ID `json:"Id,omitempty" gorm:"primary_key;column:id"`
	UserId          string       `json:"UserId" gorm:"type:varchar(255);column:userId;uniqueIndex:idx_name_user_policy"`
	PolicyName      string       `json:"PolicyName" gorm:"type:varchar(255);column:policyName;uniqueIndex:idx_name_user_policy"`
	Policy          AuthzPolicy  `json:"Policy" gorm:"-"`
	StatementShadow string       `json:"-" gorm:"column:statementShadow"`
	Version         string       `json:"Version" gorm:"column:version"`
}

// TableName maps to dao table name.
func (p *Policy) TableName() string {
	return "policy"
}

// String returns the string format of Policy.
func (ap AuthzPolicy) String() string {
	data, _ := json.Marshal(ap)

	return string(data)
}

func (p *Policy) BeforeCreateForRaw() error {
	p.StatementShadow = p.Policy.String()
	return nil
}

// BeforeCreate run before create database record.
func (p *Policy) BeforeCreate(tx *gorm.DB) error {
	p.StatementShadow = p.Policy.String()
	return nil
}

// BeforeUpdate run before update database record.
func (p *Policy) BeforeUpdate(tx *gorm.DB) error {
	return p.BeforeCreate(tx)
}

func (p *Policy) AfterFind(tx *gorm.DB) (err error) {
	if p.StatementShadow == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(p.StatementShadow), &p.Policy); err != nil {
		return fmt.Errorf("failed to unmarshal ladon policy statement Shadow: %w", err)
	}
	return nil
}

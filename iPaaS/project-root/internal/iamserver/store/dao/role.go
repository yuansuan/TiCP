package dao

import (
	"encoding/json"
	"fmt"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"gorm.io/gorm"
)

type Role struct {
	ID snowflake.ID `json:"Id,omitempty" gorm:"primary_key;column:id"`
	// userId 为空，代表系统基础角色，例如CloudComputeRole, CSPRole， 目前默认所有用户都开通这些策略
	UserId            string      `json:"UserId" gorm:"type:varchar(255);column:userId;uniqueIndex:idx_name_user_role"`
	RoleName          string      `json:"RoleName" gorm:"type:varchar(255);column:roleName;uniqueIndex:idx_name_user_role"`
	Description       string      `json:"Description" gorm:"column:description"`
	TrustPolicy       AuthzPolicy `json:"TrustPolicy" gorm:"-"`
	TrustPolicyShadow string      `json:"-" gorm:"column:trustPolicyShadow"`
}

// TableName maps to dao table name.
func (r *Role) TableName() string {
	return "role"
}

func (r *Role) BeforeCreateForRaw() error {
	s, _ := json.Marshal(r.TrustPolicy)
	r.TrustPolicyShadow = string(s)
	return nil
}

// BeforeCreate run before create database record.
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	s, _ := json.Marshal(r.TrustPolicy)
	r.TrustPolicyShadow = string(s)
	return nil
}

func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	return r.BeforeCreate(tx)
}

func (r *Role) AfterFind(tx *gorm.DB) (err error) {
	if r.TrustPolicyShadow == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(r.TrustPolicyShadow), &r.TrustPolicy); err != nil {
		return fmt.Errorf("failed to unmarshal trustpolicy Shadow: %w", err)
	}
	return nil
}

func (r *Role) UnmarshalPolicy() (err error) {
	if r.TrustPolicyShadow == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(r.TrustPolicyShadow), &r.TrustPolicy); err != nil {
		return fmt.Errorf("failed to unmarshal trustpolicy Shadow: %w", err)
	}
	return nil
}

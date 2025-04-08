package dao

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Secret struct {
	AccessKeyId     string    `json:"AccessKeyId"          gorm:"primary_key;column:accessKeyId"`
	AccessKeySecret string    `json:"AccessKeySecret"          gorm:"column:accessKeySecret"`
	SessionToken    string    `json:"SessionToken"    gorm:"column:sessionToken"`
	Expiration      time.Time `json:"Expiration"    gorm:"column:expiration"`
	// 目前只支持主账号AK，ParentUser即YSId
	ParentUser    string                 `json:"ParentUser"    gorm:"column:parentUser"`
	Claims        map[string]interface{} `json:"Claims" gorm:"-"`
	ClaimsShadows string                 `json:"-" gorm:"claimShadows"`
	Description   string                 `json:"Description"    gorm:"column:description"`
	Status        bool                   `json:"Status" gorm:"column:status"`
	// 用于标记是否是远算的产品账号AK，还是普通AK。 YSProductAccount就是产品账号
	Tag       string    `json:"Tag" gorm:"column:tag"`
	CreatedAt time.Time `json:"CreatedAt,omitempty" gorm:"column:createdAt"`
	UpdatedAt time.Time `json:"UpdatedAt,omitempty" gorm:"column:updatedAt"`
}

// TableName maps to dao table name.
func (s *Secret) TableName() string {
	return "secret"
}

// BeforeCreate run before create database record.
func (s *Secret) BeforeCreate(tx *gorm.DB) error {
	str, _ := json.Marshal(s.Claims)
	s.ClaimsShadows = string(str)
	return nil
}

func (s *Secret) BeforeUpdate(tx *gorm.DB) error {
	return s.BeforeCreate(tx)
}

func (s *Secret) AfterFind(tx *gorm.DB) (err error) {
	if s.ClaimsShadows == "" {
		return nil
	}
	if err := json.Unmarshal([]byte(s.ClaimsShadows), &s.Claims); err != nil {
		return fmt.Errorf("failed to unmarshal claims Shadow: %w", err)
	}
	return nil
}

package model

// Role Role
type Role struct {
	Id      int64  `yaml:"id" json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name    string `yaml:"name" json:"name" xorm:"not null unique VARCHAR(64)"`
	Comment string `yaml:"comment" json:"comment" xorm:"not null VARCHAR(256)"`

	// 0: custom
	// 1: super_admin
	Type int32 `yaml:"type" json:"type" xorm:"not null TINYINT(4)"`
}

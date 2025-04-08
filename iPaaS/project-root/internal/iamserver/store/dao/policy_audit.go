package dao

type PolicyAudit struct {
	ID           int64  `gorm:"column:id;primary_key" json:"id"`
	Subject      string `gorm:"column:subject" json:"subject"`
	PolicyShadow string `gorm:"column:policyShadow" json:"policyShadow"`
}

// TableName sets the insert table name for this struct type
func (p *PolicyAudit) TableName() string {
	return "policy_audit"
}

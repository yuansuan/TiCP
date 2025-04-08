package model

type Resource struct {
	Id          int64  `json:"id" xorm:"pk autoincr BIGINT(20)"`
	Name        string `json:"name" xorm:"name"`
	Action      string `json:"action" xorm:"VARCHAR(128)"`
	Type        string `json:"type" xorm:"VARCHAR(64)"`
	DisplayName string `json:"display_name" xorm:"VARCHAR(128)"`
	Custom      int8   `json:"custom" xorm:"TINYINT(1)"`
	ExternalId  int64  `json:"external_id" xorm:"BIGINT(20)"`
	ParentId    int64  `json:"parent_id" xorm:"BIGINT(20)"`
}

package model

import "github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"

type OrgStructure struct {
	Id       snowflake.ID `json:"id" xorm:"pk BIGINT(20)"`
	Name     string       `json:"name" xorm:"not null varchar(128)"`
	Comment  string       `json:"comment" xorm:"varchar(255)"`
	ParentId snowflake.ID `json:"parent_id" xorm:"BIGINT(20) default 0"`
}

type OrgUserRelation struct {
	Id     snowflake.ID `json:"id" xorm:"pk BIGINT(20)"`
	UserId snowflake.ID `json:"user_id" xorm:"not null BIGINT(20)"`
	OrgId  snowflake.ID `json:"org_id" xorm:"not null BIGINT(20)"`
}

type OrgAndUserStructure struct {
	UserList []OrgUser
	OrgList  []OrgStructure
}

type OrgUser struct {
	Id   snowflake.ID
	Name string
}

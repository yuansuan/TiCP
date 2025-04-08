package dao

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// OrgStructureDaoImpl OrgStructureDaoImpl
type OrgStructureDaoImpl struct {
}

func NewOrgStructureDaoImpl() *OrgStructureDaoImpl {
	return &OrgStructureDaoImpl{}
}

// Add Add
func (dao OrgStructureDaoImpl) Add(org *model.OrgStructure) (int64, error) {
	session := GetSession()
	defer session.Close()

	node, err := snowflake.GetInstance()
	if err != nil {
		return 0, err
	}

	org.Id = node.Generate()

	if _, err := session.Insert(org); err != nil {
		return 0, err
	}

	return org.Id.Int64(), nil
}

// Delete Delete
func (dao OrgStructureDaoImpl) Delete(orgId snowflake.ID) error {
	session := GetSession()
	defer session.Close()

	org := model.OrgStructure{Id: orgId}
	_, err := session.Delete(org)

	return err
}

// Update Update
func (dao OrgStructureDaoImpl) Update(org *model.OrgStructure) error {
	session := GetSession()
	defer session.Close()

	_, err := session.MustCols("name", "type", "remark", "parentId").ID(org.Id).Update(&org)

	return err
}

func (dao OrgStructureDaoImpl) GetChildrenOrgStructure(parentOrgId snowflake.ID) (*model.OrgAndUserStructure, error) {
	session := GetSession()
	defer session.Close()

	childrenOrgStructure := make([]model.OrgStructure, 0)

	err := session.Where("parent_id = ?", parentOrgId).Find(&childrenOrgStructure)

	if err != nil {
		return nil, err
	}

	userList := make([]model.OrgUser, 0)

	err = session.Table("org_user_relation").Alias("ou").
		Join("LEFT OUTER", "user", "ou.user_id = user.id").Where("ou.org_id = ?", parentOrgId).
		Cols("user.id", "user.name").Find(&userList)
	if err != nil {
		return nil, err
	}

	childrenOrgAndUser := model.OrgAndUserStructure{
		UserList: userList,
		OrgList:  childrenOrgStructure,
	}

	return &childrenOrgAndUser, nil
}

// DeleteOrgUserByUserIds DeleteOrgUserByUserIds
func (dao OrgStructureDaoImpl) DeleteOrgUserByUserIds(userIds ...snowflake.ID) error {
	session := GetSession()
	defer session.Close()

	_, err := session.In("user_id", userIds).Delete(&model.OrgUserRelation{})

	return err
}

// DeleteOrgUserByOrgIds DeleteOrgUserByOrgIds
func (dao OrgStructureDaoImpl) DeleteOrgUserByOrgIds(orgIds ...snowflake.ID) error {
	session := GetSession()
	defer session.Close()

	_, err := session.In("org_id", orgIds).Delete(&model.OrgUserRelation{})

	return err
}

// DeleteOrgUserByIds DeleteOrgUserByIds
func (dao OrgStructureDaoImpl) DeleteOrgUserByIds(ouIds ...snowflake.ID) error {
	session := GetSession()
	defer session.Close()

	_, err := session.In("id", ouIds).Delete(&model.OrgUserRelation{})

	return err
}

// InsertOrgUser InsertOrgUser
func (dao OrgStructureDaoImpl) InsertOrgUser(ou ...*model.OrgUserRelation) error {
	session := GetSession()
	defer session.Close()

	_, err := session.InsertMulti(ou)
	return err
}

func (dao OrgStructureDaoImpl) ListOrgMember(orgID snowflake.ID) ([]*model.User, error) {
	userList := make([]*model.User, 0)
	session := GetSession()
	defer session.Close()

	err := session.Table("org_user_relation").Alias("ou").
		Join("LEFT OUTER", "user", "ou.user_id = user.id").
		Where("ou.org_id = ?", orgID).
		Cols("user.id", "user.name").
		Find(&userList)
	return userList, err
}

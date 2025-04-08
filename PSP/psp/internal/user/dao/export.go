package dao

import (
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

// UserDao UserDao
type UserDao interface {
	Add(user model.User) (id int64, err error)
	GetByID(id int64) (model.User, bool, error)
	Get(id int64) (model.User, bool, error)
	GetIncludeDeleted(id int64) (model.User, bool, error)
	GetByUserID(id int64) (model.User, error)
	BatchGet(ids []int64) (users []*model.User, err error)

	//BatchGet(UserIdentity)
	GetIdByName(name string) (id int64, err error)
	Exist(id int64) bool
	Exists(ids ...int64) error
	List(enabled bool, query string, order string, desc bool, page int64, pageSize int64, isInternal bool) (users []*model.User, total int64, err error)
	ListAll(index, pageSize int64) (users []*model.User, total int64, err error)
	Update(user model.User) error
	Delete(id int64) error
	UpdateOrInsertUser(user model.User) error
	InsertUser(user model.User) (*model.User, error)
	ListAllUser() ([]*model.User, error)
	ListAllAdminUser() ([]*model.User, error)
	ListUserByName(names []string) ([]model.User, error)
	ListUserLikeName(names string, ids []int64) ([]model.User, error)
	DeleteUserByID(ids []int64) error
	GetUserByUserId(id snowflake.ID) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
}

// OrgStructureDao OrgStructureDao
type OrgStructureDao interface {
	Add(org *model.OrgStructure) (int64, error)
	Delete(orgId snowflake.ID) error
	Update(org *model.OrgStructure) error
	GetChildrenOrgStructure(parentOrgId snowflake.ID) (*model.OrgAndUserStructure, error)
	DeleteOrgUserByUserIds(userIds ...snowflake.ID) error
	DeleteOrgUserByOrgIds(orgIds ...snowflake.ID) error
	DeleteOrgUserByIds(ouIds ...snowflake.ID) error
	InsertOrgUser(ou ...*model.OrgUserRelation) error
	ListOrgMember(id snowflake.ID) ([]*model.User, error)
}

type CertificateDao interface {
	Add(certificate *model.OpenapiUserCertificate) (snowflake.ID, error)
	DelByUserID(userID snowflake.ID) error
	GetByUserID(userID snowflake.ID) (cert model.OpenapiUserCertificate, exist bool, err error)
	CheckCert(certificate string) (*model.User, bool, error)
}

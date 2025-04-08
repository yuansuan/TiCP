package service

import (
	"context"
	"github.com/go-ldap/ldap/v3"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/ptype"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
)

type AuthService interface {
	CheckUserPass(ctx context.Context, user dto.UserRequest) error
	CheckLdapUserPass(ctx context.Context, user dto.UserRequest) (bool, error)
	GetConn(serverAddr string) (*ldap.Conn, error)
	GetOnlineList(ctx *gin.Context, req dto.OnlineListRequest) (*dto.OnlineUserListResponse, error)
	GetOnlineListByUser(ctx *gin.Context, req dto.OnlineListByUserRequest) (*dto.OnlineUserInfoListResponse, error)
}

type UserService interface {
	GetIdByName(ctx context.Context, userName string) (int64, error)
	LoginCheck(ctx context.Context, userId int64, userName string) (*dto.LoginSuccessResponse, error)
	Get(ctx context.Context, userId int64) (*model.User, error)
	GetIncludeDeleted(ctx context.Context, userId int64) (*model.User, error)
	Add(ctx context.Context, user model.User) (int64, error)
	Exist(ctx context.Context, userId int64) (bool, error)
	QueryByCond(ctx *gin.Context, query dto.QueryByCondRequest) ([]*model.User, int64, error)
	ListAll(ctx context.Context, page *ptype.PageReq) ([]*model.User, int64, error)
	BatchGetUser(ctx context.Context, ids []int64) ([]*model.User, error)
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, userId int64) error
	ActiveUser(ctx context.Context, userId int64) error
	InactiveUser(ctx context.Context, userId int64) error
	UpdatePassword(ctx context.Context, req dto.UpdatePassRequest) error
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	AddUserWithRole(ctx context.Context, req dto.UserAddRequest) (int64, error)
	QueryUserRole(ctx *gin.Context, req dto.QueryByCondRequest) (*dto.UserListResponse, error)
	Detail(ctx context.Context, userId int64) (*dto.UserDetailResponse, error)
	OptionList(ctx *gin.Context, name string, filterPerm int64) (userOptionList []*dto.UserOptionResponse, err error)
	GetAllUser(ctx context.Context) ([]*model.User, error)
	GetUserRoleNames(ctx context.Context, userId string) string
	ResetPassword(ctx context.Context, userId string) (string, error)
	GenOpenapiCertificate(ctx context.Context, userId string, over bool) (string, error)
	CheckOpenapiCertificate(ctx context.Context, certificate string) (*model.User, bool, error)
}

type LicenseService interface {
	GetMachineID(ctx context.Context) (string, error)
	GetLicense(ctx context.Context) (*dto.License, error)
	UpdateLicense(ctx context.Context, license *dto.License) error
	CheckLicenseExpired(ctx context.Context) error
}

type OrgService interface {
	CreateOrg(ctx context.Context, req dto.CreateOrgRequest) error
	Delete(ctx context.Context, orgID snowflake.ID) error
	UpdateOrg(ctx context.Context, req dto.UpdateOrgRequest) error
	AddOrgMember(ctx context.Context, orgID string, userIDs []string) error
	DeleteOrgMemberByID(ctx context.Context, ids []string) error
	DeleteOrgMemberByUserID(ctx context.Context, userIDs []string) error
	UpdateOrgMember(ctx context.Context, req dto.UpdateOrgMemberRequest) error
	ListOrgMember(ctx context.Context, orgID snowflake.ID) ([]*dto.ListMemberResponse, error)
}

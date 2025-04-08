package service

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dto"
)

type RoleService interface {
	AddRole(ctx context.Context, role *model.Role) error
	GetRole(ctx context.Context, roleId int64) (*model.Role, error)
	GetRoles(ctx context.Context, roleId []int64) ([]*model.Role, error)
	GetRoleDetail(ctx context.Context, roleId int64) (*dto.RoleDetail, error)
	ListRolePerms(ctx context.Context, roleId int64) (*dto.Resources, error)
	ListByType(ctx context.Context, typ int) (roles []*model.Role, err error)
	UpdateRole(ctx context.Context, role *model.Role) error
	DeleteRole(ctx context.Context, roleId int64) error
	ListRole(ctx *gin.Context, req *boring.ListRequest) ([]*model.Role, int64, error)
	AddRolePerms(ctx context.Context, roleId int64, perms []int64) error
	UpdateRolePerms(ctx context.Context, roleId int64, perms []int64) error
	RemoveRolePerms(ctx context.Context, roleId int64, perms []int64) error
	InternalAddRolePerms(ctx context.Context, roleId int64, perms []int64) error
	InternalRemoveRolePerms(ctx context.Context, roleId int64, perms []int64) error
	AddObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error
	ListObjectsRoles(ctx context.Context, req []*rbac.ObjectID, needImplicitRoles bool) (*rbac.ListObjectsRolesResp, error)
	UpdateObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error
	RemoveObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error
	GetRoleByObjectID(ctx context.Context, id *rbac.ObjectID) (roleList []*model.Role, err error)
	GetRoleByName(ctx context.Context, name string) (*model.Role, error)
	SetLdapUserDefRole(ctx context.Context, roleID int64) error
}

type PermService interface {
	AddResource(ctx context.Context, res *model.Resource) (*model.Resource, error)
	GetResource(ctx context.Context, resId int64) (*model.Resource, error)
	GetResources(ctx context.Context, resIds []int64) ([]*model.Resource, error)
	FindResourceByNameOrExId(ctx context.Context, resources ...*rbac.ResourceIdentity) ([]*model.Resource, []*rbac.ResourceIdentity, error)
	UpdateResource(ctx context.Context, resId int64, res *model.Resource) error
	DeleteResource(ctx context.Context, resId int64) error
	ListResource(ctx context.Context, req *boring.ListRequest) ([]*model.Resource, int64, error)
	CheckPermAllExists(ctx context.Context, ids ...int64) error
	AddRolePerms(ctx context.Context, roleID int64, permIDs []int64) error
	ListObjectResources(ctx context.Context, objectID *rbac.ObjectID, resType []string) ([]*model.Resource, error)
	ListObjectPermissions(ctx context.Context, objectID *rbac.ObjectID) ([]*model.Resource, error)
	CheckResourcesPerm(ctx context.Context, req *rbac.CheckResourcesPermRequest) (*rbac.PermCheckResponse, error)
	CheckSelfPermissions(ctx context.Context, resources ...*rbac.SimpleResource) error
	CheckApiPermission(ctx context.Context, userID int64, url, method string) (bool, error)
	GetObjectsByResource(ctx context.Context, req *rbac.ResourceID) (*rbac.ObjectIDs, error)
}

package impl

import (
	"context"
	"fmt"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/gin-gonic/gin"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/collection"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/casbin"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/resource"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/role"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/util"
	usesconf "github.com/yuansuan/ticp/PSP/psp/internal/user/config"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

type RoleServiceImpl struct {
	resourceDao      resource.Dao
	roleDao          role.Dao
	casDao           *casbin.Dao
	superAdminRoleID int64
}

func (srv *RoleServiceImpl) SetLdapUserDefRole(ctx context.Context, roleID int64) error {
	role, err := srv.GetRole(ctx, roleID)
	if err != nil {
		return err
	}
	if role.Type == int32(rbac.RoleType_ROLE_SUPER_ADMIN) {
		return status.Error(errcode.ErrRBACSetSuperAdmin, "")
	}
	_, err = client.GetInstance().SysConfig.SetRBACDefaultRoleId(ctx, &sysconfig.SetRBACDefaultRoleIdRequest{
		RoleId: roleID,
	})
	tracelog.Info(ctx, fmt.Sprintf("set ldap defRole success, roleID[%v]", roleID))
	return err
}

func (srv *RoleServiceImpl) GetRoleByObjectID(ctx context.Context, id *rbac.ObjectID) (roleList []*model.Role, err error) {
	roleIds := srv.casDao.GetObjectRoles(id)
	for _, roleId := range roleIds {
		if role, err := srv.GetRole(ctx, roleId); err == nil {
			roleList = append(roleList, role)
		}
	}

	return
}

func NewRoleService() (service.RoleService, error) {

	ctx := context.TODO()

	resourceDao := resource.NewResourceDaoImpl()
	roleDao := role.NewRoleDaoImpl()

	superAdminRoleID, err := roleDao.ListByType(ctx, int(rbac.RoleType_ROLE_SUPER_ADMIN))
	if err != nil {
		return nil, err
	}

	logging.Default().Infof("rbac config file path: %s", config.Custom.RBACConfigPath)

	roleService := &RoleServiceImpl{
		resourceDao:      resourceDao,
		roleDao:          roleDao,
		casDao:           casbin.New(boot.MW.DefaultORMEngine(), config.Custom.RBACConfigPath),
		superAdminRoleID: superAdminRoleID[0].Id,
	}

	return roleService, nil
}

func (srv *RoleServiceImpl) AddRole(ctx context.Context, role *model.Role) error {
	oldRole, _ := srv.roleDao.GetByName(ctx, role.Name)
	if oldRole != nil && oldRole.Name == role.Name {
		return status.Error(errcode.ErrRBACRoleNameExist, "")
	}

	if err := srv.roleDao.Add(ctx, role); err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("add role success, role[%+v]", role))
	return nil
}

func (srv *RoleServiceImpl) GetRole(ctx context.Context, roleId int64) (*model.Role, error) {
	return srv.roleDao.Get(ctx, roleId)
}

func (srv *RoleServiceImpl) GetRoles(ctx context.Context, roleIds []int64) ([]*model.Role, error) {
	return srv.roleDao.Gets(ctx, roleIds)
}

func (srv *RoleServiceImpl) GetRoleByName(ctx context.Context, name string) (*model.Role, error) {
	return srv.roleDao.GetByName(ctx, name)
}

func (srv *RoleServiceImpl) GetRoleDetail(ctx context.Context, roleId int64) (*dto.RoleDetail, error) {

	roleInfo, err := srv.GetRole(ctx, roleId)
	if err != nil {
		return &dto.RoleDetail{}, err
	}

	objectsID := srv.casDao.GetRoleObjects(roleId)
	if err != nil {
		return &dto.RoleDetail{}, err
	}

	perms, err := srv.ListRolePerms(ctx, roleId)
	if err != nil {
		return &dto.RoleDetail{}, err
	}

	resp, err := client.GetInstance().SysConfig.GetRBACDefaultRoleId(ctx, &sysconfig.GetRBACDefaultRoleIdRequest{})
	if err != nil {
		return &dto.RoleDetail{}, err
	}

	return &dto.RoleDetail{
		Role:      util.ToDTO.RoleInfo(roleInfo, resp.RoleId),
		Objects:   objectsID,
		Resources: perms.Perms,
	}, nil
}

func (srv *RoleServiceImpl) ListRolePerms(ctx context.Context, roleId int64) (*dto.Resources, error) {
	if err := srv.roleDao.ShouldAllExists(ctx, roleId); err != nil {
		return nil, err
	}

	permsID := srv.casDao.GetRolePerms(roleId)
	if len(permsID) <= 0 {
		return &dto.Resources{}, nil
	}

	resList, err := srv.resourceDao.GetCustomResList(ctx, permsID)
	if err != nil {
		return &dto.Resources{}, err
	}
	return &dto.Resources{
		Perms: util.ToDTO.Resources(resList),
		Total: int64(len(resList)),
	}, nil
}

func (srv *RoleServiceImpl) ListByType(ctx context.Context, typ int) (roles []*model.Role, err error) {
	return srv.roleDao.ListByType(ctx, typ)
}

func (srv *RoleServiceImpl) UpdateRole(ctx context.Context, role *model.Role) error {
	if err := srv.ShouldNoSuperAdmin(role.Id); err != nil {
		return err
	}

	if err := srv.roleDao.Update(ctx, role.Id, role); err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update role success, role[%+v]", role))
	return nil
}

// ShouldNoSuperAdmin 超级管理员角色不允许修改
func (srv *RoleServiceImpl) ShouldNoSuperAdmin(ids ...int64) error {
	if len(ids) == 0 {
		return nil
	}

	if collection.ContainInt64(ids, srv.superAdminRoleID) {
		return status.Error(errcode.ErrRBACNotAllowOpSuperAdminRole, "")
	}

	return nil
}

func (srv *RoleServiceImpl) DeleteRole(ctx context.Context, roleId int64) error {
	if err := srv.ShouldNoSuperAdmin(roleId); err != nil {
		return err
	}

	if usesconf.GetConfig().LdapConf.Enable {
		resp, err := client.GetInstance().SysConfig.GetRBACDefaultRoleId(ctx, &sysconfig.GetRBACDefaultRoleIdRequest{})
		if err != nil {
			return err
		}

		if usesconf.GetConfig().LdapConf.Enable && roleId == resp.RoleId {
			return status.Error(errcode.ErrRBACRoleUsedByLdap, "")
		}
	}

	userIds := srv.casDao.GetRoleObjects(roleId)

	if len(userIds) > 0 {
		return status.Error(errcode.ErrRBACRoleUsed, "")
	}

	err := srv.roleDao.Delete(ctx, roleId)
	if err != nil {
		return err
	}
	srv.casDao.RemoveRole(roleId)

	tracelog.Info(ctx, fmt.Sprintf("delete role success, roleID[%v]", roleId))
	return nil
}

func (srv *RoleServiceImpl) ListRole(ctx *gin.Context, req *boring.ListRequest) ([]*model.Role, int64, error) {

	userID := ginutil.GetUserID(ctx)
	roleIDs := srv.casDao.GetObjectRoles(&rbac.ObjectID{
		Id:   snowflake.ID(userID).String(),
		Type: rbac.ObjectType_USER,
	})
	// 非超级管理员不能在列表中看到超级管理员角色
	isAdmin := len(roleIDs) > 0 && collection.ContainInt64(roleIDs, srv.superAdminRoleID)

	return srv.roleDao.List(ctx, req, isAdmin)
}

func (srv *RoleServiceImpl) AddRolePerms(ctx context.Context, roleId int64, perms []int64) error {

	if err := srv.ShouldNoSuperAdmin(roleId); err != nil {
		return status.Errorf(errcode.ErrRBACNotAllowOpSuperAdminRole, "")
	}
	if err := srv.resourceDao.NoInternalPerm(ctx, perms...); err != nil {
		return status.Errorf(errcode.ErrRBACInternalPermOnlyGiveAdmin, "")
	}
	if err := srv.rolePermsAllExists(ctx, roleId, perms); err != nil {
		return status.Errorf(errcode.ErrRBACPermissionNotFound, "")
	}

	perms, err := srv.resourceDao.FindChildResIds(ctx, perms)
	if err != nil {
		return status.Errorf(errcode.ErrRBACChildPerm, "")
	}

	srv.casDao.AddRolePerms(roleId, perms)

	tracelog.Info(ctx, fmt.Sprintf("add role perms success, roleID[%v], perms:[%v]", roleId, perms))
	return nil
}

func (srv *RoleServiceImpl) InternalAddRolePerms(ctx context.Context, roleId int64, perms []int64) error {
	if err := srv.rolePermsAllExists(ctx, roleId, perms); err != nil {
		return err
	}

	perms, err := srv.resourceDao.FindChildResIds(ctx, perms)
	if err != nil {
		return status.Errorf(errcode.ErrRBACChildPerm, "查询关联子级权限失败")
	}

	srv.casDao.AddRolePerms(roleId, perms)

	tracelog.Info(ctx, fmt.Sprintf("add internal role perms success, roleID[%v], perms:[%v]", roleId, perms))
	return nil
}

func (srv *RoleServiceImpl) rolePermsAllExists(ctx context.Context, roleId int64, perms []int64) error {

	if err := srv.roleDao.ShouldAllExists(ctx, roleId); err != nil {
		return err
	}
	if len(perms) <= 0 {
		return nil
	}
	return srv.resourceDao.AllExists(ctx, perms...)
}

func (srv *RoleServiceImpl) UpdateRolePerms(ctx context.Context, roleId int64, perms []int64) error {

	if err := srv.ShouldNoSuperAdmin(roleId); err != nil {
		return err
	}
	if err := srv.resourceDao.NoInternalPerm(ctx, perms...); err != nil {
		return err
	}
	if err := srv.rolePermsAllExists(ctx, roleId, perms); err != nil {
		return err
	}

	perms, err := srv.resourceDao.FindChildResIds(ctx, perms)
	if err != nil {
		return status.Errorf(errcode.ErrRBACChildPerm, "查询关联子级权限失败")
	}

	srv.casDao.UpdateRolePerms(roleId, perms)

	tracelog.Info(ctx, fmt.Sprintf("update role perms success, roleID[%v], perms:[%v]", roleId, perms))
	return nil
}

func (srv *RoleServiceImpl) RemoveRolePerms(ctx context.Context, roleId int64, perms []int64) error {

	if err := srv.ShouldNoSuperAdmin(roleId); err != nil {
		return err
	}

	if err := srv.resourceDao.NoInternalPerm(ctx, perms...); err != nil {
		return err
	}
	if err := srv.rolePermsAllExists(ctx, roleId, perms); err != nil {
		return err
	}

	perms, err := srv.resourceDao.FindChildResIds(ctx, perms)
	if err != nil {
		return status.Errorf(errcode.ErrRBACChildPerm, "查询关联子级权限失败")
	}

	srv.casDao.RemoveRolePerms(roleId, perms)

	tracelog.Info(ctx, fmt.Sprintf("remove role perms success, roleID[%v], perms:[%v]", roleId, perms))
	return nil
}

func (srv *RoleServiceImpl) InternalRemoveRolePerms(ctx context.Context, roleId int64, perms []int64) error {
	if err := srv.rolePermsAllExists(ctx, roleId, perms); err != nil {
		return err
	}

	perms, err := srv.resourceDao.FindChildResIds(ctx, perms)
	if err != nil {
		return status.Errorf(errcode.ErrRBACChildPerm, "查询关联子级权限失败")
	}

	srv.casDao.RemoveRolePerms(roleId, perms)

	tracelog.Info(ctx, fmt.Sprintf("remove internal role perms success, roleID[%v], perms:[%v]", roleId, perms))
	return nil
}

// AddObjectRoles 给sub赋角色
func (srv *RoleServiceImpl) AddObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error {

	if err := srv.objectRolesAllExists(ctx, req); err != nil {
		return err
	}

	srv.casDao.AddObjectRoles(req)

	tracelog.Info(ctx, fmt.Sprintf("add user role success, userID[%v], roles:[%v]", req.Id.Id, req.Roles))
	return nil
}

func (srv *RoleServiceImpl) objectRolesAllExists(ctx context.Context, req *rbac.ObjectRoles) error {
	if err := srv.objectIDAllExists(ctx, req.Id); err != nil {
		return err
	}

	return srv.roleDao.ShouldAllExists(ctx, req.Roles...)
}

// 查询用户是否存在
func (srv *RoleServiceImpl) objectIDAllExists(ctx context.Context, ids ...*rbac.ObjectID) error {
	var users []*user.UserIdentity
	for _, id := range ids {
		if id.Type == rbac.ObjectType_USER {
			users = append(users, util.ToUserManagement.UserID(id))
		}
	}
	if len(users) != 0 {
		if err := srv.userIdentityAllExists(ctx, users...); err != nil {
			return err
		}
	}

	return nil
}

func (srv *RoleServiceImpl) userIdentityAllExists(ctx context.Context, ids ...*user.UserIdentity) error {
	if len(ids) == 0 {
		return nil
	}
	users, err := client.GetInstance().User.BatchGetUser(ctx, &user.UserIdentities{
		UserIdentities: ids,
	})
	if err != nil {
		return err
	}

	if len(users.UserObj) != len(ids) {
		return status.Errorf(errcode.ErrUserNotFound, "%d users not found", len(ids)-len(users.UserObj))
	}

	return nil
}

func (srv *RoleServiceImpl) ListObjectsRoles(ctx context.Context, ids []*rbac.ObjectID, needImplicitRoles bool) (*rbac.ListObjectsRolesResp, error) {
	if err := srv.objectIDAllExists(ctx, ids...); err != nil {
		return nil, err
	}

	roleIDSet := map[int64]bool{}

	objectRolesList := make([]*rbac.ObjectRoles, 0, len(ids))
	objectImplicitRolesList := make([]*rbac.ObjectRoles, 0, len(ids))

	for _, id := range ids {
		roleIDs := srv.casDao.GetObjectRoles(id)
		for _, roleID := range roleIDs {
			roleIDSet[roleID] = true
		}
		objectRolesList = append(objectRolesList, &rbac.ObjectRoles{
			Id:    id,
			Roles: roleIDs,
		})

		if needImplicitRoles {
			implicitRoles := []int64{}

			objectImplicitRolesList = append(objectImplicitRolesList, &rbac.ObjectRoles{
				Id:    id,
				Roles: collection.UniqueInt64Array(implicitRoles),
			})
		}
	}

	roleIDs := make([]int64, 0, len(roleIDSet))
	for id := range roleIDSet {
		roleIDs = append(roleIDs, id)
	}

	roles, err := srv.roleDao.Gets(ctx, roleIDs)
	if err != nil {
		return nil, err
	}
	return &rbac.ListObjectsRolesResp{
		ObjectRolesList:         objectRolesList,
		ObjectImplicitRolesList: objectImplicitRolesList,
		Roles:                   util.ToGRPC.Roles(roles),
	}, nil

}

// UpdateObjectRoles 修改sub的角色
func (srv *RoleServiceImpl) UpdateObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error {
	if err := srv.preCheckUpdateObjectRoles(ctx, req); err != nil {
		return err
	}

	srv.casDao.UpdateObjectRoles(req)

	tracelog.Info(ctx, fmt.Sprintf("add user role success, userID[%v], roles:[%v]", req.Id.Id, req.Roles))
	return nil
}

func (srv *RoleServiceImpl) preCheckUpdateObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error {

	if err := srv.objectRolesAllExists(ctx, req); err != nil {
		return err
	}

	return nil
}

func (srv *RoleServiceImpl) RemoveObjectRoles(ctx context.Context, req *rbac.ObjectRoles) error {
	// internal user can't remove admin role
	if req.Id.Type == rbac.ObjectType_USER {
		if err := srv.ShouldNoSuperAdmin(req.Roles...); err != nil {
			return err
		}
	}

	// 如果没传具体角色就全删
	if len(req.Roles) == 0 {
		req.Roles = srv.casDao.GetObjectRoles(req.Id)
	}

	srv.casDao.RemoveObjectRoles(req)

	tracelog.Info(ctx, fmt.Sprintf("remove user role success, userID[%v], roles:[%v]", req.Id.Id, req.Roles))
	return nil
}

package impl

import (
	"context"
	"encoding/json"
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/casbin"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/resource"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/role"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type PermServiceImpl struct {
	resourceDao resource.Dao
	roleDao     role.Dao
	casDao      *casbin.Dao
}

func NewPermService() (service.PermService, error) {

	permService := &PermServiceImpl{
		resourceDao: resource.NewResourceDaoImpl(),
		roleDao:     role.NewRoleDaoImpl(),
		casDao:      casbin.New(boot.MW.DefaultORMEngine(), config.Custom.RBACConfigPath),
	}

	return permService, nil
}

func (srv PermServiceImpl) CheckPermAllExists(ctx context.Context, ids ...int64) error {
	return srv.resourceDao.AllExists(ctx, ids...)
}

func (srv PermServiceImpl) AddRolePerms(ctx context.Context, roleID int64, permIDs []int64) error {

	// Add perms to the role
	// Previous operations can ensure this function will not return false
	// 获取关联的子级权限(包括父级)
	allResIds, err := srv.resourceDao.FindChildResIds(ctx, permIDs)
	if err != nil {
		return status.Errorf(errcode.ErrRBACChildPerm, "查询关联子级权限失败")
	}

	srv.casDao.AddRolePerms(roleID, allResIds)
	return nil
}

func (srv PermServiceImpl) AddResource(ctx context.Context, res *model.Resource) (*model.Resource, error) {
	if strutil.IsEmpty(res.Action) {
		res.Action = common.ResourceActionNONE
	}

	err := srv.resourceDao.Add(ctx, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (srv PermServiceImpl) GetResource(ctx context.Context, resId int64) (*model.Resource, error) {
	return srv.resourceDao.Get(ctx, resId)
}

func (srv PermServiceImpl) GetResources(ctx context.Context, resIds []int64) ([]*model.Resource, error) {
	return srv.resourceDao.Gets(ctx, resIds)
}

func (srv PermServiceImpl) FindResourceByNameOrExId(ctx context.Context, resources ...*rbac.ResourceIdentity) ([]*model.Resource, []*rbac.ResourceIdentity, error) {
	return srv.resourceDao.FindResourceByNameOrExId(ctx, resources...)
}

func (srv PermServiceImpl) UpdateResource(ctx context.Context, resId int64, res *model.Resource) error {
	return srv.resourceDao.Update(ctx, resId, res)
}

func (srv PermServiceImpl) DeleteResource(ctx context.Context, resId int64) error {

	if res, err := srv.resourceDao.Get(ctx, resId); err != nil {
		return err
	} else if res.Type == common.PermissionResourceTypeSystem ||
		res.Type == common.PermissionResourceTypeInternal ||
		res.Type == common.PermissionResourceTypeApi {
		return status.Error(errcode.ErrRBACCantRemoveInternalPerm, "")
	}

	err := srv.resourceDao.Delete(ctx, resId)
	if err != nil {
		return err
	}

	srv.casDao.RemovePerm(resId)
	return nil
}

func (srv PermServiceImpl) ListResource(ctx context.Context, req *boring.ListRequest) ([]*model.Resource, int64, error) {
	return srv.resourceDao.List(ctx, req, common.ENABLE_CUSTOM)
}

func (srv PermServiceImpl) ListObjectResources(ctx context.Context, objectID *rbac.ObjectID, resType []string) ([]*model.Resource, error) {

	if err := srv.objectIDAllExists(ctx, objectID); err != nil {
		return nil, err
	}

	//if err := srv.havePermOrSelfInfo(ctx, objectID); err != nil {
	//	return nil, err
	//}

	resList, err := srv.ListObjectPermissions(ctx, objectID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Resource, 0)
	for _, res := range resList {
		for _, v := range resType {
			if res.Type == v {
				result = append(result, res)
			}
		}
	}

	return result, nil
}

// 查询用户是否存在
func (srv *PermServiceImpl) objectIDAllExists(ctx context.Context, ids ...*rbac.ObjectID) error {
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

func (srv *PermServiceImpl) userIdentityAllExists(ctx context.Context, ids ...*user.UserIdentity) error {
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

var errNoPermission = status.Error(errcode.ErrRBACNoPermission, errcode.NoPermission)

// 查询请求者与查询用户是否是同一人，或者是否是拥有用户管理权限的人
func (srv *PermServiceImpl) havePermOrSelfInfo(ctx context.Context, objectIDs ...*rbac.ObjectID) error {

	if !(len(objectIDs) == 1 && objectIDs[0].Type == rbac.ObjectType_USER) {
		return errNoPermission
	}

	// 请求者与查询用户是否是同一人
	uid, err := boot.Middleware.GetUserIDString(ctx)
	if err != nil {
		return status.Error(errcode.ErrRBACFailedGetUserID, err.Error())
	}
	if uid == objectIDs[0].Id {
		return nil
	}
	return errNoPermission
}

func (srv *PermServiceImpl) CheckSelfPermissions(ctx context.Context, resources ...*rbac.SimpleResource) error {
	uid, err := boot.Middleware.GetUserIDString(ctx)
	if err != nil {
		return status.Error(errcode.ErrRBACFailedGetUserID, err.Error())
	}

	resp, err := srv.CheckResourcesPerm(ctx, &rbac.CheckResourcesPermRequest{
		Id: &rbac.ObjectID{
			Id:   uid,
			Type: rbac.ObjectType_USER,
		},
		Resources: util.ResourceIdentitiesBySimple(&rbac.SimpleResources{
			Resources: resources,
		}),
	})
	if err != nil {
		return err
	}
	if !resp.Pass {
		return errNoPermission
	}
	return nil
}

// CheckResourcesPerm CheckResourcesPerm
func (srv *PermServiceImpl) CheckResourcesPerm(ctx context.Context, req *rbac.CheckResourcesPermRequest) (*rbac.PermCheckResponse, error) {
	resList, notFound, err := srv.resourceDao.FindResourceByNameOrExId(ctx, req.Resources...)
	if err != nil {
		return nil, err
	}
	if len(notFound) != 0 {
		return nil, status.Errorf(errcode.ErrRBACPermissionNotFound, "%+v not found", notFound)
	}
	check := &rbac.CheckPermissionsRequest{
		Id: req.Id,
	}
	for _, res := range resList {
		check.PermissionIds = append(check.PermissionIds, res.Id)
	}
	return srv.checkPermissions(ctx, check)
}

// CheckPermissions CheckPermissions
func (srv *PermServiceImpl) checkPermissions(ctx context.Context, req *rbac.CheckPermissionsRequest) (*rbac.PermCheckResponse, error) {
	return &rbac.PermCheckResponse{Pass: srv.casDao.CheckPermissions(req)}, nil
}

// ListObjectPermissions ListObjectPermissions

func (srv *PermServiceImpl) ListObjectPermissions(ctx context.Context, objectID *rbac.ObjectID) ([]*model.Resource, error) {
	if err := srv.objectIDAllExists(ctx, objectID); err != nil {
		return nil, err
	}

	// todo 是否开启三元管理
	//auditConfig, _ := client.GetInstance().Audit.Approve`Management.GetAuditConfig(ctx, &empty.Empty{})
	//if !auditConfig.EnableThreeMembers {
	//	if err := srv.havePermOrSelfInfo(ctx, objectID); err != nil {
	//		return nil, err
	//	}
	//}
	// 查询用户具有的所有权限id
	resIds := srv.casDao.GetObjectPermissions(objectID)
	// 查询对应权限信息
	resList, err := srv.resourceDao.GetCustomResList(ctx, resIds)
	if err != nil {
		return nil, err
	}
	return resList, nil
}

func (srv *PermServiceImpl) CheckApiPermission(ctx context.Context, userID int64, url, method string) (bool, error) {
	// 检查api鉴权开关是否打开
	if !config.Custom.EnableApiAuthorize {
		return true, nil
	}

	// 查询接口是否配置于资源控制表
	redisClient := boot.Middleware.DefaultRedis()
	apiResStr, _ := redisClient.Get(consts.REDIS_KEY_APIRES).Result()
	var perms []*model.Resource
	if apiResStr == "" {
		perms, _ = srv.resourceDao.ListResourceTypePermissions(ctx, common.PermissionResourceTypeApi)

		apiResByte, _ := json.Marshal(perms)

		redisClient.Set(consts.REDIS_KEY_APIRES, apiResByte, 1*time.Hour)

	} else {
		_ = json.Unmarshal([]byte(apiResStr), &perms)
	}

	var hitPermID int64
	for _, perm := range perms {
		if perm.Name == url && perm.Action == method {
			hitPermID = perm.Id
		}
	}

	if hitPermID == 0 {
		return true, nil
	}

	// 校验用户是否具有接口权限
	rsp, err := srv.checkPermissions(ctx, &rbac.CheckPermissionsRequest{
		Id: &rbac.ObjectID{
			Id:   snowflake.ID(userID).String(),
			Type: rbac.ObjectType_USER,
		},
		PermissionIds: []int64{hitPermID},
	})

	return rsp.Pass, err
}

func (srv PermServiceImpl) GetObjectsByResource(ctx context.Context, req *rbac.ResourceID) (*rbac.ObjectIDs, error) {
	objs, err := srv.casDao.GetObjectsForPermission(req.Id)

	if err != nil {
		return nil, err
	}

	return &rbac.ObjectIDs{
		Ids: objs,
	}, nil
}

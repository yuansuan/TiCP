package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	approveconf "github.com/yuansuan/ticp/PSP/psp/internal/approve/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/jwt"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/ptype"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/config"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type UserServiceImpl struct {
	userDao         dao.UserDao
	certDao         dao.CertificateDao
	orgStructureDao dao.OrgStructureDao
}

func (srv UserServiceImpl) ResetPassword(ctx context.Context, userID string) (string, error) {
	randPass := util.GenerateRandomPassword(consts.PasswordLength)

	user, exists, err := srv.userDao.Get(snowflake.MustParseString(userID).Int64())

	if err != nil {
		return "", status.Error(errcode.ErrInternalServer, err.Error())
	}
	if !exists {
		return "", status.Error(errcode.ErrUserNotFound, "user not found")
	}
	user.Password = util.PasswdCrypto(randPass)

	err = srv.userDao.Update(user)
	if err != nil {
		logging.Default().Errorf("Password could not be changed for user: %v , error: %v", userID, err.Error())
		return "", status.Errorf(errcode.ErrUserResetPass, err.Error())
	}
	//tracelog.Info(ctx, fmt.Sprintf("reset user password success, userID[%v]", userID))

	return randPass, nil
}

func UnhandledApprove(ctx context.Context, userId int64) error {
	// 如果开启三元管理，检查下用户是否有未处理的审批
	if approveconf.GetConfig().ThreePersonManagement {
		rsp, err := client.GetInstance().Approve.CheckUnhandledApprove(ctx, &approve.CheckUnhandledApproveRequest{
			UserId: snowflake.ID(userId).String(),
		})
		if err != nil {
			return err
		}
		if rsp.Unhandled {
			return status.Error(errcode.ErrApproveUnhandledExist, "")
		}
	}
	return nil
}

func (srv UserServiceImpl) GetUserRoleNames(ctx context.Context, userId string) string {
	userDetail, err := srv.Detail(ctx, snowflake.MustParseString(userId).Int64())
	if err != nil {
		logging.GetLogger(ctx).Errorf("get user detial error, err:%v", err)
		return ""
	}
	roleNames := lo.Map[*dto.Role, string](userDetail.RoleInfo, func(role *dto.Role, _ int) string {
		return role.Name
	})

	return strings.Join(roleNames, ",")
}

func (srv UserServiceImpl) GetAllUser(ctx context.Context) ([]*model.User, error) {
	return srv.userDao.ListAllUser()
}

func (srv UserServiceImpl) OptionList(ctx *gin.Context, filterName string, filterPerm int64) (userOptionList []*dto.UserOptionResponse, err error) {
	filterIds := make([]int64, 0)
	if filterPerm > 0 {
		userObj, err := client.GetInstance().Perm.GetObjectsByResource(ctx, &rbac.ResourceID{
			Id: filterPerm,
		})

		if err != nil {
			return nil, err
		}

		if userObj.Ids != nil {
			for _, IdObj := range userObj.Ids {
				if IdObj.Type == rbac.ObjectType_USER {
					filterIds = append(filterIds, snowflake.MustParseString(IdObj.Id).Int64())
				}
			}
		}
	}

	userList, err := srv.userDao.ListUserLikeName(filterName, filterIds)

	if err != nil {
		return nil, status.Error(errcode.ErrUserOptionList, err.Error())
	}

	for _, user := range userList {
		userOptionList = append(userOptionList, &dto.UserOptionResponse{
			Key:   snowflake.ID(user.Id).String(),
			Title: user.Name,
		})
	}

	return
}

func (srv UserServiceImpl) Detail(ctx context.Context, userId int64) (*dto.UserDetailResponse, error) {

	userInfo, err := srv.Get(ctx, userId)

	if err != nil {
		return nil, status.Error(errcode.ErrUserGetFailed, err.Error())
	}

	rsp, err := client.GetInstance().Role.ListObjectsRoles(
		ctx,
		&rbac.ListObjectsRolesReq{
			Ids: []*rbac.ObjectID{
				{
					Id:   snowflake.ID(userId).String(),
					Type: rbac.ObjectType_USER,
				},
			},
			NeedImplicitRoles: false,
		},
	)

	if err != nil {
		return nil, err
	}

	var roleInfo []*dto.Role
	if rsp != nil && len(rsp.ObjectRolesList) > 0 && len(rsp.ObjectRolesList[0].Roles) > 0 {
		roles := rsp.ObjectRolesList[0].Roles

		// 将 Roles 转换为 map
		rolesMap := make(map[int64]string)
		for _, role := range rsp.Roles {
			rolesMap[role.Id] = role.Name
		}

		for _, rID := range roles {
			roleInfo = append(roleInfo, &dto.Role{
				Id:   rID,
				Name: rolesMap[rID],
			})
		}
	}

	myPerm, err := client.GetInstance().Perm.ListObjectPermissions(ctx, &rbac.ObjectID{
		Id:   snowflake.ID(userId).String(),
		Type: rbac.ObjectType_USER,
	})

	if err != nil {
		return nil, err
	}

	// 角色所持有的权限
	hasPermId := make(map[int64]bool, 0)

	if myPerm != nil && len(myPerm.Perms) > 0 {
		for _, perm := range myPerm.Perms {
			hasPermId[perm.Id] = true
		}
	}

	system, localApp, cloudApp, visualSoftware := make([]*dto.CustomPerm, 0), make([]*dto.CustomPerm, 0), make([]*dto.CustomPerm, 0), make([]*dto.CustomPerm, 0)

	permission, err := client.GetInstance().Perm.ListPermission(ctx, &rbac.ListQuery{})
	if permission != nil && permission.Total > 0 {
		for _, perm := range permission.Perms {
			customPerm := &dto.CustomPerm{
				Id:         perm.Id,
				Key:        perm.ResourceName,
				Name:       perm.DisplayName,
				ExternalId: snowflake.ID(perm.ResourceId).String(),
				Has:        hasPermId[perm.Id],
			}

			switch perm.ResourceType {
			case common.PermissionResourceTypeSystem:
				system = append(system, customPerm)
				break
			case common.PermissionResourceTypeLocalApp:
				localApp = append(localApp, customPerm)
				break
			case common.PermissionResourceTypeVisualSoftware:
				visualSoftware = append(visualSoftware, customPerm)
				break
			case common.PermissionResourceTypeAppCloudApp:
				cloudApp = append(cloudApp, customPerm)
				break
			default:
				break
			}
		}
	}

	var openapiCertificate string
	if config.GetConfig().Openapi.Enable {
		cert, exist, err := srv.certDao.GetByUserID(snowflake.ID(userInfo.Id))
		if err != nil {
			return nil, err
		}
		if exist {
			openapiCertificate = cert.Certificate
		}
	}

	return &dto.UserDetailResponse{
		UserInfo: util.ToDTO.UserObj(userInfo, nil),
		RoleInfo: roleInfo,
		Perm: &dto.Perm{
			LocalApp:       localApp,
			CloudApp:       cloudApp,
			VisualSoftware: visualSoftware,
			System:         system,
		},
		Conf: &dto.Conf{
			LdapEnable:    config.GetConfig().LdapConf.Enable,
			OpenapiSwitch: config.GetConfig().Openapi.Enable,
		},
		OpenapiCertificate: openapiCertificate,
	}, nil
}

func (srv UserServiceImpl) QueryUserRole(ctx *gin.Context, req dto.QueryByCondRequest) (*dto.UserListResponse, error) {
	// 查询用户
	users, total, err := srv.QueryByCond(ctx, req)

	if total == 0 || err != nil {
		return &dto.UserListResponse{
			Success: true,
			Total:   total,
			UserObj: make([]*dto.UserInfo, 0),
		}, nil
	}

	// 查询用户对应角色
	var ids []*rbac.ObjectID

	if len(users) > 0 {
		for _, user := range users {
			ids = append(ids, &rbac.ObjectID{
				Id:   snowflake.ID(user.Id).String(),
				Type: rbac.ObjectType_USER,
			})
		}
	}
	roles, err := client.GetInstance().Role.ListObjectsRoles(ctx, &rbac.ListObjectsRolesReq{
		Ids: ids,
	})

	if err != nil {
		return nil, status.Error(errcode.ErrUserQueryFailed, err.Error())
	}

	var userList []*dto.UserInfo

	// 将 ObjectRolesList 转换为 map key-用户id value-用户角色id
	rolesMap := make(map[int64][]int64)
	if roles != nil && len(roles.ObjectRolesList) > 0 {
		for _, userRole := range roles.ObjectRolesList {
			rolesMap[snowflake.MustParseString(userRole.Id.Id).Int64()] = userRole.Roles
		}
	}
	// 匹配用户id，
	if len(rolesMap) > 0 {
		for _, user := range users {
			userList = append(userList, util.ToDTO.UserObj(user, rolesMap[user.Id]))
		}
	}

	return &dto.UserListResponse{
		Success: true,
		Total:   total,
		UserObj: userList,
	}, nil
}

func (srv UserServiceImpl) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	if len(name) == 0 {
		return nil, status.Error(errcode.ErrUserNameEmpty, "name is empty")
	}

	user, err := srv.userDao.GetUserByName(name)
	if err != nil {
		logging.Default().Error(err)
		return nil, err
	}
	return user, nil
}

func (srv UserServiceImpl) UpdatePassword(ctx context.Context, in dto.UpdatePassRequest) error {

	id, err := srv.userDao.GetIdByName(in.Name)
	notFoundErr := status.Error(errcode.ErrUserNotFound, "user not found")
	if err != nil || id < 0 {
		return notFoundErr
	}
	user, exists, err := srv.userDao.Get(id)
	if err != nil {
		return status.Error(errcode.ErrInternalServer, err.Error())
	}
	if !exists {
		return notFoundErr
	}

	//if user.IsDeleted {
	//	return status.Error(errcode.ErrUserDeleted, "user deleted")
	//}

	if util.PasswdCrypto(in.Password) != user.Password {
		return status.Error(errcode.ErrUserInvalidOldPassword, "invalid password")
	}

	user.Password = util.PasswdCrypto(in.NewPassword)
	err = srv.userDao.Update(user)
	if err != nil {
		logging.Default().Errorf("Password could not be changed for user: %v , error: %v", in.Name, err.Error())
		return status.Error(errcode.ErrUserUpdatePasswordFailed, err.Error())
	}
	tracelog.Info(ctx, fmt.Sprintf("update user password success, userID[%v]", id))

	return nil
}

func (srv UserServiceImpl) ActiveUser(ctx context.Context, userId int64) error {
	return srv.setUserActive(ctx, userId, true)
}

// InactiveUser inActive user
func (srv UserServiceImpl) InactiveUser(ctx context.Context, userId int64) error {
	userInfo, err := srv.Get(ctx, userId)
	if err != nil {
		return status.Error(errcode.ErrUserNotExist, "")
	}

	// 如果开启三元管理，检查下是否被设置为默认审批人
	if approveconf.GetConfig().ThreePersonManagement {
		rsp, err := client.GetInstance().SysConfig.GetThreePersonDefaultUserId(ctx, &sysconfig.GetThreePersonDefaultUserIdRequest{})
		if err != nil {
			return err
		}
		if rsp != nil && rsp.UserId == userId {
			return status.Error(errcode.ErrUserCantDeleteDefaultUser, "")
		}

		err = UnhandledApprove(ctx, userId)
		if err != nil {
			return err
		}
	}

	err = srv.setUserActive(ctx, userId, false)
	if err != nil {
		return err
	}

	err = srv.setEnableOpenapi(ctx, userId, false)
	if err != nil {
		return err
	}

	// 清除被禁用用户的token
	jwt.CleanWhiteListByUserName(userInfo.Name)
	return nil
}

func (srv UserServiceImpl) setEnableOpenapi(ctx context.Context, userID int64, enable bool) error {
	user, exists, err := srv.userDao.Get(userID)
	if err != nil {
		logging.Default().Error(err)
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	user.EnableOpenapi = enable
	err = srv.updateUser(user)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update user enable openapi success, userID[%v], enable[%v]", userID, enable))
	return nil
}

func (srv UserServiceImpl) setUserActive(ctx context.Context, id int64, active bool) error {
	user, exists, err := srv.userDao.Get(id)
	if err != nil {
		logging.Default().Error(err)
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	if user.IsInternal && !active {
		return status.Errorf(errcode.ErrUserCantRemoveInternalUser, "internal user %srv is not allowed to be deleted", user.Name)
	}

	user.Enabled = active
	err = srv.updateUser(user)
	if err != nil {
		return err
	}

	tracelog.Info(ctx, fmt.Sprintf("update user enable state success, userID[%v], enable[%v]", id, active))
	return nil
}

func (srv UserServiceImpl) Delete(ctx context.Context, userId int64) error {
	user, exists, err := srv.userDao.Get(userId)
	if err != nil {
		logging.Default().Error(err)
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	if user.IsInternal {
		logging.Default().Errorf("Internal user %v is not allowed to be deleted", user.Name)
		return status.Errorf(errcode.ErrUserCantRemoveInternalUser, "internal user %srv is not allowed to be deleted", user.Name)
	}

	// 如果开启三元管理，检查下是否被设置为默认审批人
	if approveconf.GetConfig().ThreePersonManagement {
		rsp, err := client.GetInstance().SysConfig.GetThreePersonDefaultUserId(ctx, &sysconfig.GetThreePersonDefaultUserIdRequest{})
		if err != nil {
			return err
		}
		if rsp != nil && rsp.UserId == userId {
			return status.Error(errcode.ErrUserCantDeleteDefaultUser, "")
		}

		err = UnhandledApprove(ctx, userId)
		if err != nil {
			return err
		}
	}

	// 如开启openapi则删除用户对应openapi凭证
	if config.GetConfig().Openapi.Enable {
		certModel, exist, err := srv.certDao.GetByUserID(snowflake.ID(userId))
		if err != nil {
			return err
		}
		if exist {
			srv.loseEfficacyCert(certModel.Certificate, snowflake.ID(userId))
			err = srv.certDao.DelByUserID(snowflake.ID(userId))
			if err != nil {
				return err
			}
		}
	}

	err = srv.userDao.Delete(userId)

	if err != nil {
		return err
	}
	// 清除被禁用用户的token
	jwt.CleanWhiteListByUserName(user.Name)

	tracelog.Info(ctx, fmt.Sprintf("delete user success, userID[%v]", userId))

	return nil
}

func (srv UserServiceImpl) Update(ctx context.Context, updateUser model.User) error {
	err := UnhandledApprove(ctx, updateUser.Id)
	if err != nil {
		return err
	}

	user, exists, err := srv.userDao.Get(updateUser.Id)
	if err != nil {
		logging.Default().Error(err)
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	if !user.Enabled && updateUser.EnableOpenapi {
		return status.Error(errcode.ErrUserOpenapiCertNotEnable, "")
	}

	// 从开启到禁用openapi需要删除缓存
	if user.EnableOpenapi && !updateUser.EnableOpenapi {
		srv.loseEfficacyCert("", snowflake.ID(updateUser.Id))
	}

	user.Mobile = updateUser.Mobile
	user.Email = updateUser.Email
	user.EnableOpenapi = updateUser.EnableOpenapi
	err = srv.updateUser(user)
	if err != nil {
		return err
	}

	if config.GetConfig().Openapi.Enable && updateUser.EnableOpenapi {
		srv.GenOpenapiCertificate(ctx, snowflake.ID(updateUser.Id).String(), false)
	}

	return nil
}

func (srv UserServiceImpl) updateUser(user model.User) error {
	return srv.userDao.Update(user)
}

func (srv UserServiceImpl) BatchGetUser(ctx context.Context, userIds []int64) ([]*model.User, error) {
	return srv.userDao.BatchGet(userIds)
}

func (srv UserServiceImpl) ListAll(ctx context.Context, page *ptype.PageReq) ([]*model.User, int64, error) {
	return srv.userDao.ListAll(page.Index, page.Size)
}

func (srv UserServiceImpl) QueryByCond(ctx *gin.Context, in dto.QueryByCondRequest) ([]*model.User, int64, error) {
	// 非内部用户无法看到内部用户
	userID := ginutil.GetUserID(ctx)
	user, exist, _ := srv.userDao.Get(userID)
	var isInternal bool
	if exist {
		isInternal = user.IsInternal
	}

	return srv.userDao.List(in.Enabled, in.Query, in.Order, in.Desc, in.Page.Index, in.Page.Size, isInternal)
}

func (srv UserServiceImpl) Exist(ctx context.Context, userId int64) (bool, error) {
	exit := srv.userDao.Exist(userId)
	return exit, nil
}

func (srv UserServiceImpl) Add(ctx context.Context, user model.User) (int64, error) {
	time.Now().UnixNano()
	logger := logging.GetLogger(ctx)

	user.Enabled = true

	id, err := srv.userDao.Add(user)
	user.Id = id

	if err != nil {
		return 0, status.Error(errcode.ErrUserAddFailed, err.Error())
	}

	logger.Info("add user success ", user)
	tracelog.Info(ctx, fmt.Sprintf("add user success, user[%+v]", user))

	return id, nil
}

func (srv UserServiceImpl) GetIdByName(ctx context.Context, name string) (int64, error) {
	if len(name) == 0 {
		return 0, status.Error(errcode.ErrUserNameEmpty, "name is empty")
	}

	id, err := srv.userDao.GetIdByName(name)
	if err != nil {
		logging.Default().Error(err)
		return 0, status.Error(errcode.ErrInternalServer, err.Error())
	}
	return id, nil
}

func (srv UserServiceImpl) LoginCheck(ctx context.Context, userId int64, userName string) (*dto.LoginSuccessResponse, error) {

	identity, err := srv.GetIdByName(ctx, userName)
	if err != nil {
		return nil, err
	}
	userObj, err := srv.Get(ctx, identity)
	if err != nil {
		return nil, err
	}
	if !userObj.Enabled {
		return nil, status.Error(errcode.ErrAuthUserDisabled, "")
	}

	return &dto.LoginSuccessResponse{User: util.ToDTO.UserObj(userObj, nil)}, nil
}

// Get get user
func (srv UserServiceImpl) Get(ctx context.Context, id int64) (*model.User, error) {

	user, exists, err := srv.userDao.Get(id)
	if err != nil {
		logging.Default().Error(err)
		return nil, status.Error(errcode.ErrInternalServer, err.Error())
	}
	if !exists {
		return nil, status.Error(errcode.ErrUserNotExist, err.Error())
	}

	return &user, nil
}

// GetIncludeDeleted get user
func (srv UserServiceImpl) GetIncludeDeleted(ctx context.Context, id int64) (*model.User, error) {
	user, exists, err := srv.userDao.GetIncludeDeleted(id)
	if err != nil {
		logging.Default().Error(err)
		return nil, status.Error(errcode.ErrInternalServer, err.Error())
	}
	if !exists {
		return nil, status.Error(errcode.ErrUserNotExist, err.Error())
	}

	return &user, nil
}

func (srv UserServiceImpl) AddUserWithRole(ctx context.Context, req dto.UserAddRequest) (int64, error) {

	if req.Name == common.ProjectFolderPath || req.Name == common.PublicFolderPath {
		return 0, status.Error(errcode.ErrUserNameCollideBuiltin, "")
	}

	oldUser, _ := srv.userDao.GetUserByName(req.Name)
	if oldUser != nil {
		return 0, status.Error(errcode.ErrUserNameExist, "")
	}

	userId, err := srv.Add(ctx, model.User{
		Name:          req.Name,
		Password:      req.Password,
		Email:         req.Email,
		Mobile:        req.Mobile,
		RealName:      req.RealName,
		EnableOpenapi: req.EnableOpenapi,
	})

	if err != nil {
		return 0, err
	}

	if len(req.Roles) > 0 {
		_, err = client.GetInstance().Role.AddObjectRoles(ctx, &rbac.ObjectRoles{
			Id: &rbac.ObjectID{
				Id:   snowflake.ID(userId).String(),
				Type: rbac.ObjectType_USER,
			},
			Roles: req.Roles,
		})
		if err != nil {
			return 0, status.Error(errcode.ErrUserAddFailed, err.Error())
		}
	}

	//_, _ = client.GetInstance().Storage.InitUserHome(ctx, &storage.InitUserHomeReq{
	//	UserName: req.Name,
	//})

	logging.Default().Infof("配置开关:{%v}", config.GetConfig().Openapi.Enable)

	if config.GetConfig().Openapi.Enable && req.EnableOpenapi {
		srv.GenOpenapiCertificate(ctx, snowflake.ID(userId).String(), false)
	}

	return userId, nil
}

func (srv UserServiceImpl) GenOpenapiCertificate(ctx context.Context, userId string, over bool) (string, error) {

	id := snowflake.MustParseString(userId)
	userInfo, err := srv.Get(ctx, id.Int64())
	if err != nil {
		return "", err
	}

	// 未开启openapi访问权限
	if !userInfo.EnableOpenapi {
		return "", status.Error(errcode.ErrUserOpenapiCertDisable, "")
	}

	certModel, exist, err := srv.certDao.GetByUserID(id)
	if err != nil {
		return "", err
	}
	if exist {
		if over {
			// 覆盖则删除缓存
			srv.loseEfficacyCert(certModel.Certificate, id)
		} else {
			return certModel.Certificate, nil
		}
	}

	err = srv.certDao.DelByUserID(id)
	if err != nil {
		return "", err
	}

	node, err := snowflake.GetInstance()
	if err != nil {
		return "", err
	}

	cert := node.Generate().String()
	_, err = srv.certDao.Add(&model.OpenapiUserCertificate{
		UserId:      id,
		Certificate: cert,
	})
	logging.Default().Infof("cert:{%s}, err:{%v}", cert, err)

	return cert, err
}

func (srv UserServiceImpl) CheckOpenapiCertificate(ctx context.Context, certificate string) (user *model.User, exist bool, err error) {
	redisClient := boot.Middleware.DefaultRedis()
	userResStr, err := redisClient.Get(GetCertKey(certificate)).Result()

	user = &model.User{}
	if userResStr == "" {
		user, exist, err = srv.certDao.CheckCert(certificate)
		if err != nil || !exist {
			return nil, false, err
		}

		userResByte, _ := json.Marshal(user)
		redisClient.Set(certificate, userResByte, 1*time.Hour)
	} else {
		err = json.Unmarshal([]byte(userResStr), &user)
	}

	return
}

func (srv UserServiceImpl) loseEfficacyCert(certificate string, userID snowflake.ID) {
	if strutil.IsEmpty(certificate) {
		certModel, exist, _ := srv.certDao.GetByUserID(userID)
		if !exist {
			return
		}
		certificate = certModel.Certificate
	}

	redisClient := boot.Middleware.DefaultRedis()
	if redisClient != nil {
		redisClient.Del(GetCertKey(certificate))
	}
}

// GetCertKey 获取cert缓存
func GetCertKey(cert string) string {
	return fmt.Sprintf("PSP:%s:%s:%s", "OPENAPI-CERT", cert)
}

func NewUserService() service.UserService {
	return &UserServiceImpl{
		userDao:         dao.NewUserDaoImpl(),
		certDao:         dao.NewCertificateDaoImpl(),
		orgStructureDao: dao.NewOrgStructureDaoImpl(),
	}
}

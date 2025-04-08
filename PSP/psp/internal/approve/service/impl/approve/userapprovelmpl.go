package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

type UserApproveImpl[T any] struct {
	approveDao dao.ApproveDao
}

type UserApproveInfo struct {
	ApproveType     dto.ApproveType `json:"approve_type"` // 审批类型
	OperateUserInfo OperateUserInfo `json:"operate_user_info"`
	Id              int64           `json:"id"`             // 用户id
	Name            string          `json:"name"`           // 用户名
	Password        string          `json:"password"`       // 密码
	Email           string          `json:"email"`          // 邮箱
	Mobile          string          `json:"mobile"`         // 手机号
	RealName        string          `json:"real_name"`      // 真实姓名
	Roles           []int64         `json:"roles"`          // 所属角色信息
	EnableOpenapi   bool            `json:"enable_openapi"` // 是否开启openapi
}

type OperateUserInfo struct {
	UserID    int64  `json:"user_id"`    // 用户id
	UserName  string `json:"user_name"`  // 用户名称
	IpAddress string `json:"ip_address"` // 用户请求ip
}

func NewUserApproveImpl() *UserApproveImpl[any] {
	return &UserApproveImpl[any]{
		approveDao: dao.NewApproveDao(),
	}
}

func (srv *UserApproveImpl[T]) PreCheck(ctx context.Context, req *dto.ApplyApproveRequest) error {
	if strutil.IsNotEmpty(req.UserApproveInfoRequest.Id) && req.UserApproveInfoRequest.Id == req.ApproveUserID {
		// 审批人不能审批自己
		return status.Error(errcode.ErrUserApproveSelf, "")
	}

	switch req.ApproveType {
	// 新增情况检查用户名是否已存在
	case dto.ApproveTypeAddUser:
		rsp, err := client.GetInstance().User.GetUserByName(ctx, &user.NameCondRequest{
			Name: req.UserApproveInfoRequest.Name,
		})
		if err != nil && errcode.ErrUserNotFound != status.Code(err) {
			return err
		}
		if rsp != nil && rsp.Name == req.UserApproveInfoRequest.Name {
			return status.Error(errcode.ErrUserNameExist, "")
		}
		break
	// 删除、禁用情况检查下是否被设置为默认审批人
	case dto.ApproveTypeDelUser, dto.ApproveTypeDisableUser:
		rsp, err := client.GetInstance().SysConfig.GetThreePersonDefaultUserId(ctx, &sysconfig.GetThreePersonDefaultUserIdRequest{})
		if err != nil {
			logging.Default().Errorf("=============rsp:%+v, err:%v", rsp, err)
			return err
		}
		if rsp != nil && rsp.UserId == snowflake.MustParseString(req.UserApproveInfoRequest.Id).Int64() {
			return status.Error(errcode.ErrUserCantDeleteDefaultUser, "")
		}
	default:
		break
	}

	return nil
}

func (srv *UserApproveImpl[T]) BuildContent(ctx context.Context, approveInfo any) (string, error) {
	req := approveInfo.(*UserApproveInfo)
	var content string

	oldUser, err := getOldUser(ctx, req.ApproveType, req.Id)
	if err != nil {
		return "", err
	}

	switch req.ApproveType {
	case dto.ApproveTypeAddUser:
		roleNames, err := srv.GetRoleNames(ctx, req.Roles)
		if err != nil {
			return "", err
		}
		content = fmt.Sprintf("新建用户[%s]，用户角色为[%s]", req.Name, roleNames)
		break
	case dto.ApproveTypeDelUser:
		content = fmt.Sprintf("删除用户[%s]", oldUser.Name)
		break
	case dto.ApproveTypeEditUser:
		roleNames, err := srv.GetRoleNames(ctx, req.Roles)
		if err != nil {
			return "", err
		}
		content = fmt.Sprintf("编辑用户[%s]，用户角色更改为[%s]", req.Name, roleNames)
		break
	case dto.ApproveTypeEnableUser:
		content = fmt.Sprintf("启用用户[%s]", oldUser.Name)
		break
	case dto.ApproveTypeDisableUser:
		content = fmt.Sprintf("禁用用户[%s]", oldUser.Name)
		break
	}

	return content, nil
}

func (srv *UserApproveImpl[T]) GetRoleNames(ctx context.Context, roleIds []int64) (string, error) {
	var roleNameList []string
	var roleNames string
	rsp, err := client.GetInstance().Rbac.GetRoles(ctx, &rbac.RoleIDs{
		Ids: roleIds,
	})
	if err != nil {
		return "", err
	}

	for _, role := range rsp.GetRoles() {
		roleNameList = append(roleNameList, role.Name)
	}
	if len(roleNameList) > 0 {
		roleNames = strings.Join(roleNameList, ",")
	}
	return roleNames, nil
}

func (srv *UserApproveImpl[T]) BuildApproveInfo(ctx *gin.Context, req *dto.ApplyApproveRequest) (approveInfo any, err error) {
	userReq := req.UserApproveInfoRequest
	return &UserApproveInfo{
		ApproveType: req.ApproveType,
		OperateUserInfo: OperateUserInfo{
			UserID:    ginutil.GetUserID(ctx),
			UserName:  ginutil.GetUserName(ctx),
			IpAddress: ctx.ClientIP(),
		},
		Id:            snowflake.MustParseString(userReq.Id).Int64(),
		Name:          userReq.Name,
		Password:      userReq.Password,
		Email:         userReq.Email,
		Mobile:        userReq.Mobile,
		RealName:      userReq.RealName,
		Roles:         userReq.Roles,
		EnableOpenapi: userReq.EnableOpenapi,
	}, nil
}

func (srv *UserApproveImpl[T]) GenSign(ctx context.Context, req *dto.ApplyApproveRequest) (string, error) {
	var sign string

	switch req.ApproveType {
	case dto.ApproveTypeAddUser:
		sign = fmt.Sprintf("%s-%s-%s", UserApproveType, SignTypeName, req.UserApproveInfoRequest.Name)
		break
	default:
		sign = fmt.Sprintf("%s-%s-%v", UserApproveType, SignTypeID, snowflake.MustParseString(req.UserApproveInfoRequest.Id).String())
		break
	}
	return sign, nil
}

func (srv *UserApproveImpl[T]) CheckNecessary(ctx context.Context, approveType dto.ApproveType, approveInfo any) (bool, error, error) {
	req := approveInfo.(*UserApproveInfo)

	switch approveType {
	// 新增修改用户先检查给用户分配的角色是否已经不存在了
	case dto.ApproveTypeAddUser, dto.ApproveTypeEditUser:
		if len(req.Roles) == 0 {
			break
		}
		_, err := client.GetInstance().Rbac.GetRoles(ctx, &rbac.RoleIDs{
			Ids: req.Roles,
		})
		if err != nil && status.Code(err) == errcode.ErrRBACRoleNotFound {
			return false, status.Error(errcode.ErrApproveNecessaryRoleNotExist, ""), nil
		}
		break
	default:
		break
	}

	return true, nil, nil
}

func (srv *UserApproveImpl[T]) AfterPass(ctx context.Context, approveInfo any) error {
	req := approveInfo.(*UserApproveInfo)
	operateUser := req.OperateUserInfo

	oldUser, err := getOldUser(ctx, req.ApproveType, req.Id)
	if err != nil {
		return err
	}

	switch req.ApproveType {
	case dto.ApproveTypeAddUser:
		_, err = client.GetInstance().User.AddUserWithRole(ctx, &user.AddUserWithRoleRequest{
			Name:          req.Name,
			Password:      req.Password,
			Email:         req.Email,
			Mobile:        req.Mobile,
			RealName:      req.RealName,
			RoleIds:       req.Roles,
			EnableOpenapi: req.EnableOpenapi,
		})
		if err != nil {
			return err
		}

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_USER_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v新建用户[%v]", operateUser.UserName, req.Name))
		break
	case dto.ApproveTypeEditUser:

		sourceRoleNames, _ := client.GetInstance().User.GetUserRoleNames(ctx, &user.UserIdentity{Id: snowflake.ID(req.Id).String()})
		_, err = client.GetInstance().User.UpdateUser(ctx, &user.UpdateUserRequest{
			Id:            snowflake.ID(req.Id).String(),
			Email:         req.Email,
			Mobile:        req.Mobile,
			RoleIds:       req.Roles,
			EnableOpenapi: req.EnableOpenapi,
		})
		if err != nil {
			return err
		}
		targetRoleNames, _ := client.GetInstance().User.GetUserRoleNames(ctx, &user.UserIdentity{Id: snowflake.ID(req.Id).String()})

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_USER_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v修改用户[%v]【%v -》%v】", operateUser.UserName, req.Name, sourceRoleNames.RoleNames, targetRoleNames.RoleNames))
		break
	case dto.ApproveTypeDelUser:
		_, err = client.GetInstance().User.DelUser(ctx, &user.UserIdentity{
			Id: snowflake.ID(req.Id).String(),
		})
		if err != nil {
			return err
		}

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_USER_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v删除用户[%v]", operateUser.UserName, oldUser.Name))
		break
	case dto.ApproveTypeEnableUser:
		_, err = client.GetInstance().User.EnableUser(ctx, &user.EnableUserRequest{
			Id:     snowflake.ID(req.Id).String(),
			Enable: true,
		})
		if err != nil {
			return err
		}

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_USER_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v启用用户[%v]", operateUser.UserName, oldUser.Name))
		break
	case dto.ApproveTypeDisableUser:
		_, err = client.GetInstance().User.EnableUser(ctx, &user.EnableUserRequest{
			Id:     snowflake.ID(req.Id).String(),
			Enable: false,
		})
		if err != nil {
			return err
		}

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_USER_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v禁用用户[%v]", operateUser.UserName, oldUser.Name))
		break
	}

	return nil
}

func getOldUser(ctx context.Context, approveType dto.ApproveType, userID int64) (*user.UserObj, error) {
	switch approveType {
	case dto.ApproveTypeDelUser, dto.ApproveTypeEditUser, dto.ApproveTypeEnableUser, dto.ApproveTypeDisableUser:
		return client.GetInstance().User.Get(ctx, &user.UserIdentity{
			Id: snowflake.ID(userID).String(),
		})
	default:
		break
	}
	return nil, nil
}

func (srv *UserApproveImpl[T]) ParseObject(jsonString string) (any, error) {
	var userApproveInfo interface{} = &UserApproveInfo{}
	err := json.Unmarshal([]byte(jsonString), &userApproveInfo)
	if err != nil {
		return userApproveInfo, err
	}

	return userApproveInfo, nil
}

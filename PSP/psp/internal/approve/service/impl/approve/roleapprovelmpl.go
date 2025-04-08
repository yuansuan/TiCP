package approve

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	usesconf "github.com/yuansuan/ticp/PSP/psp/internal/user/config"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

type RoleApproveImpl[T any] struct {
	approveDao dao.ApproveDao
}

type RoleApproveInfo struct {
	ApproveType     dto.ApproveType `json:"approve_type"` // 审批类型
	OperateUserInfo OperateUserInfo `json:"operate_user_info"`
	Id              int64           `json:"id"`      // 角色id
	Name            string          `json:"name"`    // 角色名
	Comment         string          `json:"comment"` // 描述
	Perms           []int64         `json:"perms"`   // 权限id
}

func NewRoleApproveImpl() *RoleApproveImpl[any] {
	return &RoleApproveImpl[any]{
		approveDao: dao.NewApproveDao(),
	}
}

func (srv *RoleApproveImpl[T]) PreCheck(ctx context.Context, req *dto.ApplyApproveRequest) error {
	switch req.ApproveType {
	case dto.ApproveTypeAddRole, dto.ApproveTypeEditRole:
		rsp, err := client.GetInstance().Rbac.GetRoleByName(ctx, &rbac.RoleName{
			Name: req.RoleApproveInfoRequest.Name,
		})
		if err != nil && errcode.ErrRBACRoleNotFound != status.Code(err) {
			return err
		}

		if req.ApproveType == dto.ApproveTypeAddRole {
			// 新增情况检查角色名是否已存在
			if rsp != nil && rsp.Name == req.RoleApproveInfoRequest.Name {
				return status.Error(errcode.ErrRBACRoleNameExist, "")
			}
		} else if req.ApproveType == dto.ApproveTypeEditRole {
			// 修改情况检查除自己以外的角色名是否重复
			if rsp != nil && rsp.Id != req.RoleApproveInfoRequest.Id {
				return status.Error(errcode.ErrRBACRoleNameExist, "")
			}
		}
	case dto.ApproveTypeDelRole:
		if usesconf.GetConfig().LdapConf.Enable {
			rsp, err := client.GetInstance().SysConfig.GetRBACDefaultRoleId(ctx, &sysconfig.GetRBACDefaultRoleIdRequest{})
			if err != nil {
				return err
			}
			if rsp != nil && rsp.RoleId == req.RoleApproveInfoRequest.Id {
				return status.Error(errcode.ErrRBACRoleUsedByLdap, "")
			}
		}
		break
	default:
		break
	}
	return nil
}

func (srv *RoleApproveImpl[T]) BuildContent(ctx context.Context, approveInfo any) (string, error) {
	req := approveInfo.(*RoleApproveInfo)
	var content string

	switch req.ApproveType {
	case dto.ApproveTypeAddRole:
		permNames, err := srv.GetPermNames(ctx, req.Perms)
		if err != nil {
			return "", err
		}
		content = fmt.Sprintf("新建角色[%s]，角色权限为 %s", req.Name, permNames)
		break
	case dto.ApproveTypeDelRole:
		oldRole, err := client.GetInstance().Rbac.GetRole(ctx, &rbac.RoleID{
			Id: req.Id,
		})
		if err != nil {
			return "", err
		}

		content = fmt.Sprintf("删除角色[%s]", oldRole.Name)
		break
	case dto.ApproveTypeEditRole:
		permNames, err := srv.GetPermNames(ctx, req.Perms)
		if err != nil {
			return "", err
		}
		content = fmt.Sprintf("编辑角色[%s]，角色权限更改为 %s", req.Name, permNames)
		break
	case dto.ApproveTypeSetLdapDefRole:
		oldRole, err := client.GetInstance().Rbac.GetRole(ctx, &rbac.RoleID{
			Id: req.Id,
		})
		if err != nil {
			return "", err
		}
		content = fmt.Sprintf("设置角色[%s]为LDAP用户首次登录默认用户", oldRole.Name)
		break
	}

	return content, nil
}
func (srv *RoleApproveImpl[T]) BuildApproveInfo(ctx *gin.Context, req *dto.ApplyApproveRequest) (approveInfo any, err error) {
	roleReq := req.RoleApproveInfoRequest
	return &RoleApproveInfo{
		ApproveType: req.ApproveType,
		OperateUserInfo: OperateUserInfo{
			UserID:    ginutil.GetUserID(ctx),
			UserName:  ginutil.GetUserName(ctx),
			IpAddress: ctx.ClientIP(),
		},
		Id:      roleReq.Id,
		Name:    roleReq.Name,
		Comment: roleReq.Comment,
		Perms:   roleReq.Perms,
	}, nil
}

func (srv *RoleApproveImpl[T]) GenSign(ctx context.Context, req *dto.ApplyApproveRequest) (string, error) {
	var sign string

	switch req.ApproveType {
	case dto.ApproveTypeAddRole:
		sign = fmt.Sprintf("%s-%s-%s", RoleApproveType, SignTypeName, req.RoleApproveInfoRequest.Name)
		break
	default:
		sign = fmt.Sprintf("%s-%s-%v", RoleApproveType, SignTypeID, req.RoleApproveInfoRequest.Id)
		break
	}
	return sign, nil
}

func (srv *RoleApproveImpl[T]) GetPermNames(ctx context.Context, permIds []int64) (string, error) {
	var permNames string
	system, localApp, cloudApp, visualSoftware := make([]string, 0), make([]string, 0), make([]string, 0), make([]string, 0)

	if len(permIds) > 0 {
		rsp, err := client.GetInstance().Perm.GetPermissions(ctx, &rbac.PermissionIDs{
			Ids: permIds,
		})
		if err != nil {
			return "", err
		}
		for _, perm := range rsp.GetPerms() {
			switch perm.ResourceType {
			case common.PermissionResourceTypeSystem:
				system = append(system, perm.DisplayName)
				break
			case common.PermissionResourceTypeLocalApp:
				localApp = append(localApp, perm.DisplayName)
				break
			case common.PermissionResourceTypeVisualSoftware:
				visualSoftware = append(visualSoftware, perm.DisplayName)
				break
			case common.PermissionResourceTypeAppCloudApp:
				cloudApp = append(cloudApp, perm.DisplayName)
				break
			default:
				break
			}
		}
	}

	permNames = fmt.Sprintf("%s %s %s %s",
		permLog(system, "系统权限"), permLog(localApp, "本地应用"),
		permLog(cloudApp, "云应用"), permLog(visualSoftware, "3D可视化镜像"))
	return permNames, nil
}

func permLog(permList []string, typeName string) string {
	return fmt.Sprintf("%s:[%s]", typeName, strings.Join(permList, ","))
}

func (srv *RoleApproveImpl[T]) CheckNecessary(ctx context.Context, approveType dto.ApproveType, approveInfo any) (bool, error, error) {
	req := approveInfo.(*RoleApproveInfo)

	switch approveType {
	// 新增修改角色先检查给角色分配的权限是否已经不存在了
	case dto.ApproveTypeAddRole, dto.ApproveTypeEditRole:
		if len(req.Perms) == 0 {
			break
		}
		_, err := client.GetInstance().Rbac.GetRoles(ctx, &rbac.RoleIDs{
			Ids: req.Perms,
		})
		if err != nil && status.Code(err) == errcode.ErrRBACPermissionNotFound {
			return false, status.Error(errcode.ErrApproveNecessaryPermNotExist, ""), nil
		}

		rsp, err := client.GetInstance().Rbac.GetRoleByName(ctx, &rbac.RoleName{
			Name: req.Name,
		})
		if err != nil && errcode.ErrRBACRoleNotFound != status.Code(err) {
			return false, nil, err
		}

		if req.ApproveType == dto.ApproveTypeAddRole {
			// 新增情况检查角色名是否已存在
			if rsp != nil && rsp.Name == req.Name {
				return false, status.Error(errcode.ErrApproveNecessaryRoleNameExist, ""), nil
			}
		} else if req.ApproveType == dto.ApproveTypeEditRole {
			// 修改情况检查除自己以外的角色名是否重复
			if rsp != nil && rsp.Id != req.Id {
				return false, status.Error(errcode.ErrApproveNecessaryRoleNameExist, ""), nil
			}
		}

		break
	default:
		break
	}

	return true, nil, nil
}

func (srv *RoleApproveImpl[T]) AfterPass(ctx context.Context, approveInfo any) error {
	req := approveInfo.(*RoleApproveInfo)
	operateUser := req.OperateUserInfo

	switch req.ApproveType {
	case dto.ApproveTypeAddRole:
		if _, err := client.GetInstance().Rbac.AddRole(ctx, &rbac.AddRoleReq{
			Name:    req.Name,
			Comment: req.Comment,
			Perms:   req.Perms,
		}); err != nil {
			return err
		}

		permNames, _ := srv.GetPermNames(ctx, req.Perms)
		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_RBAC_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v新增角色[%v], 权限为 %v", operateUser.UserName, req.Name, permNames))
		break
	case dto.ApproveTypeEditRole:
		if _, err := client.GetInstance().Rbac.UpdateRole(ctx, &rbac.UpdateRoleReq{
			Id:      req.Id,
			Name:    req.Name,
			Comment: req.Comment,
			Perms:   req.Perms,
		}); err != nil {
			return err
		}

		permNames, _ := srv.GetPermNames(ctx, req.Perms)
		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_RBAC_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v修改角色[%v], 权限更改为 %v", operateUser.UserName, req.Name, permNames))
		break
	case dto.ApproveTypeDelRole:
		oldRole, err := client.GetInstance().Rbac.GetRole(ctx, &rbac.RoleID{
			Id: req.Id,
		})
		if err != nil {
			return err
		}

		if _, err = client.GetInstance().Rbac.DelRole(ctx, &rbac.RoleID{
			Id: req.Id,
		}); err != nil {
			return err
		}

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_RBAC_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v删除角色[%v]", operateUser.UserName, oldRole.Name))
		break
	case dto.ApproveTypeSetLdapDefRole:
		oldRole, err := client.GetInstance().Rbac.GetRole(ctx, &rbac.RoleID{
			Id: req.Id,
		})
		if err != nil {
			return err
		}

		if _, err = client.GetInstance().Rbac.SetLdapUserDefRole(ctx, &rbac.RoleID{Id: req.Id}); err != nil {
			return err
		}

		oplog.GetInstance().SaveAuditLogInfoGrpc(ctx, approve.OperateTypeEnum_RBAC_MANAGER, snowflake.ID(operateUser.UserID), operateUser.UserName, operateUser.IpAddress,
			fmt.Sprintf("用户%v设置角色[%v]为LDAP用户首次登录默认用户", operateUser.UserName, oldRole.Name))
		break
	}

	return nil
}

func (srv *RoleApproveImpl[T]) ParseObject(jsonString string) (any, error) {
	var roleApproveInfo interface{} = &RoleApproveInfo{}
	err := json.Unmarshal([]byte(jsonString), &roleApproveInfo)
	if err != nil {
		return roleApproveInfo, err
	}

	return roleApproveInfo, nil
}

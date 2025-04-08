package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/sysconfig"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// AddRole
//
//	@Summary		添加自定义角色接口
//	@Description	添加自定义角色接口
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.AddRole	true	"入参"
//	@Success		200		{int}	roleId		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/role/add [post]
func (s *RouteService) AddRole(ctx *gin.Context) {
	var req dto.AddRole

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	role := &model.Role{
		Name:    req.Name,
		Comment: req.Comment,
		Type:    consts.RoleTypeCustom,
	}

	err := s.RoleService.AddRole(ctx, role)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACAddRoleError)
		return
	}

	// Add perms to the role
	if len(req.Perms) > 0 {
		err = s.RoleService.AddRolePerms(ctx, role.Id, req.Perms)
		if err != nil {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACAddRoleError)
			return
		}
	}

	s.saveRolePermAuditLogInfo(ctx, role.Id, "用户%v新增角色[%v], 权限为 %s %s %s %s")

	ginutil.Success(ctx, role.Id)
}

// QueryRole
//
//	@Summary		查询角色接口
//	@Description	查询角色接口
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.ListQueryRequest	true	"入参"
//	@Success		200		{object}	dto.ListQueryResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/role/query [post]
func (s *RouteService) QueryRole(ctx *gin.Context) {
	var req dto.ListQueryRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	roles, total, err := s.RoleService.ListRole(ctx, &boring.ListRequest{
		NameFilter: req.NameFilter,
		Page:       req.Page.Index,
		PageSize:   req.Page.Size,
		Desc:       req.Desc,
		OrderBy:    req.OrderBy,
	})

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACQueryError)
		return
	}

	resp, err := client.GetInstance().SysConfig.GetRBACDefaultRoleId(ctx, &sysconfig.GetRBACDefaultRoleIdRequest{})
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACQueryError)
		return
	}

	ginutil.Success(ctx, &dto.ListQueryResponse{
		Role:  util.ToDTO.RoleInfos(roles, resp.RoleId),
		Total: total,
	})
}

// GetRoleDetail
//
//	@Summary		获取角色详情
//	@Description	获取角色详情
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			id	query		int				true	"角色id"	default(0)
//	@Success		200	{object}	dto.RoleDetail	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/role/detail [get]
func (s *RouteService) GetRoleDetail(ctx *gin.Context) {
	id := ctx.Query("id")

	if strutil.IsEmpty(id) {
		http.Errf(ctx, errcode.ErrInvalidParam, "id can't empty")
		return
	}

	roleId, _ := strconv.Atoi(id)
	roleDetail := &dto.RoleDetail{}
	var err error
	if roleId > 0 {
		roleDetail, err = s.RoleService.GetRoleDetail(ctx, int64(roleId))
		if err != nil {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACGetRole)
			return
		}
	}

	// 查询到所有资源
	resource, total, err := s.PermService.ListResource(ctx, &boring.ListRequest{})
	// 角色所持有的权限
	hasPermId := make(map[int64]bool, 0)

	if len(roleDetail.Resources) > 0 {
		for _, res := range roleDetail.Resources {
			hasPermId[res.Id] = true
		}
	}

	system, localApp, cloudApp, visualSoftware := make([]*dto.CustomPerm, 0), make([]*dto.CustomPerm, 0), make([]*dto.CustomPerm, 0), make([]*dto.CustomPerm, 0)

	if total > 0 {
		for _, res := range resource {

			customPerm := &dto.CustomPerm{
				Id:         res.Id,
				Key:        res.Name,
				Name:       res.DisplayName,
				ExternalId: snowflake.ID(res.ExternalId).String(),
				Has:        hasPermId[res.Id],
			}

			switch res.Type {
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

	roleDetail.Perm = &dto.Perm{
		LocalApp:       localApp,
		CloudApp:       cloudApp,
		VisualSoftware: visualSoftware,
		System:         system,
	}

	ginutil.Success(ctx, roleDetail)
}

// UpdateRole
//
//	@Summary		修改角色接口
//	@Description	修改角色接口
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.RoleInfo	true	"入参"
//	@Success		200		"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/role/update [put]
func (s *RouteService) UpdateRole(ctx *gin.Context) {
	var req *dto.RoleInfo

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}
	err := s.RoleService.UpdateRole(ctx, util.FromDTO.Role(req))

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACUpdateError)
		return
	}

	err = s.RoleService.UpdateRolePerms(ctx, req.Id, req.Perms)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACUpdateError)
		return
	}
	s.saveRolePermAuditLogInfo(ctx, req.Id, "用户%v修改角色[%v], 权限更改为  %s %s %s %s")

	ginutil.Success(ctx, nil)
}

func (s *RouteService) saveRolePermAuditLogInfo(ctx *gin.Context, roleId int64, format string) {
	detail, _ := s.RoleService.GetRoleDetail(ctx, roleId)
	if detail != nil {
		system, localApp, cloudApp, visualSoftware := make([]string, 0), make([]string, 0), make([]string, 0), make([]string, 0)
		resources := detail.Resources
		for _, res := range resources {
			switch res.ResourceType {
			case common.PermissionResourceTypeSystem:
				system = append(system, res.DisplayName)
				break
			case common.PermissionResourceTypeLocalApp:
				localApp = append(localApp, res.DisplayName)
				break
			case common.PermissionResourceTypeVisualSoftware:
				visualSoftware = append(visualSoftware, res.DisplayName)
				break
			case common.PermissionResourceTypeAppCloudApp:
				cloudApp = append(cloudApp, res.DisplayName)
				break
			default:
				break
			}
		}

		oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_RBAC_MANAGER, fmt.Sprintf(format,
			ginutil.GetUserName(ctx), detail.Role.Name, util.PermLog(system, "系统权限"), util.PermLog(localApp, "本地应用"),
			util.PermLog(cloudApp, "云应用"), util.PermLog(visualSoftware, "3D可视化镜像")))
	}
}

// DeleteRole
//
//	@Summary		删除角色接口
//	@Description	删除角色接口
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			id	query	int	true	"角色id"	default(0)
//	@Success		200	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/role/delete [delete]
func (s *RouteService) DeleteRole(ctx *gin.Context) {
	id := ctx.Query("id")

	if strutil.IsEmpty(id) {
		http.Errf(ctx, errcode.ErrInvalidParam, "id can't empty")
		return
	}

	roleId, _ := strconv.Atoi(id)
	roleInfo, err := s.RoleService.GetRole(ctx, int64(roleId))
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACDeleteError)
		return
	}

	err = s.RoleService.DeleteRole(ctx, int64(roleId))
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrRBACDeleteError)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_RBAC_MANAGER, fmt.Sprintf("用户%v删除角色[%v]", ginutil.GetUserName(ctx), roleInfo.Name))

	ginutil.Success(ctx, nil)
}

// SetLdapUserDefRole
//
//	@Summary		修改ldap用户默认角色
//	@Description	修改ldap用户默认角色
//	@Tags			角色
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.LdapDefRoleRequest	true	"入参"
//	@Success		200
//	@Router			/role/setLdapUserDefRole [put]
func (s *RouteService) SetLdapUserDefRole(ctx *gin.Context) {

	var req dto.LdapDefRoleRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if req.ID == 0 {
		http.Errf(ctx, errcode.ErrInvalidParam, "roleID can't empty")
		return
	}

	err := s.RoleService.SetLdapUserDefRole(ctx, req.ID)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserLDAPDefRoleFailed)
		return
	}

	role, err := s.RoleService.GetRole(ctx, req.ID)
	if err != nil {
		oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_RBAC_MANAGER, fmt.Sprintf("用户%v设置角色[%v]为LDAP用户首次登录默认用户", ginutil.GetUserName(ctx), role.Name))
	}

	ginutil.Success(ctx, nil)
}

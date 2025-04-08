package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/jwt"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

// AddUser
//
//	@Summary		添加用户接口
//	@Description	添加用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.UserAddRequest	true	"入参"
//	@Success		200		{string}	userID
//	@Router			/user/add [post]
func (s *RouteService) AddUser(ctx *gin.Context) {

	var req dto.UserAddRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	userId, err := s.UserService.AddUserWithRole(ctx, req)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserAddFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_USER_MANAGER, fmt.Sprintf("用户%v新建用户[%v]", ginutil.GetUserName(ctx), req.Name))

	ginutil.Success(ctx, snowflake.ID(userId).String())
}

// Query
//
//	@Summary		查询用户接口
//	@Description	查询用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body		dto.QueryByCondRequest	true	"入参"
//	@Success		200		{object}	dto.UserListResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/user/query [post]
func (s *RouteService) Query(ctx *gin.Context) {

	var req dto.QueryByCondRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	userList, err := s.UserService.QueryUserRole(ctx, req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserGetFailed)
		return
	}

	ginutil.Success(ctx, userList)
}

// Get
//
//	@Summary		查询用户接口
//	@Description	查询用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			id	query		string					true	"用户id"	default("")
//	@Success		200	{object}	dto.UserDetailResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/user/get [get]
func (s *RouteService) Get(ctx *gin.Context) {
	ID := ctx.Query("id")

	if strutil.IsEmpty(ID) {
		errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrInvalidParam)
		return
	}

	userInfo, err := s.UserService.Detail(ctx, snowflake.MustParseString(ID).Int64())
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserGetFailed)
		return
	}

	ginutil.Success(ctx, userInfo)
}

// Current
//
//	@Summary		当前用户详细信息接口
//	@Description	当前用户详细信息接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UserDetailResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/user/current [get]
func (s *RouteService) Current(ctx *gin.Context) {

	userID := ginutil.GetUserID(ctx)
	if userID <= 0 {
		http.Errf(ctx, errcode.ErrInvalidParam, "id can't empty")
		return
	}

	current, err := s.UserService.Detail(ctx, userID)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserGetFailed)
		return
	}

	ginutil.Success(ctx, current)
}

// Active
//
//	@Summary		启用用户接口
//	@Description	启用用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UserIDRequest	true	"入参"
//	@Success		200
//	@Router			/user/active [put]
func (s *RouteService) Active(ctx *gin.Context) {

	var req dto.UserIDRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if strutil.IsEmpty(req.Id) {
		errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrInvalidParam)
		return
	}

	err := s.UserService.ActiveUser(ctx, snowflake.MustParseString(req.Id).Int64())

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserActiveFailed)
		return
	}

	user, _ := s.UserService.Get(ctx, snowflake.MustParseString(req.Id).Int64())
	var userName string
	if user != nil {
		userName = user.Name
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_USER_MANAGER, fmt.Sprintf("用户%v启用用户[%v]", ginutil.GetUserName(ctx), userName))

	ginutil.Success(ctx, nil)
}

// Inactive
//
//	@Summary		禁用用户接口
//	@Description	禁用用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UserIDRequest	true	"入参"
//	@Success		200
//	@Router			/user/inactive [put]
func (s *RouteService) Inactive(ctx *gin.Context) {
	var req dto.UserIDRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	if strutil.IsEmpty(req.Id) {
		errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrInvalidParam)
		return
	}

	userID := snowflake.MustParseString(req.Id).Int64()
	user, err := s.UserService.Get(ctx, userID)

	if err != nil || user == nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserNotExist)
		return
	}

	err = s.UserService.InactiveUser(ctx, userID)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserInActiveFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_USER_MANAGER, fmt.Sprintf("用户%v禁用用户[%v]", ginutil.GetUserName(ctx), user.Name))

	ginutil.Success(ctx, nil)
}

// Delete
//
//	@Summary		删除用户接口
//	@Description	删除用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			id	query	string	true	"用户id"	default("")
//	@Success		200
//	@Router			/user/delete [delete]
func (s *RouteService) Delete(ctx *gin.Context) {
	ID := ctx.Query("id")

	if strutil.IsEmpty(ID) {
		errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrInvalidParam)
		return
	}
	userID := snowflake.MustParseString(ID).Int64()
	user, err := s.UserService.Get(ctx, userID)

	if err != nil || user == nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserNotExist)
		return
	}

	err = s.UserService.Delete(ctx, userID)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserDeleted)
		return
	}

	err = s.OrgService.DeleteOrgMemberByUserID(ctx, []string{ID})
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgMemberDeletedFailed)
		return
	}

	client.GetInstance().Role.RemoveObjectRoles(ctx, &rbac.ObjectRoles{
		Id: &rbac.ObjectID{
			Id:   ID,
			Type: rbac.ObjectType_USER,
		},
	})

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_USER_MANAGER, fmt.Sprintf("用户%v删除用户[%v]", ginutil.GetUserName(ctx), user.Name))

	ginutil.Success(ctx, nil)
}

// UpdatePassword
//
//	@Summary		修改用户密码接口
//	@Description	修改用户密码接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UpdatePassRequest	true	"入参"
//	@Success		200
//	@Router			/user/updatePassword [put]
func (s *RouteService) UpdatePassword(ctx *gin.Context) {

	var req dto.UpdatePassRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	err := s.UserService.UpdatePassword(ctx, req)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserUpdatePasswordFailed)
		return
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_USER_MANAGER, fmt.Sprintf("用户%v修改用户[%v]密码", ginutil.GetUserName(ctx), req.Name))
	// 清除此用户所有生效的token
	jwt.CleanWhiteListByUserName(ginutil.GetUserName(ctx))

	ginutil.Success(ctx, nil)
}

// GetDataConfig
//
//	@Summary		获取用户数据配置
//	@Description	获取用户数据配置
//	@Tags			用户-user
//	@Produce		json
//	@Success		200
//	@Router			/user/getDataConfig [get]
func (s *RouteService) GetDataConfig(ctx *gin.Context) {
	ginutil.Success(ctx, dto.UserDataConfigResponse{
		HomeDir: ".",
	})
}

// Update
//
//	@Summary		修改用户接口
//	@Description	修改用户接口
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.UserUpdateRequest	true	"入参"
//	@Success		200
//	@Router			/user/update [put]
func (s *RouteService) Update(ctx *gin.Context) {

	var req dto.UserUpdateRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	err := s.UserService.Update(ctx, model.User{
		Id:            snowflake.MustParseString(req.Id).Int64(),
		Email:         req.Email,
		Mobile:        req.Mobile,
		EnableOpenapi: req.EnableOpenapi,
	})

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserUpdateFailed)
		return
	}

	sourceRoleNames := s.UserService.GetUserRoleNames(ctx, req.Id)

	_, err = client.GetInstance().Role.UpdateObjectRoles(ctx, &rbac.ObjectRoles{
		Id: &rbac.ObjectID{
			Id:   req.Id,
			Type: rbac.ObjectType_USER,
		},
		Roles: req.Roles,
	})

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserUpdateFailed)
		return
	}

	user, _ := s.UserService.Get(ctx, snowflake.MustParseString(req.Id).Int64())
	var userName string
	if user != nil {
		userName = user.Name
	}

	targetRoleNames := s.UserService.GetUserRoleNames(ctx, req.Id)

	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_USER_MANAGER, fmt.Sprintf("用户%v修改用户[%v]【%v -》%v】", ginutil.GetUserName(ctx), userName, sourceRoleNames, targetRoleNames))

	ginutil.Success(ctx, nil)
}

// OptionList
//
//	@Summary		用户下拉框
//	@Description	用户下拉框
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			filterName	query	string					false	"用户名搜索"	default("")
//	@Success		200			{array}	dto.UserOptionResponse	"正常返回表示成功，否则返回错误码和错误信息"
//	@Router			/user/optionList [get]
func (s *RouteService) OptionList(ctx *gin.Context) {
	filterName := ctx.Query("filterName")
	filterPerm := ctx.Query("filterPerm")
	var filterPermI int
	if strutil.IsNotEmpty(filterPerm) {
		filterPermI, _ = strconv.Atoi(filterPerm)
	}

	optionList, err := s.UserService.OptionList(ctx, filterName, int64(filterPermI))

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserOptionList)
		return
	}

	ginutil.Success(ctx, optionList)
}

// ResetPassword
//
//	@Summary		重置密码
//	@Description	重置密码
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.ResetPassword	true	"入参"
//	@Success		200		{string}  password
//	@Router			/user/resetPassword [post]
func (s *RouteService) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPassword

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	password, err := s.UserService.ResetPassword(ctx, req.UserID)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserOptionList)
		return
	}

	ginutil.Success(ctx, password)
}

// GenOpenapiCertificate
//
//	@Summary		生成OpenAPI凭证
//	@Description	生成OpenAPI凭证
//	@Tags			用户-user
//	@Accept			json
//	@Produce		json
//	@Param			param	body	dto.GenOpenapiCertificateRequest	true	"入参"
//	@Success		200		{string}  certificate
//	@Router			/user/genOpenapiCertificate [put]
func (s *RouteService) GenOpenapiCertificate(ctx *gin.Context) {
	var req dto.GenOpenapiCertificateRequest
	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrFileFailToBindRequest, "failed to bind request, err: %v", err)
		return
	}

	certificate, err := s.UserService.GenOpenapiCertificate(ctx, req.UserID, true)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserOpenapiCertGenFailed)
		return
	}

	ginutil.Success(ctx, certificate)
}

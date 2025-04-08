package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/user/dto"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
)

func (s *RouteService) CreateOrg(ctx *gin.Context) {

	var req dto.CreateOrgRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrInvalidParam, "failed to bind request, err: %v", err)
		return
	}

	err := s.OrgService.CreateOrg(ctx, req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgCreatedFailed)
		return
	}

	ginutil.Success(ctx, nil)
}

func (s *RouteService) DeleteOrg(ctx *gin.Context) {
	ID := ctx.Query("id")

	if strutil.IsEmpty(ID) {
		errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrInvalidParam)
		return
	}
	orgID := snowflake.MustParseString(ID)
	//user, err := s.OrgService,get
	//
	//if err != nil || user == nil {
	//	errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrUserNotExist)
	//	return
	//}

	err := s.OrgService.Delete(ctx, orgID)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgDeleted)
		return
	}

	ginutil.Success(ctx, nil)
}

func (s *RouteService) UpdateOrg(ctx *gin.Context) {

	var req dto.UpdateOrgRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrInvalidParam, "failed to bind request, err: %v", err)
		return
	}

	err := s.OrgService.UpdateOrg(ctx, req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgUpdateFailed)
		return
	}

	ginutil.Success(ctx, nil)
}

func (s *RouteService) AddOrgMember(ctx *gin.Context) {
	var req dto.AddOrgMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrInvalidParam, "failed to bind request, err: %v", err)
		return
	}

	err := s.OrgService.AddOrgMember(ctx, req.OrgID, req.UserList)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgMemberAddFailed)
		return
	}

	ginutil.Success(ctx, nil)
}

func (s *RouteService) DeleteOrgMember(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	req := &dto.DeleteOrgMemberRequest{}
	if err := ctx.BindQuery(req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	err := s.OrgService.DeleteOrgMemberByID(ctx, req.IDs)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgMemberDeletedFailed)
		return
	}

	ginutil.Success(ctx, nil)
}

func (s *RouteService) UpdateOrgMember(ctx *gin.Context) {

	var req dto.UpdateOrgMemberRequest

	if err := ctx.BindJSON(&req); err != nil {
		http.Errf(ctx, errcode.ErrInvalidParam, "failed to bind request, err: %v", err)
		return
	}

	err := s.OrgService.UpdateOrgMember(ctx, req)
	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgMemberUpdateFailed)
		return
	}

	ginutil.Success(ctx, nil)

}

func (s *RouteService) ListOrgMember(ctx *gin.Context) {
	ID := ctx.Query("org_id")

	if strutil.IsEmpty(ID) {
		errcode.ResolveErrCodeMessage(ctx, nil, errcode.ErrInvalidParam)
		return
	}
	orgID := snowflake.MustParseString(ID)

	member, err := s.OrgService.ListOrgMember(ctx, orgID)

	if err != nil {
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrOrgMemberListFailed)
		return
	}

	ginutil.Success(ctx, member)
}

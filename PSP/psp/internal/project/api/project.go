package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/util"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// ProjectList
//
//	@Summary		项目列表
//	@Description	项目列表接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ProjectListRequest	true	"请求参数"
//	@Response		200		{object}	dto.ProjectListResponse
//	@Router			/project/list [post]
func (r *apiRoute) ProjectList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectListRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	resp, err := r.projectService.ProjectList(ctx, &req, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("List project err: %v", err)
		util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectList)
		return
	}

	ginutil.Success(ctx, resp)
}

// CurrentProjectList
//
//	@Summary		当前用户项目列表
//	@Description	当前用户项目列表接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.CurrentProjectListRequest	true	"请求参数"
//	@Response		200		{object}	dto.CurrentProjectListResponse
//	@Router			/project/list/current [get]
func (r *apiRoute) CurrentProjectList(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.CurrentProjectListRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	resp, err := r.projectService.CurrentProjectList(ctx, &req, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("List current project err: %v", err)
		util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectCurrentList)
		return
	}

	ginutil.Success(ctx, resp)
}

// CurrentProjectListForParam
//
//	@Summary		当前用户项目列表(只用做条件参数“项目”下拉框填充数据)
//	@Description	当前用户项目列表接口(只用做条件参数“项目”下拉框填充数据, 可按照项目最近创建时间区间去展示)
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.CurrentProjectListForParamRequest	true	"请求参数"
//	@Response		200		{object}	dto.CurrentProjectListForParamResponse
//	@Router			/project/listForParam [get]
func (r *apiRoute) CurrentProjectListForParam(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.CurrentProjectListForParamRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	resp, err := r.projectService.CurrentProjectListForParam(ctx, &req, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("List current project for param err: %v", err)
		util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectCurrentListForParam)
		return
	}

	ginutil.Success(ctx, resp)
}

// ProjectSave
//
//	@Summary		保存项目
//	@Description	保存接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ProjectAddRequest	true	"请求参数"
//	@Response		200		{object}	dto.ProjectAddResponse
//	@Router			/project/save [post]
func (r *apiRoute) ProjectSave(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectAddRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if ok, err := util.ValidOwnerID(req.ProjectOwner); !ok {
		logger.Errorf("reqeust owner id err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectOwnerID)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	tracelog.Info(ctx, fmt.Sprintf("save project: user:[%v], params: [%+v]", loginUserID, req))
	resp, err := r.projectService.ProjectSave(ctx, &req, snowflake.ID(loginUserID), ginutil.GetUserName(ctx))
	if err != nil {
		logger.Errorf("save project err: %v", err)
		if status.Code(err) == errcode.ErrProjectNameIsDefault {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectNameIsDefault)
		} else if status.Code(err) == errcode.ErrProjectSameName {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectSameName)
		} else {
			util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectAdd)
		}

		return
	}

	identities := make([]*user.UserIdentity, 0)
	for _, id := range req.Members {
		identities = append(identities, &user.UserIdentity{
			Id: id,
		})
	}
	rsp, _ := client.GetInstance().User.BatchGetUser(ctx, &user.UserIdentities{UserIdentities: identities})
	if rsp != nil && len(rsp.UserObj) > 0 {
		memberName := make([]string, 0)
		for _, userInfo := range rsp.UserObj {
			memberName = append(memberName, userInfo.Name)
		}
		oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_PROJECT_MANAGER, fmt.Sprintf("用户%v新增项目[%v]，项目成员为%v", ginutil.GetUserName(ctx), req.ProjectName, memberName))
	}

	ginutil.Success(ctx, resp)
}

// ProjectDetail
//
//	@Summary		项目详情
//	@Description	详情接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ProjectDetailRequest	true	"请求参数"
//	@Response		200		{object}	dto.ProjectDetailResponse
//	@Router			/project/detail [get]
func (r *apiRoute) ProjectDetail(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectDetailRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	loginUserID := ginutil.GetUserID(ctx)

	resp, err := r.projectService.ProjectDetail(ctx, &req, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("get project detail err: %v", err)
		util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectDetail)
		return
	}

	ginutil.Success(ctx, resp)
}

// ProjectDelete
//
//	@Summary		删除项目
//	@Description	删除项目接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ProjectDeleteRequest	true	"请求参数"
//	@Router			/project/delete [post]
func (r *apiRoute) ProjectDelete(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectDeleteRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}
	loginUserID := ginutil.GetUserID(ctx)

	tracelog.Info(ctx, fmt.Sprintf("delete project: user:[%v], projectID:[%v]", loginUserID, req.ProjectID))
	err := r.projectService.ProjectDelete(ctx, req.ProjectID, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("delete project err: %v", err)
		if status.Code(err) == errcode.ErrProjectDelBeforeExistMembers {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectDelBeforeExistMembers)
		} else if status.Code(err) == errcode.ErrProjectDeleteState {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectDeleteState)
		} else {
			util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectDelete)
		}
		return
	}

	ginutil.Success(ctx, nil)
}

// ProjectTerminate
//
//	@Summary		终止项目
//	@Description	终止项目接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ProjectTerminatedRequest	true	"请求参数"
//	@Router			/project/terminate [post]
func (r *apiRoute) ProjectTerminate(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectTerminatedRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	tracelog.Info(ctx, fmt.Sprintf("terminate project: user:[%v], projectID:[%v]", loginUserID, req.ProjectID))
	err := r.projectService.ProjectTerminate(ctx, req.ProjectID, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("terminate project err: %v", err)
		if status.Code(err) == errcode.ErrProjectTerminateState {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectTerminateState)
		} else {
			util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectTerminated)
		}

		return
	}

	projectInfo, _ := r.projectService.ProjectDetail(ctx, &dto.ProjectDetailRequest{ProjectID: req.ProjectID}, snowflake.ID(loginUserID))
	var projectName string
	if projectInfo != nil {
		projectName = projectInfo.ProjectName
	}
	oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_PROJECT_MANAGER, fmt.Sprintf("用户%v终止项目[%v]", ginutil.GetUserName(ctx), projectName))

	ginutil.Success(ctx, nil)
}

// ProjectEdit
//
//	@Summary		编辑项目
//	@Description	编辑项目接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ProjectEditRequest	true	"请求参数"
//	@Router			/project/edit [post]
func (r *apiRoute) ProjectEdit(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectEditRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if ok, err := util.ValidOwnerID(req.ProjectOwner); !ok {
		logger.Errorf("reqeust project owner id err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectOwnerID)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	tracelog.Info(ctx, fmt.Sprintf("edit project: user:[%v], params:[%+v]", loginUserID, req))
	err := r.projectService.ProjectEdit(ctx, &req, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("edit project err: %v", err)
		if status.Code(err) == errcode.ErrProjectEditState {
			errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectEditState)
		} else {
			util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectEdit)
		}
		return
	}

	projectInfo, _ := r.projectService.ProjectDetail(ctx, &dto.ProjectDetailRequest{ProjectID: req.ProjectID}, snowflake.ID(loginUserID))
	var projectName string
	if projectInfo != nil {
		projectName = projectInfo.ProjectName
	}

	identities := make([]*user.UserIdentity, 0)
	for _, id := range req.Members {
		identities = append(identities, &user.UserIdentity{
			Id: id,
		})
	}
	rsp, _ := client.GetInstance().User.BatchGetUser(ctx, &user.UserIdentities{UserIdentities: identities})
	if rsp != nil && len(rsp.UserObj) > 0 {
		memberName := make([]string, 0)
		for _, userInfo := range rsp.UserObj {
			memberName = append(memberName, userInfo.Name)
		}
		oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_PROJECT_MANAGER, fmt.Sprintf("用户%v编辑项目[%v]，项目成员为%v", ginutil.GetUserName(ctx), projectName, memberName))
	}

	ginutil.Success(ctx, nil)
}

// ProjectModifyOwner
//
//	@Summary		修改项目管理员
//	@Description	修改项目管理员接口
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query	dto.ProjectModifyOwnerRequest	true	"请求参数"
//	@Router			/project/modifyOwner [post]
func (r *apiRoute) ProjectModifyOwner(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectModifyOwnerRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	if ok, err := util.ValidOwnerID(req.TargetProjectOwnerID); !ok {
		logger.Errorf("reqeust project owner id err: %v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectOwnerID)
		return
	}

	loginUserID := ginutil.GetUserID(ctx)

	tracelog.Info(ctx, fmt.Sprintf("edit project owner: user:[%v], params: [%+v]", loginUserID, req))
	err := r.projectService.ProjectModifyOwner(ctx, &req, snowflake.ID(loginUserID))
	if err != nil {
		logger.Errorf("modify project owner err: %v", err)
		util.CheckErrorIfNoPermission(ctx, err, errcode.ErrProjectModifyOwner)
		return
	}

	ginutil.Success(ctx, nil)
}

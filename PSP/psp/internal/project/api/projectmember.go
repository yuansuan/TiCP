package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/user"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/project/service/client"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/tracelog"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

// BatchUpdateProjectMember
//
//	@Summary		保存项目成员
//	@Description	保存项目成员
//	@Tags			项目
//	@Accept			json
//	@Produce		json
//	@Param			param	query		dto.ProjectMemberRequest	true	"请求参数"
//	@Response		200		{object}	dto.ProjectMemberResponse
//	@Router			/projectMember/save [post]
func (r *apiRoute) BatchUpdateProjectMember(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	req := dto.ProjectMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("reqeust params bind err: %v", err)
		ginutil.Error(ctx, errcode.ErrInvalidParam, errcode.MsgInvalidParam)
		return
	}

	userID := ginutil.GetUserID(ctx)

	tracelog.Info(ctx, fmt.Sprintf("save project member: userID:[%v], params:[%v]", userID, req))
	resp, err := r.projectMemberService.ProjectMemberSave(ctx, &req, snowflake.ID(userID))
	if err != nil {
		logger.Errorf("add project member error! err:%v", err)
		errcode.ResolveErrCodeMessage(ctx, err, errcode.ErrProjectMemberAdd)
		return
	}

	projectInfo, _ := r.projectService.ProjectDetail(ctx, &dto.ProjectDetailRequest{ProjectID: req.ProjectId}, snowflake.ID(userID))
	var projectName string
	if projectInfo != nil {
		projectName = projectInfo.ProjectName
	}
	identities := make([]*user.UserIdentity, 0)
	for _, id := range req.UserIds {
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

		oplog.GetInstance().SaveAuditLogInfo(ctx, approve.OperateTypeEnum_PROJECT_MANAGER, fmt.Sprintf("用户%v更改项目[%v]的成员为%v", ginutil.GetUserName(ctx), projectName, memberName))
	}

	ginutil.Success(ctx, resp)
}

package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/impl/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/oplog"
	pbapprove "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	pbnotice "github.com/yuansuan/ticp/PSP/psp/internal/common/proto/notice"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/with"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type ApproveServiceImpl struct {
	sid        *snowflake.Node
	approveDao dao.ApproveDao
}

func (srv ApproveServiceImpl) CheckUnhandledApprove(ctx context.Context, userId int64) (bool, error) {
	return srv.approveDao.CheckUnhandledApprove(ctx, userId)
}

func (srv ApproveServiceImpl) Refuse(ctx *gin.Context, req *dto.HandleApproveRequest) error {
	recordID := snowflake.MustParseString(req.RecordID)
	record, err := srv.approveDao.GetRecord(ctx, recordID)
	if err != nil {
		return err
	}
	if record.Status != int8(dto.ApproveStatusWaiting) {
		return status.Error(errcode.ErrApproveStatusEnd, "")
	}

	err = with.DefaultTransaction(ctx, func(tCtx context.Context) error {
		err = srv.approveDao.UpdateApproveUser(tCtx, &model.ApproveUser{
			Id:          snowflake.MustParseString(req.ID),
			Result:      int8(dto.ApproveResultRefuse),
			Suggest:     req.Suggest,
			ApproveTime: time.Now(),
		})
		if err != nil {
			return err
		}

		err = srv.approveDao.UpdateApproveRecord(tCtx, &model.ApproveRecord{
			Id:          recordID,
			Status:      int8(dto.ApproveStatusRefuse),
			ApproveTime: time.Now(),
		})
		if err != nil {
			return err
		}
		oplog.GetInstance().SaveAuditLogInfo(ctx, pbapprove.OperateTypeEnum_SECURITY_APPROVAL, fmt.Sprintf("用户%v拒绝%v发起的审批:[%v]", ginutil.GetUserName(ctx), record.ApplyUserName, record.Content))

		return nil
	})

	srv.SendCompleteNotice(ctx, record.ApplyUserId, false, record.Content)

	return err
}

func (srv ApproveServiceImpl) Pass(ctx *gin.Context, req *dto.HandleApproveRequest) error {

	handleSrv := approve.GetImpl(req.ApproveType)
	recordID := snowflake.MustParseString(req.RecordID)

	record, err := srv.approveDao.GetRecord(ctx, recordID)
	if err != nil {
		return err
	}
	if record.Status != int8(dto.ApproveStatusWaiting) {
		return status.Error(errcode.ErrApproveStatusEnd, "")
	}

	var notReadyErr error
	var ready bool
	err = with.DefaultTransaction(ctx, func(tCtx context.Context) error {
		// 改approve_user表
		err := srv.approveDao.UpdateApproveUser(tCtx, &model.ApproveUser{
			Id:          snowflake.MustParseString(req.ID),
			Result:      int8(dto.ApproveResultPass),
			Suggest:     req.Suggest,
			ApproveTime: time.Now(),
		})
		if err != nil {
			return err
		}

		// 判断是否要改approve_record表；
		finish, err := srv.approveDao.AllApproved(tCtx, recordID)
		if err != nil {
			return err
		}
		if !finish {
			return nil
		}

		// handleSrv.CheckNecessary 返回true false，如果是false直接把approve_record状态改为失败然后return
		info, err := handleSrv.ParseObject(record.ApproveInfo)
		if err != nil {
			return err
		}

		// 检查审批数据落地的必要条件
		ready, notReadyErr, err = handleSrv.CheckNecessary(tCtx, req.ApproveType, info)
		if err != nil {
			return err
		}
		// 改approve_record表
		if !ready {
			err = srv.approveDao.UpdateApproveRecord(tCtx, &model.ApproveRecord{
				Id:          recordID,
				Status:      int8(dto.ApproveStatusFailed),
				ApproveTime: time.Now(),
			})
			if err != nil {
				return err
			}
			return nil
		}
		err = srv.approveDao.UpdateApproveRecord(tCtx, &model.ApproveRecord{
			Id:          recordID,
			Status:      int8(dto.ApproveStatusPass),
			ApproveTime: time.Now(),
		})

		err = handleSrv.AfterPass(ctx, info)
		if err != nil {
			return err
		}
		oplog.GetInstance().SaveAuditLogInfo(ctx, pbapprove.OperateTypeEnum_SECURITY_APPROVAL, fmt.Sprintf("用户%v同意%v发起的审批:[%v]", ginutil.GetUserName(ctx), record.ApplyUserName, record.Content))

		return nil
	})
	if notReadyErr != nil {
		return notReadyErr
	}

	// 给发起审批用户发个消息
	if err != nil {
		return err
	}
	srv.SendCompleteNotice(ctx, record.ApplyUserId, true, record.Content)

	return nil
}

func (srv ApproveServiceImpl) SendCompleteNotice(ctx context.Context, userID snowflake.ID, isSuccess bool, content string) {
	logger := logging.Default()
	pass := "同意"
	if !isSuccess {
		pass = "拒绝"
	}

	msg := &pbnotice.WebsocketMessage{
		UserId:  userID.String(),
		Type:    common.ApproveEventType,
		Content: fmt.Sprintf("您提交的审批:[%s]已完成，结果为[%s]", content, pass),
	}

	if _, err := client.GetInstance().Notice.SendWebsocketMessage(ctx, msg); err != nil {
		logger.Errorf("approve complete send ws message err: %v", err)
	}
}

func (srv ApproveServiceImpl) ApplyApprove(ctx *gin.Context, req *dto.ApplyApproveRequest) error {

	handleSrv := approve.GetImpl(req.ApproveType)
	// 前置检查,参数校验等
	err := handleSrv.PreCheck(ctx, req)
	if err != nil {
		return err
	}

	// 生成唯一审批签名
	sign, err := handleSrv.GenSign(ctx, req)
	if err != nil {
		return err
	}

	// 检查签名是否重复
	exist, err := srv.approveDao.CheckSign(ctx, sign)
	if err != nil {
		return err
	}
	if exist {
		return status.Error(errcode.ErrApproveHasConflict, "")
	}

	// 构建审批详情内容
	info, err := handleSrv.BuildApproveInfo(ctx, req)
	if err != nil {
		return err
	}
	// 生成审批文案
	content, err := handleSrv.BuildContent(ctx, info)
	if err != nil {
		return err
	}

	approveInfo, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// 生成审批记录与关联审批人
	userName := ginutil.GetUserName(ctx)
	userID := ginutil.GetUserID(ctx)

	err = with.DefaultTransaction(ctx, func(tCtx context.Context) error {
		approveRecord := &model.ApproveRecord{
			Id:            srv.sid.Generate(),
			Type:          int8(req.ApproveType),
			ApproveInfo:   string(approveInfo),
			Status:        int8(dto.ApproveStatusWaiting),
			ApplyUserId:   snowflake.ID(userID),
			ApplyUserName: userName,
			Sign:          sign,
			Content:       content,
			CreateTime:    time.Now(),
		}
		err = srv.approveDao.AddApproveRecord(tCtx, approveRecord)
		if err != nil {
			return err
		}

		approveUser := &model.ApproveUser{
			Id:              srv.sid.Generate(),
			ApproveRecordId: approveRecord.Id,
			ApproveUserId:   snowflake.MustParseString(req.ApproveUserID),
			ApproveUserName: req.ApproveUserName,
			Result:          int8(dto.ApproveResultDefault),
			CreateTime:      time.Now(),
		}
		err = srv.approveDao.AddApproveUser(tCtx, approveUser)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	oplog.GetInstance().SaveAuditLogInfo(ctx, pbapprove.OperateTypeEnum_SECURITY_APPROVAL, fmt.Sprintf("用户%v发起审批:[%v]", ginutil.GetUserName(ctx), content))
	return nil
}

func (a ApproveServiceImpl) CancelApprove(ctx *gin.Context, reqId int64) error {
	err := a.approveDao.CancelApprove(ctx, reqId)
	if err != nil {
		return err
	}

	record, _ := a.approveDao.GetRecord(ctx, snowflake.ID(reqId))
	if record != nil {
		oplog.GetInstance().SaveAuditLogInfo(ctx, pbapprove.OperateTypeEnum_SECURITY_APPROVAL, fmt.Sprintf("用户%v撤销审批:[%v]", ginutil.GetUserName(ctx), record.Content))
	}
	return nil
}

func (a ApproveServiceImpl) GetApproveList(ctx context.Context, UserId int64, request *dto.GetApproveListRequest) (*dto.ApproveLogListAllResponse, error) {
	if UserId <= 0 {
		return nil, fmt.Errorf("userid not exist")
	}
	resp, err := client.GetInstance().Rbac.GetRoleByObjectID(ctx, &rbac.ObjectID{
		Id:   snowflake.ID(UserId).String(),
		Type: rbac.ObjectType_USER,
	})

	if err != nil {
		return nil, err
	}
	isAdmin := false
	for _, role := range resp.Roles {
		if isAdmin = role.Type == rbac.RoleType_ROLE_SUPER_ADMIN; isAdmin {
			break
		}
	}
	if isAdmin {
		UserId = 0
	}

	list, total, err := a.approveDao.ApproveList(ctx, request.Page, &dto.ApproveListCondition{
		ApplyId:    UserId,
		RecordType: request.Type,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Status:     request.Status,
	})
	infos := make([]*dto.ApproveInfo, 0)

	if total > 0 {
		for _, log := range list {
			infos = append(infos, &dto.ApproveInfo{
				Id:                log.ApproveRecord.Id.String(),
				RecordID:          log.ApproveRecord.Id.String(),
				ApplicationName:   log.ApplyUserName,
				Type:              log.Type,
				Content:           log.Content,
				ApproveUserName:   log.ApproveUserName,
				Status:            log.ApproveRecord.Status,
				Suggest:           log.ApproveUser.Suggest,
				ApproveCreateTime: timeutil.FormatTime(log.ApproveRecord.CreateTime, common.DatetimeFormat),
				ApproveTime:       timeutil.FormatTime(log.ApproveRecord.ApproveTime, common.DatetimeFormat),
			})
		}
	}
	if err != nil {
		return nil, err
	}
	return &dto.ApproveLogListAllResponse{
		Page: &xtype.PageResp{
			Total: total,
			Index: request.Page.Index,
			Size:  request.Page.Size,
		},
		LogInfo: infos,
	}, nil
}

func (a ApproveServiceImpl) GetApprovePendingList(ctx context.Context, UserId int64, request *dto.GetApprovePendingRequest) (*dto.ApproveLogListAllResponse, error) {
	if UserId <= 0 {
		return nil, fmt.Errorf("userid not exist")
	}
	list, total, err := a.approveDao.ApplicationList(ctx, request.Page, &dto.ApplicationListCondition{
		UserId:     UserId,
		ApplyName:  request.ApplicationName,
		RecordType: request.Type,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Status:     []int8{int8(dto.ApproveStatusWaiting)},
	})
	if err != nil {
		return nil, err
	}
	infos := make([]*dto.ApproveInfo, 0)

	if total > 0 {
		for _, log := range list {
			infos = append(infos, &dto.ApproveInfo{
				Id:                log.ApproveUser.Id.String(),
				RecordID:          log.ApproveRecord.Id.String(),
				ApplicationName:   log.ApproveRecord.ApplyUserName,
				Type:              log.Type,
				Content:           log.Content,
				ApproveUserName:   log.ApproveUser.ApproveUserName,
				Status:            log.ApproveRecord.Status,
				Suggest:           log.ApproveUser.Suggest,
				ApproveCreateTime: timeutil.FormatTime(log.ApproveUser.CreateTime, common.DatetimeFormat),
				ApproveTime:       timeutil.FormatTime(log.ApproveUser.ApproveTime, common.DatetimeFormat),
			})
		}
	}
	return &dto.ApproveLogListAllResponse{
		Page: &xtype.PageResp{
			Total: total,
			Index: request.Page.Index,
			Size:  request.Page.Size,
		},
		LogInfo: infos,
	}, nil
}

func (a ApproveServiceImpl) GetApprovedList(ctx context.Context, UserId int64, request *dto.GetApproveCompleteRequest) (*dto.ApproveLogListAllResponse, error) {
	if UserId <= 0 {
		return nil, fmt.Errorf("userid not exist")
	}

	var defaultStatus []int8
	defaultStatus = []int8{int8(dto.ApproveStatusPass), int8(dto.ApproveStatusFailed), int8(dto.ApproveStatusRefuse)}

	if request.Status != 0 {
		exist := false

		// request.Status in defaultStatus
		for _, status := range defaultStatus {
			if status == request.Status {
				exist = true
			}
		}
		if exist {
			defaultStatus = []int8{int8(request.Status)}
		} else {
			return nil, fmt.Errorf("can't get waiting list")
		}
	}

	list, total, err := a.approveDao.ApplicationList(ctx, request.Page, &dto.ApplicationListCondition{
		UserId:     UserId,
		ApplyName:  request.ApplicationName,
		RecordType: request.Type,
		StartTime:  request.StartTime,
		EndTime:    request.EndTime,
		Status:     defaultStatus,
	})
	if err != nil {
		return nil, err
	}
	infos := make([]*dto.ApproveInfo, 0)

	if total > 0 {
		for _, log := range list {
			infos = append(infos, &dto.ApproveInfo{
				Id:                log.ApproveUser.Id.String(),
				RecordID:          log.ApproveRecord.Id.String(),
				ApplicationName:   log.ApproveRecord.ApplyUserName,
				Type:              log.Type,
				Content:           log.Content,
				ApproveUserName:   log.ApproveUser.ApproveUserName,
				ApproveCreateTime: timeutil.FormatTime(log.ApproveUser.CreateTime, common.DatetimeFormat),
				ApproveTime:       timeutil.FormatTime(log.ApproveUser.ApproveTime, common.DatetimeFormat),
				Status:            log.ApproveRecord.Status,
				Suggest:           log.ApproveUser.Suggest,
			})
		}
	}

	return &dto.ApproveLogListAllResponse{
		Page: &xtype.PageResp{
			Total: total,
			Index: request.Page.Index,
			Size:  request.Page.Size,
		},
		LogInfo: infos,
	}, nil
}

func NewApproveService() (service.ApproveService, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("new snowflake node err: %v", err)
		return nil, err
	}

	return &ApproveServiceImpl{
		sid:        node,
		approveDao: dao.NewApproveDao(),
	}, nil
}

package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/dto"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/client"
	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/pkg/snowflake"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/csvutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type AuditLogServiceImpl struct {
	sid         *snowflake.Node
	auditLogDao dao.AuditLogDao
}

func (srv AuditLogServiceImpl) List(ctx *gin.Context, req *dto.AuditLogListRequest) (*dto.AuditLogListResponse, error) {
	userId := ginutil.GetUserID(ctx)
	resp, err := client.GetInstance().Rbac.GetRoleByObjectID(ctx, &rbac.ObjectID{
		Id:   snowflake.ID(userId).String(),
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

	filterUserId := userId
	if isAdmin {
		filterUserId = 0
	}

	list, total, err := srv.auditLogDao.List(ctx, req.Page, filterUserId, req.UserName, req.IpAddress, dto.OperateTypeString(req.OperateType), req.StartTime, req.EndTime)

	infos := make([]*dto.AuditLogInfo, 0)

	if err != nil {
		return nil, status.Error(errcode.ErrApproveAuditLogAdd, "")
	}

	if total > 0 {
		for _, log := range list {
			infos = append(infos, &dto.AuditLogInfo{
				Id:             log.Id.String(),
				UserName:       log.UserName,
				OperateType:    log.OperateType,
				OperateContent: log.OperateContent,
				IpAddress:      log.IpAddress,
				OperateTime:    log.OperateTime,
			})
		}
	}

	return &dto.AuditLogListResponse{
		Page: &xtype.PageResp{
			Index: req.Page.Index,
			Size:  req.Page.Size,
			Total: total,
		},
		AuditLogInfo: infos,
	}, nil
}

func (srv AuditLogServiceImpl) ListAll(ctx *gin.Context, req *dto.AuditLogListAllRequest) (*dto.AuditLogListAllResponse, error) {
	userId := ginutil.GetUserID(ctx)
	list, total, err := srv.auditLogDao.ListAll(ctx, req.Page, userId, req.UserName, req.IpAddress, dto.OperateTypeString(req.OperateType), req.StartTime, req.EndTime, req.OperateUserType)
	infos := make([]*dto.AuditLogExportInfo, 0)

	if err != nil {
		return nil, status.Error(errcode.ErrApproveAuditLogAdd, "")
	}

	if total > 0 {
		for _, log := range list {
			infos = append(infos, &dto.AuditLogExportInfo{
				Id:             log.Id.String(),
				UserName:       log.UserName,
				OperateType:    log.OperateType,
				OperateContent: log.OperateContent,
				IpAddress:      log.IpAddress,
				OperateTime:    log.OperateTime,
			})
		}
	}

	var auditLogListAllResponse dto.AuditLogListAllResponse
	auditLogListAllResponse.AuditLogInfo = infos
	auditLogListAllResponse.Page = &xtype.PageResp{
		Index: req.Page.Index,
		Size:  req.Page.Size,
		Total: total,
	}

	return &auditLogListAllResponse, nil
}

func buildAuditLogHeader() []csvutil.CsvHeaderEntity {
	return []csvutil.CsvHeaderEntity{
		{
			Name:   "用户名称",
			Column: "UserName",
		},
		{
			Name:   "IP地址",
			Column: "IpAddress",
		},
		{
			Name:   "操作时间",
			Column: "OperateTime",
			Converter: func(i interface{}) string {
				return csvutil.CSVFormatTime(i.(time.Time), common.DatetimeFormat, common.Bar)
			},
		},
		{
			Name:   "操作类型",
			Column: "OperateType",
		},
		{
			Name:   "操作内容",
			Column: "OperateContent",
		},
	}
}

func (srv AuditLogServiceImpl) Export(ctx *gin.Context, req *dto.AuditLogListRequest) error {

	header := buildAuditLogHeader()

	csvFileName := fmt.Sprintf("%s", "用户操作记录")

	return csvutil.LargeDataExportCsv(ctx, header, func(page *xtype.Page) ([]interface{}, error) {
		req.Page = page
		rsp, err := srv.List(ctx, req)
		list := rsp.AuditLogInfo
		if err != nil {
			return nil, err
		}
		data := make([]interface{}, len(list))
		for i := range list {
			data[i] = *list[i]
		}

		return data, nil
	}, csvFileName)

}

func (srv AuditLogServiceImpl) ExportAll(ctx *gin.Context, req *dto.AuditLogListAllRequest) error {

	header := buildAuditLogHeader()

	csvFileName := fmt.Sprintf("%s", "用户操作记录")

	return csvutil.LargeDataExportCsv(ctx, header, func(page *xtype.Page) ([]interface{}, error) {
		req.Page = page
		rsp, err := srv.ListAll(ctx, req)
		list := rsp.AuditLogInfo
		if err != nil {
			return nil, err
		}
		data := make([]interface{}, len(list))
		for i := range list {
			data[i] = *list[i]
		}

		return data, nil
	}, csvFileName)

}

func (srv AuditLogServiceImpl) SaveLog(ctx context.Context, req *dto.SaveLogRequest) error {

	if approve.OperateTypeEnum_FILE_MANAGER == req.OperateType {
		// 如果是操作.tmp_upload临时目录，则不需要记录
		// todo 之后加入到配置文件中
		if strings.Contains(req.OperateContent, "[.tmp_upload/") {
			return nil
		}
	}

	// 系统管理员>安全管理员>普通用户（系统管理员和安全管理员根据权限标识判断，其余为普通用户）
	UserId := req.UserId
	var OperateUserType dto.OperateUserType

	permissions, err := client.GetInstance().Perm.ListObjectPermissions(ctx, &rbac.ObjectID{Id: UserId.String(), Type: rbac.ObjectType_USER})
	if err != nil {
		return err
	}
	var UserTypeSecurity, UserTypeAdmin bool
	for _, permission := range permissions.Perms {

		if permission.ResourceName == common.ResourceSecurityManagerName {
			UserTypeSecurity = true
		}

		if permission.ResourceName == common.ResourceSysManagerName {
			UserTypeAdmin = true
		}
	}

	if UserTypeAdmin {
		OperateUserType = dto.OperateUserTypeAdmin
	} else if UserTypeSecurity {
		OperateUserType = dto.OperateUserTypeSecurity
	}

	if OperateUserType == 0 {
		OperateUserType = dto.OperateUserTypeUser
	}

	err = srv.auditLogDao.Add(ctx, &model.AuditLog{
		Id:              srv.sid.Generate(),
		UserId:          req.UserId,
		UserName:        req.UserName,
		IpAddress:       req.IpAddress,
		OperateType:     dto.OperateTypeString(req.OperateType),
		OperateContent:  req.OperateContent,
		OperateTime:     time.Now(),
		OperateUserType: OperateUserType,
	})

	if err != nil {
		return status.Error(errcode.ErrAppAddAppFailed, "")
	}

	return nil
}

func NewAuditLogService() (service.AuditLogService, error) {
	node, err := snowflake.GetInstance()
	if err != nil {
		logging.Default().Errorf("new snowflake node err: %v", err)
		return nil, err
	}

	return &AuditLogServiceImpl{
		sid:         node,
		auditLogDao: dao.NewLogDao(),
	}, nil
}

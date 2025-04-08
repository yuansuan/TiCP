package api

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/impl"
)

type RouteService struct {
	AuditLogService service.AuditLogService
	ApproveService  service.ApproveService
}

func NewAuditLogService() (*RouteService, error) {
	auditLogService, err := impl.NewAuditLogService()
	if err != nil {
		logging.Default().Errorf("init sys config server service err: %v", err)
		return nil, err
	}
	approveService, err := impl.NewApproveService()
	if err != nil {
		logging.Default().Errorf("init sys config server service err: %v", err)
		return nil, err
	}

	return &RouteService{
		AuditLogService: auditLogService,
		ApproveService:  approveService,
	}, nil
}

func InitAPI(drv *http.Driver) {
	logger := logging.Default()

	s, err := NewAuditLogService()
	if err != nil {
		logger.Errorf("init api service err: %v", err)
		panic(err)
	}

	auditLogGroup := drv.Group("/api/v1/auditlog")
	{
		auditLogGroup.POST("/list", s.List)
		auditLogGroup.POST("/listAll", s.ListAll)
		auditLogGroup.POST("/export", s.Export)
		auditLogGroup.POST("/exportAll", s.ExportAll)
	}

	approveGroup := drv.Group("/api/v1/approve")
	{
		approveGroup.POST("/apply", s.ApplyApprove)
		approveGroup.POST("/cancel", s.CancelApprove)

		approveGroup.GET("/threePersonManagement", s.ThreePersonManagement)

		listGroup := approveGroup.Group("/list")
		{
			listGroup.POST("/application", s.Application)
			listGroup.POST("/pending", s.Pending)
			listGroup.POST("/complete", s.Complete)
		}
		approveGroup.POST("/pass", s.Pass)
		approveGroup.POST("/refuse", s.Refuse)
	}
}

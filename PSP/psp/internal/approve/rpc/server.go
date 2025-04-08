package rpc

import (
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service"
	"github.com/yuansuan/ticp/PSP/psp/internal/approve/service/impl"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/approve"
)

type GRPCService struct {
	AuditLogService service.AuditLogService
	ApproveService  service.ApproveService
}

func NewGRPCService() (*GRPCService, error) {
	auditLogService, err := impl.NewAuditLogService()
	if err != nil {
		logging.Default().Errorf("init auditlog server service err: %v", err)
		return nil, err
	}
	approveService, err := impl.NewApproveService()
	if err != nil {
		logging.Default().Errorf("init approve server service err: %v", err)
		return nil, err
	}

	return &GRPCService{
		AuditLogService: auditLogService,
		ApproveService:  approveService,
	}, nil
}

func InitGRPC(drv *http.Driver) {
	s, err := boot.GRPC.DefaultServer()
	if err != nil {
		panic(err)
	}

	grpcServer, err := NewGRPCService()
	if err != nil {
		panic(err)
	}

	approve.RegisterAuditLogManagementServer(s.Driver(), grpcServer)
	approve.RegisterApproveManagementServer(s.Driver(), grpcServer)
}

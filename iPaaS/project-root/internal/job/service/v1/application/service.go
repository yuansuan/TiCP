//go:generate mockgen -destination mock_application_srv.go -package application github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application Service,AppSrv,AppQuotaSrv,AppAllowSrv,UserGeter

package application

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/module/rpc"
)

// Service 服务
type Service interface {
	Apps() AppSrv
	AppsQuota() AppQuotaSrv
	AppsAllow() AppAllowSrv
}

type service struct {
	store store.FactoryNew
}

func NewService(store store.FactoryNew) Service {
	return &service{
		store: store,
	}
}

func (s *service) Apps() AppSrv {
	return newAppService(s)
}

func (s *service) AppsQuota() AppQuotaSrv {
	return newAppQuotaService(s, rpc.GetInstance())
}

func (s *service) AppsAllow() AppAllowSrv {
	return newAppAllowService(s, rpc.GetInstance())
}

package application

import (
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/dao/store"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/application"
)

type Controller struct {
	srv srvv1.Service
}

func (app *Controller) GetService() srvv1.Service {
	return app.srv
}

func NewApplicationController(store store.FactoryNew) *Controller {
	return &Controller{
		srv: srvv1.NewService(store),
	}
}

func MockApplicationController(srv srvv1.Service) *Controller {
	return &Controller{
		srv: srv,
	}
}

package service

import (
	jobsrv "github.com/yuansuan/ticp/PSP/psp/internal/job/service"
	jobimpl "github.com/yuansuan/ticp/PSP/psp/internal/job/service/impl"
	rbacsrv "github.com/yuansuan/ticp/PSP/psp/internal/rbac/service"
	rbacimpl "github.com/yuansuan/ticp/PSP/psp/internal/rbac/service/impl"
	usersrv "github.com/yuansuan/ticp/PSP/psp/internal/user/service"
	userimpl "github.com/yuansuan/ticp/PSP/psp/internal/user/service/impl"
)

var RouteSrv *RouteService

type RouteService struct {
	JobService  jobsrv.JobService
	PermService rbacsrv.PermService
	UserService usersrv.UserService
}

func NewRouteService() (*RouteService, error) {
	jobService, err := jobimpl.NewJobService()
	if err != nil {
		return nil, err
	}
	permService, err := rbacimpl.NewPermService()
	if err != nil {
		return nil, err
	}
	userService := userimpl.NewUserService()

	return &RouteService{
		JobService:  jobService,
		PermService: permService,
		UserService: userService,
	}, nil
}

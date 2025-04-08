package v1

//go:generate mockgen -self_package=github.com/yuansuan/ticp/common/project-root-iam/internal/iamserver/service/v1 -destination mock_types.go -package=v1 github.com/yuansuan/ticp/common/project-root-iam/internal/iamserver/service/v1 Svc,SecretSvc,PolicySvc,RoleSvc

import "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"

type Svc interface {
	Secrets() SecretSvc
	Policies() PolicySvc

	Roles() RoleSvc
}

type svc struct {
	store store.Factory
}

func NewSvc(s store.Factory) Svc {
	return &svc{
		store: s,
	}
}

func (s *svc) Secrets() SecretSvc {
	return &secretService{
		store: s.store,
	}
}

func (s *svc) Policies() PolicySvc {
	return &policyService{
		store: s.store,
	}
}

func (s *svc) Roles() RoleSvc {
	return &roleService{
		store: s.store,
	}
}

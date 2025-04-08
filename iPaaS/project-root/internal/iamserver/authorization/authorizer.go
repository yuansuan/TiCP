package authorization

import (
	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
)

type Authorizer struct {
	l *ladon.Ladon
}

func NewAuthorizer(factory store.Factory) *Authorizer {
	return &Authorizer{
		l: &ladon.Ladon{
			AuditLogger: NewAuditLogger(factory),
		},
	}
}

func (a *Authorizer) DoPoliciesAllow(request *ladon.Request, policies []ladon.DefaultPolicy) (bool, error) {
	logging.Default().Debugf("authorize request: %+v", request)
	p := mergePolicy(policies)
	if err := a.l.DoPoliciesAllow(request, p); err != nil {
		logging.Default().Infof("authorize request: %+v, error: %s", request, err.Error())
		return false, err
	}
	return true, nil
}

func mergePolicy(policies []ladon.DefaultPolicy) []ladon.Policy {
	var p []ladon.Policy
	for _, policy := range policies {
		// pointers to distinct policy objects, rather than pointers to the same policy object
		newPolicy := policy
		p = append(p, &newPolicy)
	}
	return p
}

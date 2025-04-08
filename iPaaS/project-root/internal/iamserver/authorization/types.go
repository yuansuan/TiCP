package authorization

import "github.com/ory/ladon"

//go:generate mockgen -self_package=github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/authorization -destination mock_policyCheck.go -package authorization github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/authorization PolicyCheck

type PolicyCheck interface {
	DoPoliciesAllow(r *ladon.Request, policies []ladon.DefaultPolicy) (bool, error)
}

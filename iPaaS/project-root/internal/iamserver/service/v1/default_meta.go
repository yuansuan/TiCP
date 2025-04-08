package v1

import (
	"strings"

	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

const (
	arnServiceIAM string = "iam"
)

const (
	arnResourceTypeRole string = "role"
)

// 默认所有用户都开通这些role
var YSCloudComputeRole *dao.Role
var CSPRole *dao.Role
var CAE365Role *dao.Role

var YSCloudComputeRolePolicy *dao.Policy

func IsValidRoleName(roleName string) bool {
	if len(roleName) == 0 || len(roleName) > 64 {
		return false
	}
	if strings.HasPrefix(roleName, "YS_") {
		return false
	}
	return true
}

func IsYuansuanProductAccount(tag string) bool {
	return strings.HasPrefix(tag, "YS_")
}

// IsValidYrnForRole checks if the given YRN is valid for iam service and role resource type.
func IsValidYrnForRole(yrn *common.YRN) bool {
	if yrn.Service != arnServiceIAM {
		return false
	}
	if yrn.ResourceType != arnResourceTypeRole {
		return false
	}
	return true
}

func IsYusanRole(name string) bool {
	return name == YSCloudComputeRole.RoleName || name == CSPRole.RoleName || name == CAE365Role.RoleName
}

func subjectName(name string) string {
	switch name {
	case YSCloudComputeRole.RoleName:
		return "YS_CloudCompute"
	case CSPRole.RoleName:
		return "YS_CSP"
	case CAE365Role.RoleName:
		return "YS_CAE365"
	// won't happen
	default:
		return "YS_CloudCompute"
	}
}

func getYuansuanRole(name string) *dao.Role {
	if name == YSCloudComputeRole.RoleName {
		return YSCloudComputeRole
	}
	return YSCloudComputeRole
}

// FIXME: change name to managedPolicy
func getYusnuanPolicy(roleName string) []*dao.Policy {
	return PolicySaved[roleName]
}

var PolicySaved map[string][]*dao.Policy

func LoadDefaultRole() {
	PolicySaved = make(map[string][]*dao.Policy)

	YSCloudComputeRole = &dao.Role{
		RoleName: "YS_CloudComputeRole",
	}
	CSPRole = &dao.Role{
		RoleName: "YS_CSPRole",
	}
	CAE365Role = &dao.Role{
		RoleName: "YS_CAE365Role",
	}
	p1 := &dao.Policy{
		PolicyName: "YS_CloudStorageAllAccess",
		Policy: dao.AuthzPolicy{
			Policy: ladon.DefaultPolicy{
				ID:        "YS_CloudStorageAllAccess",
				Resources: []string{"yrn:ys:cs:::path/<.*>"},
				Actions:   []string{"<.*>"},
				Effect:    ladon.AllowAccess,
			},
		},
	}
	p2 := &dao.Policy{
		PolicyName: "YS_HpcStorageAllAccess",
		Policy: dao.AuthzPolicy{
			Policy: ladon.DefaultPolicy{
				ID:        "YS_HpcStorageAllAccess",
				Resources: []string{"yrn:ys:cc:::path/<.*>"},
				Actions:   []string{"<.*>"},
				Effect:    ladon.AllowAccess,
			},
		},
	}
	PolicySaved[YSCloudComputeRole.RoleName] = []*dao.Policy{p1, p2}
}

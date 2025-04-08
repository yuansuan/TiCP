package fake

import (
	"sync"
	"time"

	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

const ResourceCount = 1000

type datastore struct {
	secrets             []*dao.Secret
	policies            []*dao.Policy
	roles               []*dao.Role
	policyAudits        []*dao.PolicyAudit
	rolePolicyRelations []*dao.RolePolicyRelation

	sync.RWMutex
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

func (ds *datastore) Roles() store.RoleStore {
	return newRoles(ds)
}

func (ds *datastore) PolicyAudits() store.PolicyAuditStore {
	return newPolicyAudits(ds)
}

func (ds *datastore) RolePolicyRelations() store.RolePolicyRelationStore {
	return newRolePolicyRelations(ds)
}

func (ds *datastore) MigrateDatabase() error {
	return nil
}

func (ds *datastore) Close() error {
	return nil
}

var (
	fakeFactory store.Factory
	once        sync.Once
)

func GetFakeFactory() (store.Factory, error) {
	once.Do(func() {
		fakeFactory = &datastore{
			secrets:             FakeSecrets(ResourceCount),
			policies:            FakePolicies(ResourceCount),
			roles:               FakeRoles(ResourceCount),
			rolePolicyRelations: FakeRolePolicyRelations(ResourceCount),
			policyAudits:        FakePolicyAudits(ResourceCount),
		}
	})
	if fakeFactory == nil {
		return nil, nil
	}
	return fakeFactory, nil
}

func FakeSecrets(count int) []*dao.Secret {
	var secrets []*dao.Secret

	vipBox := &dao.Secret{
		AccessKeyId:     "FL1E9NMAL7CPJUY7NJ5O",
		AccessKeySecret: "vRvkV1yWHdMOwhEVPyUTda3cLlSydyOny3kz23++",
		Expiration:      time.Now().Add(24 * time.Hour),
		ParentUser:      "4x7nt47MXA9",
	}
	secrets = append(secrets, vipBox)

	return secrets
}

func FakePolicies(count int) []*dao.Policy {
	var policies []*dao.Policy

	p1 := &dao.Policy{
		// ID:         1673584046262718464,
		UserId:     "4x7nt47MXA9",
		PolicyName: "YS_CloudStorageAllAccess",
		Policy: dao.AuthzPolicy{
			Policy: ladon.DefaultPolicy{
				ID: " : YS_CloudStorageAllAccess",
				// Subjects:  []string{"YS_CloudCompute"},
				Resources: []string{"yrn:ys:cs::4x7nt47MXA9:path/<.*>"},
				Actions:   []string{"<.*>"},
				Effect:    ladon.AllowAccess,
			},
		},
		Version: "v1",
	}
	p2 := &dao.Policy{
		ID:         1673584046262718465,
		PolicyName: "YS_HpcStorageAllAccess",
		Policy: dao.AuthzPolicy{
			Policy: ladon.DefaultPolicy{
				ID:        "YS_HpcStorageAllAccess",
				Subjects:  []string{"YS_CloudCompute"},
				Resources: []string{"yrn:ys:cc::4x7nt47MXA9:path/<.*>"},
				Actions:   []string{"<.*>"},
				Effect:    ladon.AllowAccess,
			},
		},
	}

	policies = append(policies, p1, p2)

	return policies
}

func FakeRoles(count int) []*dao.Role {
	var roles []*dao.Role

	// [0]
	policy := ladon.DefaultPolicy{
		ID:          "trust_policy",
		Description: "The policy allows VIPBoxRole_4x7nt47MXA7 to perform 'sts:AssumeRole' action",
		Subjects:    []string{"4x7nt47MXA7"},
		Resources:   []string{"yrn:ys:iam::4TiSxuPtJEm:role/VIPBoxRole_4x7nt47MXA7"},
		Actions:     []string{"STS:AssumeRole"},
		Effect:      ladon.AllowAccess,
	}

	boxRole := &dao.Role{
		// ID:          1673559263525474304,
		UserId:      "4TiSxuPtJEm",
		RoleName:    "VIPBoxRole_4x7nt47MXA7",
		TrustPolicy: dao.AuthzPolicy{Policy: policy},
	}

	// [1]
	policy1 := ladon.DefaultPolicy{
		ID:        "trust_policy",
		Subjects:  []string{"4x7nt47MXA7"},
		Resources: []string{"yrn:ys:iam::4TiSxuPtJEm:role/VIPBoxRole_4x7nt47MXA7"},
		Actions:   []string{"STS:AssumeRole"},
		Effect:    ladon.AllowAccess,
	}

	boxRole1 := &dao.Role{
		RoleName:    "VIPBoxRole_4x7nt47MXA7",
		TrustPolicy: dao.AuthzPolicy{Policy: policy1},
	}

	policy2 := ladon.DefaultPolicy{
		ID:          "trust_policy",
		Description: "This policy allows the YS_CloudComputeRole perform 4x7nt47MXA9 the 'sts:AssumeRole' action",
		Subjects:    []string{"YS_CloudCompute"},
		Resources:   []string{"yrn:ys:iam::4x7nt47MXA9:role/YS_CloudComputeRole"},
		Actions:     []string{"STS:AssumeRole"},
		Effect:      ladon.AllowAccess,
	}

	ComputeRole := &dao.Role{
		ID:          1673559263525474304,
		UserId:      "4x7nt47MXA9",
		RoleName:    "YS_CloudComputeRole",
		TrustPolicy: dao.AuthzPolicy{Policy: policy2},
	}

	roles = append(roles, boxRole, boxRole1, ComputeRole)
	return roles
}

func FakeRolePolicyRelations(count int) []*dao.RolePolicyRelation {
	var rolePolicyRelations []*dao.RolePolicyRelation

	rolePolicyRelation1 := &dao.RolePolicyRelation{
		// ID:       1673584477948874752,
		RoleId:   1673559263525474304,
		PolicyId: 1673584046262718464,
	}
	rolePolicyRelation2 := &dao.RolePolicyRelation{
		// ID:       1673584477948874753,
		RoleId:   1673559263525474304,
		PolicyId: 1673584046262718465,
	}
	rolePolicyRelations = append(rolePolicyRelations, rolePolicyRelation1, rolePolicyRelation2)
	return rolePolicyRelations
}

func FakePolicyAudits(count int) []*dao.PolicyAudit {
	var policyAudits []*dao.PolicyAudit

	policyAudit1 := &dao.PolicyAudit{
		ID:           1673584477948874752,
		Subject:      "YS_CloudCompute",
		PolicyShadow: "YS_CloudStorageAllAccess",
	}
	policyAudits = append(policyAudits, policyAudit1)
	return policyAudits
}

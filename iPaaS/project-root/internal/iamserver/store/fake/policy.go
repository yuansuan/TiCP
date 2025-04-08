package fake

import (
	"context"

	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type policies struct {
	ds *datastore
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds: ds}
}

func (p *policies) Create(ctx context.Context, policy *dao.Policy) (bool, error) {
	p.ds.Lock()
	defer p.ds.Unlock()
	policy.StatementShadow = policy.Policy.String()
	p.ds.policies = append(p.ds.policies, policy)
	return true, nil
}

func (p *policies) Update(ctx context.Context, id snowflake.ID, policy *dao.Policy) error {
	p.ds.Lock()
	defer p.ds.Unlock()
	for _, v := range p.ds.policies {
		if v.ID == id {
			v.Policy = policy.Policy
			v.StatementShadow = policy.Policy.String()
			return nil
		}
	}
	return nil
}

func (p *policies) Delete(ctx context.Context, policyID snowflake.ID) error {
	p.ds.Lock()
	defer p.ds.Unlock()
	var tmp []*dao.Policy
	for _, v := range p.ds.policies {
		if policyID != v.ID {
			tmp = append(tmp, v)
		}
	}
	p.ds.policies = tmp
	return nil
}

func (p *policies) Get(ctx context.Context, userId string, name string) (*dao.Policy, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()
	for _, v := range p.ds.policies {
		if v.UserId == userId && v.PolicyName == name {
			return v, nil
		}
	}
	return nil, nil
}

func (p *policies) GetByIds(ctx context.Context, ids []snowflake.ID) ([]*dao.Policy, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()
	var policies []*dao.Policy
	for _, v := range p.ds.policies {
		for _, id := range ids {
			if v.ID == id {
				policies = append(policies, v)
			}
		}
	}
	return policies, nil
}

func (p *policies) List(ctx context.Context, userId string, offset, limit int) ([]*dao.Policy, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()
	var policies []*dao.Policy
	for _, v := range p.ds.policies {
		if v.UserId == userId {
			policies = append(policies, v)
		}
	}
	return policies, nil
}

func (p *policies) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()
	var policies []*dao.Policy
	for _, v := range p.ds.policies {
		if v.PolicyName == key {
			policies = append(policies, v)
		}
	}
	return nil, nil
}

func (p *policies) BatchCreate(ctx context.Context, policies []*dao.Policy) error {
	p.ds.Lock()
	defer p.ds.Unlock()
	p.ds.policies = append(p.ds.policies, policies...)
	return nil
}

func (p *policies) ListByNameAndUserId(ctx context.Context, userId string, names []string) ([]*dao.Policy, error) {
	p.ds.RLock()
	defer p.ds.RUnlock()
	var policies []*dao.Policy
	for _, v := range p.ds.policies {
		if v.UserId == userId {
			for _, name := range names {
				if v.PolicyName == name {
					policies = append(policies, v)
				}
			}
		}
	}
	return policies, nil
}

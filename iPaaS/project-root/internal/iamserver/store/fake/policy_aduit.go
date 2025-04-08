package fake

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type policyAudit struct {
	ds *datastore
}

func newPolicyAudits(ds *datastore) *policyAudit {
	return &policyAudit{ds: ds}
}

func (p *policyAudit) Create(ctx context.Context, aduit *dao.PolicyAudit) error {
	p.ds.Lock()
	defer p.ds.Unlock()
	p.ds.policyAudits = append(p.ds.policyAudits, aduit)
	return nil
}

func (p *policyAudit) CleanThreeMonthAgoData(ctx context.Context) error {
	p.ds.Lock()
	defer p.ds.Unlock()
	p.ds.policyAudits = nil
	return nil
}

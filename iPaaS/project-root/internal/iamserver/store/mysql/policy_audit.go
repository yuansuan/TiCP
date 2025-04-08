package mysql

import (
	context "context"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	dao "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"gorm.io/gorm"
)

type policyAudit struct {
	db *gorm.DB
}

func newPolicyAudit(ds *datastore) *policyAudit {
	return &policyAudit{db: ds.db}
}

func (p *policyAudit) Create(ctx context.Context, policy *dao.PolicyAudit) error {
	return p.db.Create(&policy).Error
}

func (p *policyAudit) CleanThreeMonthAgoData(ctx context.Context) error {
	interval := config.GetConfig().CleanAuditLogInterval
	if interval <= 0 {
		interval = 3
	}
	return p.db.Exec("delete from policy_audit where createdAt < DATE_SUB(NOW(), INTERVAL ? DAY)", interval).Error
}

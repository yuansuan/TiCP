package store

import (
	context "context"

	dao "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type PolicyAuditStore interface {
	Create(ctx context.Context, policy *dao.PolicyAudit) error

	CleanThreeMonthAgoData(ctx context.Context) error
}

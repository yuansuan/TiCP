package store

import (
	"context"

	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	Create(ctx context.Context, policy *dao.Policy) (bool, error)
	Update(ctx context.Context, id snowflake.ID, policy *dao.Policy) error
	Delete(ctx context.Context, policyID snowflake.ID) error
	Get(ctx context.Context, userId string, name string) (*dao.Policy, error)
	GetByIds(ctx context.Context, ids []snowflake.ID) ([]*dao.Policy, error)
	List(ctx context.Context, userID string, offset, limit int) ([]*dao.Policy, error)
	GetPolicy(key string) ([]*ladon.DefaultPolicy, error)

	BatchCreate(ctx context.Context, managedPolicies []*dao.Policy) error
	ListByNameAndUserId(ctx context.Context, userId string, names []string) ([]*dao.Policy, error)
}

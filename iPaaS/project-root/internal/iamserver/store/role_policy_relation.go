package store

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type RolePolicyRelationStore interface {
	Create(ctx context.Context, relation *dao.RolePolicyRelation) (bool, error)

	Delete(ctx context.Context, relation *dao.RolePolicyRelation) error

	ListPolicyByRoleId(ctx context.Context, roleId snowflake.ID, offset, limit int) ([]snowflake.ID, error)

	CreateBatch(ctx context.Context, relations []*dao.RolePolicyRelation) error

	DeleteByRoleID(ctx context.Context, roleID snowflake.ID) error

	DeleteByPolicyID(ctx context.Context, policyID snowflake.ID) error
}

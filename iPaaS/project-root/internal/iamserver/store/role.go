package store

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type RoleStore interface {
	Create(ctx context.Context, role *dao.Role) (bool, error)
	Update(ctx context.Context, role *dao.Role, roleID snowflake.ID) error
	Delete(ctx context.Context, roleID snowflake.ID) error
	Get(ctx context.Context, userId, roleName string) (*dao.Role, error)
	List(ctx context.Context, userId string) ([]*dao.Role, error)
}

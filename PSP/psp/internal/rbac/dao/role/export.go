package role

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/consts"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
)

// Dao Dao
type Dao interface {
	Add(ctx context.Context, role *model.Role) (err error)
	GetByName(ctx context.Context, name string) (role *model.Role, err error)
	Get(ctx context.Context, id int64) (role *model.Role, err error)
	Gets(ctx context.Context, ids []int64) (roles []*model.Role, err error)
	List(ctx context.Context, request *boring.ListRequest, isAdmin bool) (roles []*model.Role, total int64, err error)
	Update(ctx context.Context, id int64, role *model.Role) error
	Delete(ctx context.Context, id int64) error
	ListByType(ctx context.Context, typ int) (roles []*model.Role, err error)
	ShouldAllExists(ctx context.Context, ids ...int64) error
}

// RoleDaoImpl RoleDaoImpl
type RoleDaoImpl struct {
	boring *boring.Dao

	adminRoleID      int64
	normalUserRoleID int64
}

// NewRoleDaoImpl NewRoleDaoImpl
func NewRoleDaoImpl() *RoleDaoImpl {
	result := &RoleDaoImpl{
		boring: boring.NewBoringDao(
			&model.Role{},
			boring.Codes{
				errcode.ErrRBACRoleNotFound,
				errcode.ErrRBACRoleAlreadyExist,
				errcode.ErrRBACInsertError,
				errcode.ErrRBACQueryError,
				errcode.ErrRBACDeleteError,
				errcode.ErrRBACUpdateError,
				errcode.ErrRBACInvalidOrderError,
				errcode.ErrRBACPageShouldPositive,
			},
			"name",
			[]string{"id", "type"},
			[]string{"name", "id"},
		),
	}
	return result
}

// Add Add
func (b *RoleDaoImpl) Add(ctx context.Context, role *model.Role) (err error) {
	return b.boring.Add(ctx, role)
}

// Get Get
func (b *RoleDaoImpl) Get(ctx context.Context, id int64) (role *model.Role, err error) {
	role = new(model.Role)
	err = b.boring.Get(ctx, id, role)
	return
}

// Get Get
func (b *RoleDaoImpl) GetByName(ctx context.Context, roleName string) (role *model.Role, err error) {
	role = new(model.Role)
	err = b.boring.GetByName(ctx, roleName, role)
	return
}

// Gets Gets
func (b *RoleDaoImpl) Gets(ctx context.Context, ids []int64) (roles []*model.Role, err error) {
	roles = make([]*model.Role, 0)
	err = b.boring.Gets(ctx, ids, &roles)
	return
}

// List List
func (b *RoleDaoImpl) List(ctx context.Context, listReq *boring.ListRequest, isAdmin bool) (roles []*model.Role, total int64, err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	if listReq.OrderBy != "" && len(b.boring.AllowOrderBy) != 0 && !b.boring.AllowOrderBy[listReq.OrderBy] {
		return nil, 0, status.Errorf(b.boring.Codes.InvalidOrderError, "unsupported order [%v], only support: %v",
			listReq.OrderBy, b.boring.AllowOrderBy)
	}

	if listReq.NameFilter != "" {
		session = session.Where("name LIKE ?", fmt.Sprintf("%%%s%%", listReq.NameFilter))
	}

	if !isAdmin {
		session = session.Where("type != ?", consts.RoleTypeSuperAdmin)
	}

	// if pageSize <= 0, no limit
	if listReq.PageSize > 0 {
		if listReq.Page < 1 {
			return nil, 0, status.Error(b.boring.Codes.PageShouldPositive, "")
		}
		session = session.Limit(int(listReq.PageSize), int(listReq.Page-1)*int(listReq.PageSize))
	}

	if listReq.OrderBy != "" {
		if listReq.Desc {
			session.Desc(listReq.OrderBy)
		} else {
			session.Asc(listReq.OrderBy)
		}
	}

	total, err = session.FindAndCount(&roles)

	return
}

// Update Update
func (b *RoleDaoImpl) Update(ctx context.Context, id int64, role *model.Role) error {
	return b.boring.Update(ctx, id, role)
}

// Delete Delete
func (b *RoleDaoImpl) Delete(ctx context.Context, id int64) error {
	return b.boring.Delete(ctx, id)
}

// ShouldAllExists ShouldAllExists
func (b *RoleDaoImpl) ShouldAllExists(ctx context.Context, ids ...int64) error {
	return b.boring.AllExists(ctx, ids...)
}

// ListByType ListByType
func (b *RoleDaoImpl) ListByType(ctx context.Context, typ int) (roles []*model.Role, err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	roles = []*model.Role{}

	err = session.Where("type=?", typ).Find(&roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

package resource

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc/status"
	"xorm.io/xorm"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/errcode"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/proto/rbac"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/boring"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/dao/model"
	"github.com/yuansuan/ticp/PSP/psp/internal/rbac/util"
)

// Dao Dao
type Dao interface {
	Add(ctx context.Context, res *model.Resource) (err error)
	Get(ctx context.Context, id int64) (res *model.Resource, err error)
	Gets(ctx context.Context, ids []int64) (resList []*model.Resource, err error)
	GetCustomResList(ctx context.Context, ids []int64) ([]*model.Resource, error)
	List(ctx context.Context, request *boring.ListRequest, custom int32) (resList []*model.Resource, total int64, err error)
	Update(ctx context.Context, id int64, res *model.Resource) error
	Delete(ctx context.Context, id int64) error
	AllExists(ctx context.Context, ids ...int64) error
	// FindResourceByNameOrExId 通过资源名称或外部关联id获取资源信息
	FindResourceByNameOrExId(ctx context.Context, ress ...*rbac.ResourceIdentity) (result []*model.Resource, notFound []*rbac.ResourceIdentity, err error)
	ListResourceTypePermissions(ctx context.Context, resourceType ...string) (perms []*model.Resource, err error)
	NoInternalPerm(ctx context.Context, ids ...int64) error

	FindChildResIds(ctx context.Context, ids []int64) ([]int64, error)
}

// DaoImpl DaoImpl
type ResourceDaoImpl struct {
	boring *boring.Dao
}

// NewResourceDaoImpl NewResourceDaoImpl
func NewResourceDaoImpl() *ResourceDaoImpl {
	return &ResourceDaoImpl{
		boring: boring.NewBoringDao(
			&model.Resource{},
			boring.Codes{
				errcode.ErrRBACPermissionNotFound,
				errcode.ErrRBACPermissionAlreadyExist,
				errcode.ErrRBACInsertError,
				errcode.ErrRBACQueryError,
				errcode.ErrRBACDeleteError,
				errcode.ErrRBACUpdateError,
				errcode.ErrRBACInvalidOrderError,
				errcode.ErrRBACPageShouldPositive,
			},
			"display_name",
			[]string{"id", "type", "custom"},
			[]string{}),
	}
}

// Add Add
func (b *ResourceDaoImpl) Add(ctx context.Context, res *model.Resource) (err error) {
	return b.boring.Add(ctx, res)
}

// Get Get
func (b *ResourceDaoImpl) Get(ctx context.Context, id int64) (res *model.Resource, err error) {
	res = new(model.Resource)
	res = &model.Resource{}
	err = b.boring.Get(ctx, id, res)
	return
}

// Gets Gets
func (b *ResourceDaoImpl) Gets(ctx context.Context, ids []int64) (resList []*model.Resource, err error) {
	err = b.boring.Gets(ctx, ids, &resList)
	return
}

// Gets Gets
func (b *ResourceDaoImpl) GetCustomResList(ctx context.Context, ids []int64) ([]*model.Resource, error) {
	resList := []*model.Resource{}

	session := middleware.DefaultSession(ctx)
	defer session.Close()

	err := session.In("id", ids).Where("`custom` = 1").Find(&resList)

	if err != nil {
		err = status.Error(errcode.ErrRBACQueryError, "")
	}
	return resList, nil

}

// List List
func (b *ResourceDaoImpl) List(ctx context.Context, request *boring.ListRequest, custom int32) (resList []*model.Resource, total int64, err error) {

	sql := middleware.DefaultSession(ctx)
	defer sql.Close()

	if request.NameFilter != "" {
		sql = sql.Where("display_name LIKE ?", "%"+request.NameFilter+"%")
	}

	if custom != 0 {
		sql = sql.Where("custom = ?", custom)
	}

	// if pageSize <= 0, no limit
	if request.PageSize > 0 {
		if request.Page < 1 {
			total = 0
			err = status.Error(errcode.ErrRBACPageShouldPositive, "")
		}
		sql = sql.Limit(int(request.PageSize), int(request.Page-1)*int(request.PageSize))
	}

	if request.OrderBy != "" {
		if request.Desc {
			sql.Desc(request.OrderBy)
		} else {
			sql.Asc(request.OrderBy)
		}
	}

	total, err = sql.FindAndCount(&resList)

	return
}

// Update Update
func (b *ResourceDaoImpl) Update(ctx context.Context, id int64, res *model.Resource) error {
	return b.boring.Update(ctx, id, res)
}

// Delete Delete
func (b *ResourceDaoImpl) Delete(ctx context.Context, id int64) error {
	return b.boring.Delete(ctx, id)
}

// AllExists AllExists
func (b *ResourceDaoImpl) AllExists(ctx context.Context, ids ...int64) error {
	return b.boring.AllExists(ctx, ids...)
}

// ListResourceTypePermissions ListResourceTypePermissions
func (b *ResourceDaoImpl) ListResourceTypePermissions(ctx context.Context, resourceType ...string) (resList []*model.Resource, err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	resList = []*model.Resource{}
	err = session.In("type", resourceType).Find(&resList)
	if err != nil {
		return nil, err
	}
	return resList, err
}

// FindPermissionsByResource FindPermissionsByResource
// notFound is index in resources
func (b *ResourceDaoImpl) FindResourceByNameOrExId(ctx context.Context, resources ...*rbac.ResourceIdentity) (result []*model.Resource, notFound []*rbac.ResourceIdentity, err error) {

	session := middleware.DefaultSession(ctx)
	defer session.Close()

	result = []*model.Resource{}

	for _, resource := range resources {
		query := (*xorm.Session)(nil)

		switch v := resource.Identity.(type) {
		case *rbac.ResourceIdentity_Id:
			query = session.Where("`type`=? AND `external_id`=?", v.Id.Type, v.Id.Id)
		case *rbac.ResourceIdentity_Name:
			if v.Name.Name == "" {
				v.Name.Name = common.ResourceActionNONE
			}
			if v.Name.Action == "" {
				v.Name.Action = "NONE"
			}
			query = session.Where("`type`=? AND `name`=? and `action`=?", v.Name.Type, v.Name.Name, v.Name.Action)
		default:
			return nil, nil, status.Errorf(errcode.ErrRBACUnknownResourceIdentify, "type: %T; value: %v", v, v)
		}

		var records []*model.Resource
		err = query.Find(&records)
		if err != nil {
			return nil, nil, status.Error(errcode.ErrRBACQueryError, err.Error())
		}

		if len(records) == 0 {
			notFound = append(notFound, resource)
			continue
		}

		if len(records) > 1 {
			// just use first elem
			logging.GetLogger(ctx).Warnf("resource should unique by (type, name, action) or (type, external_id), one type should only choose one method,"+
				"using first elem, resource_id: %T :: %v, result with len=%v :: %+v", resource.Identity, resource.Identity, len(records), records)
		}

		result = append(result, records[0])
	}

	return result, notFound, nil
}

// NoInternalPerm NoInternalPerm
func (b *ResourceDaoImpl) NoInternalPerm(ctx context.Context, ids ...int64) error {

	if len(ids) <= 0 {
		return nil
	}

	session := middleware.DefaultSession(ctx)
	defer session.Close()

	count, err := session.In("id", ids).Where("`type` = 'internal'").Count(&model.Resource{})
	if err != nil {
		return status.Error(errcode.ErrRBACQueryError, err.Error())
	}

	if count != 0 {
		return status.Error(errcode.ErrRBACInternalPermOnlyGiveAdmin, fmt.Sprintf("[%d/%d] is internal perm", count, len(ids)))
	}

	return nil
}

// FindChildResIds 批量查询子级权限id(包含父级id)
func (b *ResourceDaoImpl) FindChildResIds(ctx context.Context, ids []int64) (allIds []int64, err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	if len(ids) <= 0 {
		return ids, nil
	}

	for _, id := range ids {
		childIds, _ := b.FindChildResId(ctx, id)
		if len(childIds) <= 0 {
			continue
		}
		allIds = append(allIds, childIds...)
	}
	allIds = append(allIds, ids...)
	allIds = util.RemoveDuplicates(allIds)
	return
}

// FindChildResId 查询子级权限id
func (b *ResourceDaoImpl) FindChildResId(ctx context.Context, id int64) (childIds []int64, err error) {
	session := middleware.DefaultSession(ctx)
	defer session.Close()

	err = session.SQL(`select id
								from (select t1.id,
											 if(find_in_set(parent_id, @pids) > 0, @pids := concat(@pids, ',', id), 0) as ischild
									  from (select id, parent_id from resource t order by parent_id, id) t1,
										   (select @pids := ?) t2) t3
								where ischild != 0 group by id`, id).Find(&childIds)

	return
}

package fake

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type roles struct {
	ds *datastore
}

func newRoles(ds *datastore) *roles {
	return &roles{ds: ds}
}

func (r *roles) Create(ctx context.Context, role *dao.Role) (bool, error) {
	r.ds.Lock()
	defer r.ds.Unlock()
	r.ds.roles = append(r.ds.roles, role)
	return true, nil
}
func (r *roles) Delete(ctx context.Context, roleID snowflake.ID) error {
	r.ds.Lock()
	defer r.ds.Unlock()
	var tmp []*dao.Role
	for _, v := range r.ds.roles {
		if roleID != v.ID {
			tmp = append(tmp, v)
		}
	}
	r.ds.roles = tmp
	return nil
}
func (r *roles) Get(ctx context.Context, userId, roleName string) (*dao.Role, error) {
	r.ds.RLock()
	defer r.ds.RUnlock()
	for _, v := range r.ds.roles {
		if v.UserId == userId && v.RoleName == roleName {
			return v, nil
		}
	}
	return nil, nil
}

func (r *roles) List(ctx context.Context, userId string) ([]*dao.Role, error) {
	r.ds.RLock()
	defer r.ds.RUnlock()
	var roles []*dao.Role
	for _, v := range r.ds.roles {
		if v.UserId == userId {
			roles = append(roles, v)
		}
	}
	return roles, nil
}

func (r *roles) Update(ctx context.Context, role *dao.Role, id snowflake.ID) error {
	r.ds.Lock()
	defer r.ds.Unlock()
	for i, v := range r.ds.roles {
		if v.ID == id {
			r.ds.roles[i] = role
			return nil
		}
	}
	return nil
}

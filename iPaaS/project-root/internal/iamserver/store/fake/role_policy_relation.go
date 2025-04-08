package fake

import (
	"context"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
)

type rolePolicyRelations struct {
	ds *datastore
}

func newRolePolicyRelations(ds *datastore) *rolePolicyRelations {
	return &rolePolicyRelations{ds: ds}
}

func (r *rolePolicyRelations) Create(ctx context.Context, relation *dao.RolePolicyRelation) (bool, error) {
	r.ds.Lock()
	defer r.ds.Unlock()
	r.ds.rolePolicyRelations = append(r.ds.rolePolicyRelations, relation)
	return true, nil
}
func (r *rolePolicyRelations) ListPolicyByRoleId(ctx context.Context, roleId snowflake.ID, offset, limit int) ([]snowflake.ID, error) {
	r.ds.RLock()
	defer r.ds.RUnlock()
	var ids []snowflake.ID
	for _, v := range r.ds.rolePolicyRelations {
		if v.RoleId == roleId {
			ids = append(ids, v.PolicyId)
		}
	}
	return ids, nil
}

func (r *rolePolicyRelations) CreateBatch(ctx context.Context, relations []*dao.RolePolicyRelation) error {
	r.ds.Lock()
	defer r.ds.Unlock()
	r.ds.rolePolicyRelations = append(r.ds.rolePolicyRelations, relations...)
	return nil
}

func (r *rolePolicyRelations) DeleteByRoleID(ctx context.Context, roleID snowflake.ID) error {
	r.ds.Lock()
	defer r.ds.Unlock()
	var tmp []*dao.RolePolicyRelation
	for _, v := range r.ds.rolePolicyRelations {
		if roleID != v.RoleId {
			tmp = append(tmp, v)
		}
	}
	r.ds.rolePolicyRelations = tmp
	return nil
}

func (r *rolePolicyRelations) DeleteByPolicyID(ctx context.Context, policyID snowflake.ID) error {
	r.ds.Lock()
	defer r.ds.Unlock()
	var tmp []*dao.RolePolicyRelation
	for _, v := range r.ds.rolePolicyRelations {
		if v.PolicyId != policyID {
			tmp = append(tmp, v)
		}
	}
	r.ds.rolePolicyRelations = tmp
	return nil
}

func (r *rolePolicyRelations) Delete(ctx context.Context, relation *dao.RolePolicyRelation) error {
	r.ds.Lock()
	defer r.ds.Unlock()
	var tmp []*dao.RolePolicyRelation
	for _, v := range r.ds.rolePolicyRelations {
		if v.PolicyId != relation.PolicyId || v.RoleId != relation.RoleId {
			tmp = append(tmp, v)
		}
	}
	r.ds.rolePolicyRelations = tmp
	return nil
}

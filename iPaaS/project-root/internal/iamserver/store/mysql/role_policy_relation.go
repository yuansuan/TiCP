package mysql

import (
	"context"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"gorm.io/gorm"
)

type rolePolicyRelations struct {
	db *gorm.DB
}

func newRolePolicyRelation(ds *datastore) *rolePolicyRelations {
	return &rolePolicyRelations{ds.db}
}

func (r *rolePolicyRelations) Create(ctx context.Context, relation *dao.RolePolicyRelation) (bool, error) {
	relation.ID = common.IdGen.Generate()
	err := r.db.Create(relation).Error
	if isDuplicateKeyErr(err) {
		return true, nil
	}
	return false, err
}

func (r *rolePolicyRelations) Delete(ctx context.Context, relation *dao.RolePolicyRelation) error {
	sql := "delete from role_policy_relation where roleId = ? and policyId = ?"
	i := r.db.Exec(sql, relation.RoleId, relation.PolicyId).RowsAffected
	if i == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

func (r *rolePolicyRelations) ListPolicyByRoleId(ctx context.Context, roleId snowflake.ID, offset, limit int) ([]snowflake.ID, error) {
	var ret []*dao.RolePolicyRelation
	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	d := r.db.Where("roleId = ? ", roleId).
		Offset(offset).
		Limit(limit).
		Order("id desc").
		Find(&ret)
	var pIds []snowflake.ID
	for _, v := range ret {
		pIds = append(pIds, v.PolicyId)
	}
	return pIds, d.Error
}

func (r *rolePolicyRelations) CreateBatch(ctx context.Context, relations []*dao.RolePolicyRelation) error {
	// return r.db.Create(&relations).Error
	setTran(ctx, r)
	values := []interface{}{}
	sql := "insert into role_policy_relation (id, roleId, policyId) values "
	for _, v := range relations {
		v.ID = common.IdGen.Generate()
		sql += "(?, ?, ?),"
		values = append(values, v.ID, v.RoleId, v.PolicyId)
	}
	sql = strings.TrimSuffix(sql, ",") // remove the trailing comma
	err := r.db.Exec(sql, values...).Error
	return err
}

func (r *rolePolicyRelations) DeleteByRoleID(ctx context.Context, roleID snowflake.ID) error {
	setTran(ctx, r)
	sql := "delete from role_policy_relation where roleId = ?"
	err := r.db.Exec(sql, roleID).Error
	return err
}

func (r *rolePolicyRelations) DeleteByPolicyID(ctx context.Context, policyId snowflake.ID) error {
	setTran(ctx, r)
	sql := "delete from role_policy_relation where policyId = ?"
	err := r.db.Exec(sql, policyId).Error
	return err
}

func setTran(ctx context.Context, r *rolePolicyRelations) {
	v := ctx.Value(common.ContextTransactionKey)
	if v != nil {
		tranDB := v.(*gorm.DB)
		if tranDB != nil {
			r.db = tranDB
		}
	}
}

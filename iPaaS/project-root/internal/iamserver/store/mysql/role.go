package mysql

import (
	"context"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/code"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"

	"github.com/go-sql-driver/mysql"
	"github.com/marmotedu/errors"

	"gorm.io/gorm"
)

type roles struct {
	db *gorm.DB
}

func newRoles(ds *datastore) *roles {
	return &roles{ds.db}
}

// Create creates a new role.
func (r *roles) Create(ctx context.Context, role *dao.Role) (bool, error) {
	role.ID = common.IdGen.Generate()
	err := role.BeforeCreateForRaw()
	if err != nil {
		logging.Default().Errorf("failed to marshal policy: %v", err)
		return false, errors.WithCode(code.ErrDatabase, err.Error())
	}
	setRoleTran(ctx, r)
	// easy to mock
	err = r.db.Exec("insert into role (id, userId, roleName, description, trustPolicyShadow) values (?, ?, ?, ?, ?)",
		role.ID, role.UserId, role.RoleName, role.Description, role.TrustPolicyShadow).Error
	if isDuplicateKeyErr(err) {
		return true, nil
	}
	return false, err
}

func isDuplicateKeyErr(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		// uniqueConstraint
		if mysqlErr.Number == 1062 {
			return true
		}
	}
	return false
}

func (r *roles) Update(ctx context.Context, role *dao.Role, roleID snowflake.ID) error {
	err := role.BeforeCreateForRaw()
	if err != nil {
		logging.Default().Errorf("failed to marshal policy: %v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	err = r.db.Exec("update role set description = ?, trustPolicyShadow = ? where id = ?",
		role.Description, role.TrustPolicyShadow, roleID).Error
	return err
}

// Delete deletes the policy by the policy identifier.
func (r *roles) Delete(ctx context.Context, roleID snowflake.ID) error {
	setRoleTran(ctx, r)
	row := r.db.Exec("DELETE FROM role WHERE id = ?", roleID).RowsAffected
	if row == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

// Get return policy by the policy identifier.
func (r *roles) Get(ctx context.Context, userid, name string) (*dao.Role, error) {
	setRoleTran(ctx, r)
	role := &dao.Role{}
	err := r.db.Where("userId = ? and roleName = ?", userid, name).First(role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	err = role.UnmarshalPolicy()
	if err != nil {
		logging.Default().Errorf("failed to unmarshal policy: %v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return role, nil
}

func (r *roles) List(ctx context.Context, userId string) ([]*dao.Role, error) {
	var ret []*dao.Role
	d := r.db.Where("userId = ?", userId).Find(&ret)
	return ret, d.Error
}

func setRoleTran(ctx context.Context, r *roles) {
	v := ctx.Value(common.ContextTransactionKey)
	if v != nil {
		tranDB := v.(*gorm.DB)
		if tranDB != nil {
			r.db = tranDB
		}
	}
}

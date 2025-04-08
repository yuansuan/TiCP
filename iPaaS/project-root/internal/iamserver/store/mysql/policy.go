package mysql

import (
	"context"
	"strings"

	"github.com/marmotedu/errors"
	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/code"
	"gorm.io/gorm"
)

type policies struct {
	db *gorm.DB
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds.db}
}

// Create creates a new ladon policy.
func (p *policies) Create(ctx context.Context, policy *dao.Policy) (bool, error) {
	policy.ID = common.IdGen.Generate()
	err := p.db.Create(&policy).Error
	if isDuplicateKeyErr(err) {
		return true, nil
	}
	return false, err
}

// Update updates policy by the policy identifier.
func (p *policies) Update(ctx context.Context, id snowflake.ID, policy *dao.Policy) error {
	sql := "update policy set statementShadow = ?, version = ? where id = ?"
	err := p.db.Exec(sql, policy.StatementShadow, policy.Version, id).Error
	return err
}

// Delete deletes the policy by the policy identifier.
func (p *policies) Delete(ctx context.Context, policyID snowflake.ID) error {
	setPolicyTran(ctx, p)
	row := p.db.Exec("DELETE FROM policy WHERE id = ?", policyID).RowsAffected
	if row == 0 {
		return common.ErrRecordNotFound
	}
	return nil
}

// DeleteCollection batch deletes policies by policies ids.
func (p *policies) DeleteCollection(
	ctx context.Context,
	username string,
	names []string,
) error {
	return p.db.Where("userId = ? and name in (?)", username, names).Delete(dao.Policy{}).Error
}

// Get return policy by the policy identifier.
func (p *policies) Get(ctx context.Context, userid, name string) (*dao.Policy, error) {
	setPolicyTran(ctx, p)
	policy := dao.Policy{}
	err := p.db.Where("userId = ? and policyName = ?", userid, name).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrRecordNotFound
		}
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return &policy, nil
}

func (p *policies) GetByIds(ctx context.Context, ids []snowflake.ID) ([]*dao.Policy, error) {
	var ret []*dao.Policy
	err := p.db.Where("id in (?)", ids).Order("id desc").Find(&ret).Error
	return ret, err
}

// List return all policies.
func (p *policies) List(ctx context.Context, userID string, offset, limit int) ([]*dao.Policy, error) {
	var ret []*dao.Policy
	if userID != "" {
		p.db = p.db.Where("userId = ?", userID)
	}
	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 1000 {
		limit = 1000
	}
	d := p.db.
		Offset(offset).
		Limit(limit).
		Order("id desc").
		Find(&ret)
	return ret, d.Error
}

func (p *policies) GetPolicy(key string) ([]*ladon.DefaultPolicy, error) {
	var ret []*ladon.DefaultPolicy
	err := p.db.Where("policyName = ?", key).Find(&ret).Error
	return ret, err
}

func (p *policies) BatchCreate(ctx context.Context, managedPolicies []*dao.Policy) error {
	// return p.db.Create(&managedPolicies).Error
	setPolicyTran(ctx, p)
	values := []interface{}{}
	sql := "insert into policy (id, userId, policyName, statementShadow) values "
	for _, policy := range managedPolicies {
		policy.ID = common.IdGen.Generate()
		err := policy.BeforeCreateForRaw()
		if err != nil {
			logging.Default().Errorf("policy.BeforeCreateForRaw error: %v", err)
			return err
		}
		sql += "(?, ?, ?, ?),"
		values = append(values, policy.ID, policy.UserId, policy.PolicyName, policy.StatementShadow)
	}
	sql = strings.TrimSuffix(sql, ",") // remove the trailing comma
	err := p.db.Exec(sql, values...).Error
	return err
}

func (p *policies) ListByNameAndUserId(ctx context.Context, userId string, names []string) ([]*dao.Policy, error) {
	setPolicyTran(ctx, p)
	var ret []*dao.Policy
	// err := p.db.Where("userId = ? and policyName in (?)", userId, names).Find(&ret).Error
	err := p.db.Raw("SELECT id,  policyName FROM policy where userId = ? and policyName in (?)", userId, names).Scan(&ret).Error
	return ret, err
}

func setPolicyTran(ctx context.Context, p *policies) {
	v := ctx.Value(common.ContextTransactionKey)
	if v != nil {
		tranDB := v.(*gorm.DB)
		if tranDB != nil {
			p.db = tranDB
		}
	}
}

package v1

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/code"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/mysql"

	"github.com/marmotedu/errors"
	"github.com/ory/ladon"
)

type PolicySvc interface {
	GetByPolicyName(c context.Context, userID, policyName string) (*dao.Policy, error)
	AddPolicy(c context.Context, userID string, policy *iam_api.AddPolicyRequest) error
	UpdatePolicy(c context.Context, userID string, policy *iam_api.AddPolicyRequest) error
	DeletePolicy(c context.Context, userID, policyName string) error
	ListPolicy(c context.Context, userID string, offset, limit int) ([]*dao.Policy, error)
}

type policyService struct {
	store store.Factory
}

var _ PolicySvc = (*policyService)(nil)

func newPolicies(s *svc) *policyService {
	return &policyService{
		store: s.store,
	}
}

func (p *policyService) GetByPolicyName(c context.Context, userID, policyName string) (*dao.Policy, error) {
	policy, err := p.store.Policies().Get(c, userID, policyName)
	if err != nil {
		logging.Default().Infof("get policy error: %v", err)
		return nil, err
	}
	return policy, nil
}

func (p *policyService) AddPolicy(c context.Context, userID string, req *iam_api.AddPolicyRequest) error {
	policy := convertStatementToPolicy(req)
	return p.addPolicy(c, userID, req.PolicyName, req.Version, policy)
}

func convertStatementToPolicy(s *iam_api.AddPolicyRequest) dao.AuthzPolicy {
	p := ladon.DefaultPolicy{
		ID:          fmt.Sprintf(s.PolicyName),
		Description: "",
		Resources:   s.Resources,
		Actions:     s.Actions,
		Effect:      s.Effect,
	}
	policy := dao.AuthzPolicy{
		Policy: p,
	}
	return policy
}

func (p *policyService) addPolicy(c context.Context, userID, policyName, v string, policy dao.AuthzPolicy) error {
	po := &dao.Policy{
		UserId:     userID,
		PolicyName: policyName,
		Version:    v,
		Policy:     policy,
	}
	exists, err := p.store.Policies().Create(c, po)
	if err != nil {
		logging.Default().Infof("create policy error: %v", err)
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	if exists {
		logging.Default().Infof("policy already exists")
		return common.ErrAlreadyExists
	}
	return nil
}

func (p *policyService) UpdatePolicy(c context.Context, userID string, req *iam_api.AddPolicyRequest) error {
	policy, err := p.store.Policies().Get(c, userID, req.PolicyName)
	if err != nil {
		logging.Default().Infof("get policy error: %v", err)
		return err
	}
	newPolicy := convertStatementToPolicy(req)
	s := newPolicy.String()

	newP := &dao.Policy{
		ID:              policy.ID,
		UserId:          userID,
		PolicyName:      req.PolicyName,
		StatementShadow: s,
		Version:         req.Version,
	}

	err = p.store.Policies().Update(c, policy.ID, newP)
	if err != nil {
		logging.Default().Infof("update policy error: %v", err)
		return err
	}
	return nil
}

func (p *policyService) DeletePolicy(c context.Context, userID, policyName string) error {
	tran, err := mysql.GetDB(p.store)
	if err != nil {
		logging.Default().Infof("DeletePolicy get db error: %v", err)
		return err
	}
	tx := tran.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		logging.Default().Infof("DeletePolicy begin transaction error: %v", err)
		return err
	}
	ctx := context.WithValue(c, common.ContextTransactionKey, tx)
	policy, err := p.store.Policies().Get(ctx, userID, policyName)
	if err != nil {
		logging.Default().Infof("get policy error: %v", err)
		tx.Rollback()
		return err
	}
	err = p.store.Policies().Delete(ctx, policy.ID)
	if err != nil {
		logging.Default().Infof("delete policy error: %v", err)
		tx.Rollback()
		return err
	}
	err = p.store.RolePolicyRelations().DeleteByPolicyID(ctx, policy.ID)
	if err != nil {
		logging.Default().Infof("delete role policy relation error: %v", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (p *policyService) ListPolicy(c context.Context, userID string, offset, limit int) ([]*dao.Policy, error) {
	policies, err := p.store.Policies().List(c, userID, offset, limit)
	if err != nil {
		logging.Default().Infof("list policy error: %v", err)
		return nil, errors.WithCode(code.ErrDatabase, err.Error())
	}
	return policies, nil
}

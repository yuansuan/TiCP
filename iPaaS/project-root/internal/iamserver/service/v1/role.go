package v1

import (
	"context"
	"fmt"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/mysql"

	// "github.com/yuansuan/ticp/common/project-root-iam/pkg/common"
	"github.com/ory/ladon"
)

type RoleSvc interface {
	AddRole(c context.Context, userID string, req *iam_api.AddRoleRequest) error
	UpdateRole(c context.Context, userID string, req *iam_api.Role) error
	GetRole(c context.Context, userID string, roleName string) (*dao.Role, error)
	Delete(c context.Context, userID string, roleName string) error
	ListRole(c context.Context, userID string) ([]*dao.Role, error)

	PatchPolicy(c context.Context, role *dao.Role, policy *dao.Policy) error

	DetachPolicy(c context.Context, role *dao.Role, policy *dao.Policy) error
}

type roleService struct {
	store store.Factory
}

var _ RoleSvc = (*roleService)(nil)

func newRoles(s *svc) *roleService {
	return &roleService{
		store: s.store,
	}
}

func (r *roleService) AddRole(c context.Context, userID string, req *iam_api.AddRoleRequest) error {
	subjects := req.TrustPolicy.Principals
	effect := req.TrustPolicy.Effect
	resources := req.TrustPolicy.Resources
	if len(resources) == 0 {
		resources = []string{fmt.Sprintf("yrn:ys:iam::%s:role/%s", userID, req.RoleName)}
	}

	policy := ladon.DefaultPolicy{
		ID:          "trust_policy",
		Description: fmt.Sprintf("The policy allows %s to perform 'sts:AssumeRole' action", req.RoleName),
		Subjects:    subjects,
		Resources:   resources,
		Actions:     []string{StsAction},
		Effect:      effect,
	}
	trustPolicy := dao.AuthzPolicy{
		Policy: policy,
	}
	role := &dao.Role{
		UserId:      userID,
		RoleName:    req.RoleName,
		Description: req.Description,
		TrustPolicy: trustPolicy,
	}

	exists, err := r.store.Roles().Create(c, role)
	if exists {
		return common.ErrAlreadyExists
	}
	if err != nil {
		logging.Default().Infof("AddRole error: %v", err)
		return err
	}
	return nil
}

func (r *roleService) UpdateRole(c context.Context, userID string, req *iam_api.Role) error {
	role, err := r.store.Roles().Get(c, userID, req.RoleName)
	if err != nil {
		logging.Default().Infof("UpdateRole get role error: %v", err)
		return err
	}

	trustPolicy := role.TrustPolicy
	trustPolicy.Policy.Description = req.Description
	trustPolicy.Policy.Subjects = req.TrustPolicy.Principals
	trustPolicy.Policy.Resources = req.TrustPolicy.Resources
	trustPolicy.Policy.Actions = req.TrustPolicy.Actions
	trustPolicy.Policy.Effect = req.TrustPolicy.Effect

	updateRole := &dao.Role{
		RoleName:    req.RoleName,
		Description: req.Description,
		TrustPolicy: trustPolicy,
	}
	err = r.store.Roles().Update(c, updateRole, role.ID)
	if err != nil {
		logging.Default().Infof("UpdateRole error: %v", err)
		return err
	}
	return nil
}

func (r *roleService) GetRole(c context.Context, userID string, roleName string) (*dao.Role, error) {
	role, err := r.store.Roles().Get(c, userID, roleName)
	if err != nil {
		logging.Default().Infof("GetRole error: %v", err)
		return nil, err
	}
	return role, nil
}

func (r *roleService) Delete(c context.Context, userID string, roleName string) error {
	tran, err := mysql.GetDB(r.store)
	if err != nil {
		logging.Default().Infof("DeleteRole get db error: %v", err)
		return err
	}
	tx := tran.Begin()
	defer func() {
		if re := recover(); re != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		logging.Default().Infof("DeleteRole begin transaction error: %v", err)
		return err
	}
	ctx := context.WithValue(c, common.ContextTransactionKey, tx)
	role, err := r.store.Roles().Get(ctx, userID, roleName)
	if err != nil {
		logging.Default().Infof("DeleteRole error: %v", err)
		tx.Rollback()
		return err
	}
	err = r.store.Roles().Delete(ctx, role.ID)
	if err != nil {
		logging.Default().Infof("DeleteRole error: %v", err)
		tx.Rollback()
		return err
	}
	err = r.store.RolePolicyRelations().DeleteByRoleID(ctx, role.ID)
	if err != nil {
		logging.Default().Infof("DeleteRolePolicyRelations error: %v", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *roleService) ListRole(c context.Context, userID string) ([]*dao.Role, error) {
	roles, err := r.store.Roles().List(c, userID)
	if err != nil {
		logging.Default().Infof("ListRole error: %v", err)
		return nil, err
	}
	return roles, nil
}

func (r *roleService) PatchPolicy(c context.Context, role *dao.Role, policy *dao.Policy) error {
	relation := &dao.RolePolicyRelation{
		RoleId:   role.ID,
		PolicyId: policy.ID,
	}

	exists, err := r.store.RolePolicyRelations().Create(c, relation)
	if exists {
		return common.ErrAlreadyExists
	}
	if err != nil {
		logging.Default().Infof("PatchPolicy error: %v", err)
		return err
	}
	return nil
}

func (r *roleService) DetachPolicy(c context.Context, role *dao.Role, policy *dao.Policy) error {
	relation := &dao.RolePolicyRelation{
		RoleId:   role.ID,
		PolicyId: policy.ID,
	}
	err := r.store.RolePolicyRelations().Delete(c, relation)
	if err != nil {
		logging.Default().Infof("DetachPolicy error: %v", err)
		return err
	}
	return nil
}

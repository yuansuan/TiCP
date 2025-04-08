package v1

import (
	context "context"
	"fmt"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	"net/http"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/authorization"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/code"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common/snowflake"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/mysql"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"github.com/ory/ladon"
	"github.com/samber/lo"
)

const (
	// Role Claim key
	roleYrnClaim = "roleYrn"

	StsAction = "STS:AssumeRole"
)

// Service defines functions used to return resource interface.
type Service interface {
	AssumeRole(c *gin.Context)
	IsAllow(c *gin.Context)
}

type service struct {
	store store.Factory
	p     authorization.PolicyCheck
}

// NewService returns Service interface.
func NewService(factory store.Factory) Service {
	LoadDefaultRole()
	svc := &service{
		store: factory,
		p:     authorization.NewAuthorizer(factory),
	}
	return svc
}

func (s *service) AssumeRole(c *gin.Context) {
	req := iam_api.AssumeRoleRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, "InvalidArgument", "invalid request body")
		return
	}
	if len(req.RoleYrn) == 0 {
		common.ErrorResp(c, http.StatusBadRequest, "InvalidArgument", "roleYrn is required")
		return
	}
	if req.DurationSeconds < 0 {
		common.ErrorResp(c, http.StatusBadRequest, "InvalidArgument", "durationSeconds should be positive")
		return
	}
	yrn, err := s.tryAddPlatformRole(c, req.RoleYrn)
	if err != nil {
		// if err happens, tryAddPlatformRole has already handled it
		return
	}
	userID := yrn.AccountID
	roleID := yrn.ResourceID
	var r *dao.Role
	r, err = s.store.Roles().Get(c, userID, roleID)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorResp(c, http.StatusNotFound, "NotFound.Role", "role not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	userInfo := common.GetUserInfo(c)
	userId := userInfo.UserID.String()
	if IsYuansuanProductAccount(userInfo.Tag) {
		// 远算自己的云产品，采用别名
		// 账号的Tag不能乱打，要对应于远算产品名称
		userId = userInfo.Tag
	}

	isAllow, err := s.isAllowSTS(r, userId, req.RoleYrn)
	if !isAllow {
		logging.Default().Infof("AssumeRoleFail, %s", err.Error())
		common.ErrorResp(c, http.StatusForbidden, "Forbidden", fmt.Sprintf("You are not allowed to assume this role: %s.", req.RoleYrn))
		return
	}

	cred, expireTime, err := s.genTmpAK(r, time.Second*time.Duration(req.DurationSeconds), req.RoleYrn)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := &iam_api.AssumeRoleResponse{
		Credentials: cred,
		ExpireTime:  expireTime,
	}
	common.SuccessResp(c, res)
}

func (s *service) isAllowSTS(r *dao.Role, userID, yrn string) (bool, error) {
	request := ladon.Request{
		Subject:  userID,
		Action:   StsAction,
		Resource: yrn,
	}

	var ps []ladon.DefaultPolicy
	ps = append(ps, r.TrustPolicy.Policy)
	return s.p.DoPoliciesAllow(&request, ps)
}

func (s *service) genTmpAK(r *dao.Role, d time.Duration, yrn string) (*iam_api.Credentials, time.Time, error) {
	maxTime := config.GetConfig().MaxAssumeRoleTime
	if maxTime <= 0 {
		maxTime = 4
	}
	if int64(d) > int64(time.Second*3600)*maxTime || d <= 0 {
		d = time.Second * 3600 * 4
	}
	claims := make(map[string]interface{})
	claims[roleYrnClaim] = yrn
	expTime := time.Now().Add(d)
	// jwt claims exp is numericDate
	claims["exp"] = expTime.Unix()
	claims["parent"] = r.UserId

	// 服务端已经保存了session token，在STS功能上其实可以不需要向客户端返回session token
	// 返回session token的原因：1. 可能是为了通信的可追溯 2. 凭证的部分可视化
	// TODO 将claims加密生成session token
	accessKey, secretKey, sessionToken, err := GetNewCredentialsWithMetadata(claims)
	// akId, akSecret, err := GenerateCredentials()
	if err != nil {
		logging.Default().Errorf("GenerateCredentialsFail, Error: %s", err.Error())
		return nil, time.Time{}, err
	}

	secret := &dao.Secret{
		AccessKeyId:     accessKey,
		AccessKeySecret: secretKey,
		ParentUser:      fmt.Sprintf("sts:%s:%s", r.UserId, r.RoleName),
		Expiration:      expTime,
		Claims:          claims,
		SessionToken:    sessionToken,
	}
	err = s.store.Secrets().Create(context.Background(), secret)
	if err != nil {
		return nil, time.Time{}, err
	}
	res := iam_api.Credentials{
		AccessKeyId:     secret.AccessKeyId,
		AccessKeySecret: secret.AccessKeySecret,
		SessionToken:    secret.SessionToken,
	}
	return &res, expTime, nil
}

// only yuansuan product account can access this
func (s *service) IsAllow(c *gin.Context) {
	req := iam_api.IsAllowRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		common.InvalidParams(c, "BindReqFail")
		return
	}
	if len(req.Subject) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, "InvalidArgument", "subject is required")
		return
	}
	if len(req.Action) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, "InvalidArgument", "action is required")
		return
	}
	if len(req.Resource) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, "InvalidArgument", "resource is required")
		return
	}
	userInfo := common.GetUserInfo(c)
	if !userInfo.IsAdmin() && !IsYuansuanProductAccount(userInfo.Tag) {
		common.ErrorRespWithAbort(c, http.StatusForbidden, "Forbidden", "Only yusansuan product account can call this")
		return
	}
	res := &iam_api.IsAllowResponse{}
	secret, err := s.store.Secrets().Get(c, req.Subject)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, "NotFound.Secret", "secret not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	if secret == nil {
		res.Message = fmt.Sprintf("AccessKeyId %s not found", req.Subject)
		common.SuccessResp(c, res)
		return
	}
	y, err := common.ParseYRN(req.Resource)
	if err != nil {
		common.ErrorResp(c, http.StatusBadRequest, "InvalidResourceName", err.Error())
		return
	}
	if y.AccountID == secret.ParentUser {
		res.Allow = true
		res.Message = "access to own resource."
		common.SuccessResp(c, res)
		return
	}
	if secret.SessionToken == "" {
		res.Message = "this accessKey is not STS, access denied."
		common.SuccessResp(c, res)
		return
	}
	if _, ok := secret.Claims[roleYrnClaim]; !ok {
		logging.Default().Errorf("Empty yrn in STS, RequestId: %s", common.GetRequestID(c))
		common.InternalServerError(c, "")
		return
	}
	roleYrn, err := common.ParseYRN(secret.Claims[roleYrnClaim].(string))
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	if roleYrn.AccountID != y.AccountID {
		res.Message = fmt.Sprintf("AccountID not same, %s, %s", roleYrn.AccountID, y.AccountID)
		common.SuccessResp(c, res)
		return
	}
	if time.Now().After(secret.Expiration) {
		res.Message = fmt.Sprintf("Time expirs: %+v", secret.Expiration)
		common.SuccessResp(c, res)
		return
	}
	r, err := s.store.Roles().Get(context.Background(), roleYrn.AccountID, roleYrn.ResourceID)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			res.Message = fmt.Sprintf("Role not found, RoleId: %s, UserId: %s", roleYrn.ResourceID, roleYrn.AccountID)
			common.SuccessResp(c, res)
			return
		}
		common.InternalServerError(c, "")
		return
	}
	pIds, err := s.store.RolePolicyRelations().ListPolicyByRoleId(c, r.ID, 0, 1000)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	if len(pIds) == 0 {
		res.Message = fmt.Sprintf("No Policy found, RoleId: %d", r.ID)
		common.SuccessResp(c, res)
		return
	}
	policys, err := s.store.Policies().GetByIds(context.Background(), pIds)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	// for check policy, set subject dynamically
	for _, p := range policys {
		subject := []string{roleYrn.ResourceID}
		p.Policy.Policy.Subjects = subject
	}
	if s.checkPolicy(roleYrn.ResourceID, &req, policys) {
		res.Allow = true
	} else {
		res.Message = fmt.Sprintf("No allow in policies, Policies: %+v", policys)
	}
	common.SuccessResp(c, res)
}

// checkPolicy
// common useage, storage use compute sts ak to check policy if allow
func (s *service) checkPolicy(sub string, req *iam_api.IsAllowRequest, ps []*dao.Policy) bool {
	r := ladon.Request{
		Resource: req.Resource,
		Action:   req.Action,
		Subject:  sub,
		Context:  ladon.Context{},
	}
	policies := lo.Map(ps, func(p *dao.Policy, index int) ladon.DefaultPolicy {
		return p.Policy.Policy
	})
	isAllow, err := s.p.DoPoliciesAllow(&r, policies)
	if err != nil {
		return false
	}
	return isAllow
}

func (s *service) tryAddPlatformRole(c *gin.Context, roleYrn string) (*common.YRN, error) {
	yrn, err := common.ParseYRN(roleYrn)
	if err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, "InvalidArgument", err.Error())
		return nil, errors.WithCode(code.ErrInvalidArgument, err.Error())
	}
	if !IsValidYrnForRole(yrn) {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, "InvalidArgument", "invalid role yrn")
		return nil, errors.WithCode(code.ErrInvalidArgument, "invalid role yrn")
	}

	roleName := yrn.ResourceID
	userID := yrn.AccountID

	// platform role is special, can't be created by user
	if IsYusanRole(roleName) {
		// GetQuota from DB use roleName and userID, may be yrn is a better choice
		roleInfo, err := s.store.Roles().Get(c, userID, roleName)
		if err != nil && !errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusInternalServerError, "InternalServerError", err.Error())
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
		if roleInfo != nil {
			// role exist, do nothing
			return yrn, nil
		}
		// role not exist, create it.
		// generally, platform role should be specified by user
		err = s.createPlatformRole(c, roleYrn, roleName, userID)
		if err != nil {
			logging.Default().Errorf("create platform role error: %v", err)
			common.ErrorRespWithAbort(c, http.StatusInternalServerError, "InternalServerError", err.Error())
			return nil, errors.WithCode(code.ErrDatabase, err.Error())
		}
	}
	return yrn, nil
}

func (s *service) createPlatformRole(c context.Context, roleYrn, roleName, userID string) error {
	tran, err := mysql.GetDB(s.store)
	if err != nil {
		return errors.WithCode(code.ErrBind, err.Error())
	}
	tx := tran.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	// set tx to context
	ctx := context.WithValue(c, common.ContextTransactionKey, tx)

	if err := s.addRoleForPlatform(ctx, roleYrn, roleName, userID); err != nil {
		logging.Default().Errorf("tran--add role for platform error: %v", err)
		tx.Rollback()
		return err
	}

	if err := s.loadManagedPolicies(ctx, roleName, userID); err != nil {
		logging.Default().Errorf("tran--load managed policies error: %v", err)
		tx.Rollback()
		return err
	}
	if err := s.attachPolicyToRole(ctx, roleName, userID); err != nil {
		logging.Default().Errorf("tran--attach policy to role error: %v", err)
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (s *service) addRoleForPlatform(c context.Context, roleYrn, roleName, userID string) error {
	policy := ladon.DefaultPolicy{
		ID:          "trust_policy",
		Description: fmt.Sprintf("This policy allows the %s perform %s the 'sts:AssumeRole' action", roleName, userID),
		Subjects:    []string{subjectName(roleName)},
		Resources:   []string{roleYrn},
		Actions:     []string{StsAction},
		Effect:      ladon.AllowAccess,
	}
	trustPolicy := dao.AuthzPolicy{
		Policy: policy,
	}

	role := &dao.Role{
		RoleName:    roleName,
		UserId:      userID,
		TrustPolicy: trustPolicy,
	}
	duplicateRole, err := s.store.Roles().Create(c, role)
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	if duplicateRole {
		logging.Default().Infof("role %s and userID %s already exist", roleName, userID)
	}
	return nil
}

func (s *service) loadManagedPolicies(c context.Context, roleName, userID string) error {
	managedPolicies := getYusnuanPolicy(roleName)
	for _, p := range managedPolicies {
		res := p.Policy.Policy.Resources[0]
		yrn, err := common.ParseYRN(res)
		if err != nil {
			logging.Default().Errorf("parse default yrn %s failed, err: %s", res, err.Error())
			continue
		}
		yrn.AccountID = userID
		p.Policy.Policy.Resources[0] = yrn.String()
		p.Policy.Policy.Subjects = []string{roleName}
		p.UserId = userID
	}
	// save to DB
	err := s.store.Policies().BatchCreate(c, managedPolicies)
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

func (s *service) attachPolicyToRole(c context.Context, roleName, userID string) error {
	// attach platform role policy to role
	managedPolicies := getYusnuanPolicy(roleName)
	policyNames := lo.Map(managedPolicies, func(p *dao.Policy, index int) string {
		return p.PolicyName
	})
	policies, err := s.store.Policies().ListByNameAndUserId(c, userID, policyNames)
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	policyIDs := lo.Map(policies, func(p *dao.Policy, index int) snowflake.ID {
		return p.ID
	})

	role, err := s.store.Roles().Get(c, userID, roleName)
	if err != nil {
		// rollback
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	var rolePolicyRelations []*dao.RolePolicyRelation
	for _, p := range policyIDs {
		rolePolicyRelations = append(rolePolicyRelations, &dao.RolePolicyRelation{
			RoleId:   role.ID,
			PolicyId: p,
		})
	}
	err = s.store.RolePolicyRelations().CreateBatch(c, rolePolicyRelations)
	if err != nil {
		return errors.WithCode(code.ErrDatabase, err.Error())
	}
	return nil
}

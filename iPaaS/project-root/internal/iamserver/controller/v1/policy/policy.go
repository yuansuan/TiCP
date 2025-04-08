package policy

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	v1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/service/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
)

type PolicyController struct {
	srv srvv1.Svc
}

func NewPolicyController(s store.Factory) *PolicyController {
	return &PolicyController{
		srv: srvv1.NewSvc(s),
	}
}

func (p *PolicyController) GetPolicy(c *gin.Context) {
	policyName := strings.TrimSpace(c.Param("policyName"))
	if len(policyName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid policy name")
		return
	}
	userInfo := common.GetUserInfo(c)
	policy, err := p.srv.Policies().GetByPolicyName(c, userInfo.UserID.String(), policyName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.PolicyNotFound, "policy not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	res := &iam_api.GetPolicyResponse{
		PolicyName: policy.PolicyName,
		Effect:     policy.Policy.Policy.Effect,
		Actions:    policy.Policy.Policy.Actions,
		Resources:  policy.Policy.Policy.Resources,
	}
	common.SuccessResp(c, res)
}

func (p *PolicyController) AddPolicy(c *gin.Context) {
	req := iam_api.AddPolicyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	if !isValidAddPolicy(c, &req) {
		return
	}
	userInfo := common.GetUserInfo(c)
	err := p.srv.Policies().AddPolicy(c, userInfo.UserID.String(), &req)
	policyErrHandler(c, err)
}

func isValidAddPolicy(c *gin.Context, req *iam_api.AddPolicyRequest) bool {
	if len(req.PolicyName) == 0 || len(req.PolicyName) > 64 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid policy name")
		return false
	}
	if req.Effect != ladon.AllowAccess && req.Effect != ladon.DenyAccess {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid effect")
		return false
	}
	if len(req.Actions) == 0 || len(req.Resources) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid action or resource")
		return false
	}
	for _, a := range req.Actions {
		if len(a) == 0 || strings.TrimSpace(a) == "*" {
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid action, wildcard action should be <.*>")
			return false
		}
	}
	for _, r := range req.Resources {
		y, err := common.ParseYRN(r)
		if err != nil {
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid yrn")
			return false
		}
		if y.ResourceID == "*" {
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid resource, wildcard resource should be like <yrn:ys:iam:::role/<.*>")
			return false
		}
	}
	return true
}

func (p *PolicyController) DeletePolicy(c *gin.Context) {
	policyName := strings.TrimSpace(c.Param("policyName"))
	if len(policyName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid policy name")
		return
	}
	userInfo := common.GetUserInfo(c)

	err := p.srv.Policies().DeletePolicy(c, userInfo.UserID.String(), policyName)
	policyErrHandler(c, err)
}

func (p *PolicyController) ListPolicy(c *gin.Context) {
	userInfo := common.GetUserInfo(c)
	policies, err := p.srv.Policies().ListPolicy(c, userInfo.UserID.String(), 0, 1000)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	ps := make([]*iam_api.GetPolicyResponse, 0, len(policies))
	for _, policy := range policies {
		ps = append(ps, &iam_api.GetPolicyResponse{
			PolicyName: policy.PolicyName,
			Effect:     policy.Policy.Policy.Effect,
			Actions:    policy.Policy.Policy.Actions,
			Resources:  policy.Policy.Policy.Resources,
		})
	}
	res := &iam_api.ListPolicyResponse{
		Policies: ps,
	}

	common.SuccessResp(c, res)
}

func (p *PolicyController) UpdatePolicy(c *gin.Context) {
	req := iam_api.UpdatePolicyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	add := req.Policy
	if !isValidAddPolicy(c, &add) {
		return
	}
	userInfo := common.GetUserInfo(c)
	err := p.srv.Policies().UpdatePolicy(c, userInfo.UserID.String(), &req.Policy)
	policyErrHandler(c, err)
}

// admin api

func (p *PolicyController) AdminAddPolicy(c *gin.Context) {
	req := iam_api.AdminAddPolicyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}

	if len(req.UserId) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid user id")
		return
	}

	add := req.AddPolicyRequest
	if !isValidAddPolicy(c, &add) {
		return
	}
	err := p.srv.Policies().AddPolicy(c, req.UserId, &req.AddPolicyRequest)
	policyErrHandler(c, err)
}

func (p *PolicyController) AdminDeletePolicy(c *gin.Context) {
	policyName := strings.TrimSpace(c.Param("policyName"))
	userId := strings.TrimSpace(c.Param("userId"))
	if len(policyName) == 0 || len(userId) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid policy name or user id")
		return
	}

	err := p.srv.Policies().DeletePolicy(c, userId, policyName)
	policyErrHandler(c, err)
}

func (p *PolicyController) AdminListPolicy(c *gin.Context) {
	userId := strings.TrimSpace(c.Param("userId"))
	if len(userId) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid user id")
		return
	}
	policies, err := p.srv.Policies().ListPolicy(c, userId, 0, 1000)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	ps := make([]*iam_api.GetPolicyResponse, 0, len(policies))
	for _, policy := range policies {
		ps = append(ps, &iam_api.GetPolicyResponse{
			PolicyName: policy.PolicyName,
			Effect:     policy.Policy.Policy.Effect,
			Actions:    policy.Policy.Policy.Actions,
			Resources:  policy.Policy.Policy.Resources,
		})
	}
	res := &iam_api.AdminListPolicyResponse{
		Policies: ps,
	}
	common.SuccessResp(c, res)
}

func (p *PolicyController) AdminGetPolicy(c *gin.Context) {
	policyName := strings.TrimSpace(c.Param("policyName"))
	userId := strings.TrimSpace(c.Param("userId"))
	if len(policyName) == 0 || len(userId) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid policy name or user id")
		return
	}

	policy, err := p.srv.Policies().GetByPolicyName(c, userId, policyName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.PolicyNotFound, "not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	res := &iam_api.GetPolicyResponse{
		PolicyName: policy.PolicyName,
		Effect:     policy.Policy.Policy.Effect,
		Actions:    policy.Policy.Policy.Actions,
		Resources:  policy.Policy.Policy.Resources,
	}
	common.SuccessResp(c, res)
}

func (p *PolicyController) AdminUpdatePolicy(c *gin.Context) {
	userID := strings.TrimSpace(c.Param("userId"))
	policyName := strings.TrimSpace(c.Param("policyName"))
	if len(userID) == 0 || len(policyName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid user id or policy name")
		return
	}
	req := iam_api.AdminUpdatePolicyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}

	if !isValidAddPolicy(c, &req.Policy) {
		return
	}
	err := p.srv.Policies().UpdatePolicy(c, req.UserId, &req.Policy)
	policyErrHandler(c, err)
}

func policyErrHandler(c *gin.Context, err error) {
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.PolicyNotFound, "policy not found")
			return
		}
		if errors.Is(err, common.ErrAlreadyExists) {
			common.ErrorRespWithAbort(c, http.StatusConflict, v1.AlreadyExists, "policy already exists")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	common.SuccessResp(c, nil)
}

func (p *PolicyController) InvalidPolicyName(c *gin.Context) {
	common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "empty policy name")
	return
}

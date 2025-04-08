package role

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ory/ladon"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	v1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/service/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"

	"net/http"

	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

type RoleController struct {
	srv srvv1.Svc
}

func NewRoleController(s store.Factory) *RoleController {
	return &RoleController{
		srv: srvv1.NewSvc(s),
	}
}

func (r *RoleController) GetRole(c *gin.Context) {
	roleName := strings.TrimSpace(c.Param("roleName"))
	if len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
		return
	}
	caller := common.GetUserInfo(c)
	role, err := r.srv.Roles().GetRole(c, caller.UserID.String(), roleName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RoleNotFound, "not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}

	res := &iam_api.GetRoleResponse{
		RoleName:    role.RoleName,
		Description: role.Description,
		TrustPolicy: &iam_api.RolePolicy{
			Actions:    role.TrustPolicy.Policy.Actions,
			Resources:  role.TrustPolicy.Policy.Resources,
			Effect:     role.TrustPolicy.Policy.Effect,
			Principals: role.TrustPolicy.Policy.Subjects,
		},
	}

	common.SuccessResp(c, res)
}

func (r *RoleController) ListRole(c *gin.Context) {
	caller := common.GetUserInfo(c)
	roles, err := r.srv.Roles().ListRole(c, caller.UserID.String())
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := make([]*iam_api.Role, 0, len(roles))
	for _, role := range roles {
		policy := &iam_api.RolePolicy{
			Actions:    role.TrustPolicy.Policy.Actions,
			Resources:  role.TrustPolicy.Policy.Resources,
			Effect:     role.TrustPolicy.Policy.Effect,
			Principals: role.TrustPolicy.Policy.Subjects,
		}
		res = append(res, &iam_api.Role{
			RoleName:    role.RoleName,
			Description: role.Description,
			TrustPolicy: policy,
		})
	}
	list := &iam_api.ListRoleResponse{
		Roles: res,
	}
	common.SuccessResp(c, list)
}

func (r *RoleController) DeleteRole(c *gin.Context) {
	roleName := strings.TrimSpace(c.Param("roleName"))
	if len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
		return
	}
	caller := common.GetUserInfo(c)
	err := r.srv.Roles().Delete(c, caller.UserID.String(), roleName)
	roleErrHandler(c, err)
}

func (r *RoleController) AddRole(c *gin.Context) {
	var req iam_api.AddRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	if !isValidAddRoleReq(c, req) {
		return
	}
	if len(req.TrustPolicy.Resources) != 0 {
		res := req.TrustPolicy.Resources
		for _, yrn := range res {
			if _, err := common.ParseYRN(yrn); err != nil {
				logging.Default().Infof("invalid yrn: %s", yrn)
				common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
				return
			}
		}
	}
	userInfo := common.GetUserInfo(c)
	err := r.srv.Roles().AddRole(c, userInfo.UserID.String(), &req)
	roleErrHandler(c, err)
}

func isValidAddRoleReq(c *gin.Context, req iam_api.AddRoleRequest) bool {
	if !srvv1.IsValidRoleName(req.RoleName) {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
		return false
	}
	if req.TrustPolicy.Effect != ladon.AllowAccess && req.TrustPolicy.Effect != ladon.DenyAccess {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid effect")
		return false
	}
	if req.TrustPolicy.Principals == nil || len(req.TrustPolicy.Principals) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid principals")
		return false
	}

	if len(req.Description) > 255 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "description too long")
		return false
	}
	if req.TrustPolicy == nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid trust policy")
		return false
	}
	if len(req.TrustPolicy.Actions) == 0 || req.TrustPolicy.Actions[0] != srvv1.StsAction {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "only support sts action")
		return false
	}
	return true
}

func (r *RoleController) UpdateRole(c *gin.Context) {
	var req iam_api.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	if !isValidUpdateRole(c, req.Role) {
		return
	}
	yrns := req.Role.TrustPolicy.Resources
	if len(yrns) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid yrn")
		return
	}
	for _, yrn := range yrns {
		if _, err := common.ParseYRN(yrn); err != nil {
			logging.Default().Infof("invalid yrn: %s", yrn)
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid yrn")
			return
		}
	}
	userInfo := common.GetUserInfo(c)
	err := r.srv.Roles().UpdateRole(c, userInfo.UserID.String(), req.Role)
	roleErrHandler(c, err)
}

func isValidUpdateRole(c *gin.Context, req *iam_api.Role) bool {
	if !srvv1.IsValidRoleName(req.RoleName) {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
		return false
	}
	if req.TrustPolicy.Effect != ladon.AllowAccess && req.TrustPolicy.Effect != ladon.DenyAccess {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid effect")
		return false
	}
	if req.TrustPolicy.Principals == nil || len(req.TrustPolicy.Principals) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid principals")
		return false
	}

	if len(req.Description) > 255 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "description too long")
		return false
	}
	if req.TrustPolicy == nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid trust policy")
		return false
	}
	if len(req.TrustPolicy.Actions) == 0 || req.TrustPolicy.Actions[0] != srvv1.StsAction {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "only support sts action")
		return false
	}
	return true
}

func (r *RoleController) PatchPolicy(c *gin.Context) {
	roleName := strings.TrimSpace(c.Param("roleName"))
	policyName := strings.TrimSpace(c.Param("policyName"))
	if len(roleName) == 0 || len(policyName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName or policyName")
		return
	}
	userInfo := common.GetUserInfo(c)
	userID := userInfo.UserID.String()
	r.patchPolicy(c, roleName, policyName, userID)
}

func (r *RoleController) DetachPolicy(c *gin.Context) {
	roleName := strings.TrimSpace(c.Param("roleName"))
	policyName := strings.TrimSpace(c.Param("policyName"))
	if len(roleName) == 0 || len(policyName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName or policyName")
		return
	}
	userInfo := common.GetUserInfo(c)
	userID := userInfo.UserID.String()
	r.detachPolicy(c, roleName, policyName, userID)
}

func (r *RoleController) detachPolicy(c *gin.Context, roleName, policyName, userID string) {
	role, err := r.srv.Roles().GetRole(c, userID, roleName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RoleNotFound, "role not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	p, err := r.srv.Policies().GetByPolicyName(c, userID, policyName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.PolicyNotFound, "policy not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	err = r.srv.Roles().DetachPolicy(c, role, p)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RelationNotFound, "relation not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	common.SuccessResp(c, nil)
}

func (r *RoleController) patchPolicy(c *gin.Context, roleName, policyName, userID string) {
	role, err := r.srv.Roles().GetRole(c, userID, roleName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RoleNotFound, "role not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	p, err := r.srv.Policies().GetByPolicyName(c, userID, policyName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.PolicyNotFound, "policy not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	err = r.srv.Roles().PatchPolicy(c, role, p)
	if err != nil {
		if errors.Is(err, common.ErrAlreadyExists) {
			common.ErrorRespWithAbort(c, http.StatusConflict, v1.AlreadyExists, "relation already exists")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	common.SuccessResp(c, nil)
}

// admin api
func (r *RoleController) AdminAddRole(c *gin.Context) {
	var req iam_api.AdminAddRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	if len(req.UserId) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid userId")
		return
	}
	addReq := req.AddRoleRequest
	if !isValidAddRoleReq(c, addReq) {
		return
	}
	if len(req.TrustPolicy.Resources) != 0 {
		res := req.TrustPolicy.Resources
		for _, yrn := range res {
			if _, err := common.ParseYRN(yrn); err != nil {
				logging.Default().Infof("invalid yrn: %s", yrn)
				common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid yrn")
				return
			}
		}
	}
	err := r.srv.Roles().AddRole(c, req.UserId, &req.AddRoleRequest)
	roleErrHandler(c, err)
}

func (r *RoleController) AdminGetRole(c *gin.Context) {
	userID := c.Param("userId")
	roleName := c.Param("roleName")
	if len(userID) == 0 || len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	role, err := r.srv.Roles().GetRole(c, userID, roleName)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RoleNotFound, "role not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	if role == nil {
		common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RoleNotFound, "role not found")
		return
	}
	res := &iam_api.AdminGetRoleResponse{
		Role: &iam_api.Role{
			RoleName:    role.RoleName,
			Description: role.Description,
			TrustPolicy: &iam_api.RolePolicy{
				Actions:    role.TrustPolicy.Policy.Actions,
				Resources:  role.TrustPolicy.Policy.Resources,
				Effect:     role.TrustPolicy.Policy.Effect,
				Principals: role.TrustPolicy.Policy.Subjects,
			},
		},
	}
	common.SuccessResp(c, res)
}

func (r *RoleController) AdminListRole(c *gin.Context) {
	userID := strings.TrimSpace(c.Param("userId"))
	if len(userID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid use id")
		return
	}
	roles, err := r.srv.Roles().ListRole(c, userID)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := make([]*iam_api.Role, 0, len(roles))
	for _, role := range roles {
		policy := &iam_api.RolePolicy{
			Actions:    role.TrustPolicy.Policy.Actions,
			Resources:  role.TrustPolicy.Policy.Resources,
			Effect:     role.TrustPolicy.Policy.Effect,
			Principals: role.TrustPolicy.Policy.Subjects,
		}
		res = append(res, &iam_api.Role{
			RoleName:    role.RoleName,
			Description: role.Description,
			TrustPolicy: policy,
		})
	}
	list := &iam_api.AdminListRoleResponse{
		Roles: res,
	}
	common.SuccessResp(c, list)
}

func (r *RoleController) AdminDeleteRole(c *gin.Context) {
	userID := strings.TrimSpace(c.Param("userId"))
	roleName := strings.TrimSpace(c.Param("roleName"))
	if len(userID) == 0 || len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	err := r.srv.Roles().Delete(c, userID, roleName)
	roleErrHandler(c, err)
}

func (r *RoleController) AdminUpdateRole(c *gin.Context) {
	var req iam_api.AdminUpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	userID := c.Param("userId")
	roleName := c.Param("roleName")
	if len(userID) == 0 || len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid userId or roleName")
		return
	}

	if !isValidUpdateRole(c, req.Role) {
		return
	}

	yrns := req.Role.TrustPolicy.Resources
	if len(yrns) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid yrn")
		return
	}
	for _, yrn := range yrns {
		if _, err := common.ParseYRN(yrn); err != nil {
			logging.Default().Infof("invalid yrn: %s", yrn)
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid yrn")
			return
		}
	}
	err := r.srv.Roles().UpdateRole(c, userID, req.Role)
	roleErrHandler(c, err)
}

func (r *RoleController) AdminPatchPolicy(c *gin.Context) {
	var req iam_api.AdminPatchPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	userID := req.UserId
	if len(userID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid userId")
		return
	}
	roleName := strings.TrimSpace(c.Param("roleName"))
	if len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
		return
	}
	policyName := req.PolicyName
	if len(policyName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid policyName")
		return
	}
	r.patchPolicy(c, roleName, policyName, userID)
}

func (r *RoleController) AdminDetachPolicy(c *gin.Context) {
	var req iam_api.AdminDetachPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	userID := req.UserId
	roleName := strings.TrimSpace(c.Param("roleName"))
	if len(roleName) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
		return
	}
	policyName := req.PolicyName
	r.detachPolicy(c, roleName, policyName, userID)
}

func roleErrHandler(c *gin.Context, err error) {
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.RoleNotFound, "role not found")
			return
		}
		if errors.Is(err, common.ErrAlreadyExists) {
			common.ErrorRespWithAbort(c, http.StatusConflict, v1.AlreadyExists, "role already exists")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	common.SuccessResp(c, nil)
}

func (s *RoleController) InvalidRoleName(c *gin.Context) {
	common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid roleName")
}

func (s *RoleController) InvalidPolicyName(c *gin.Context) {
	common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "empty policyName")
}

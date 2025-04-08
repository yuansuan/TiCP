package iam_client

import (
	"fmt"
	"net/http"
	"strings"

	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
)

type IamAdminClient struct {
	IamClient
}

func NewAdminClient(endpoint, akId, akSecret, token string) *IamAdminClient {
	return &IamAdminClient{
		IamClient{
			endpoint: strings.TrimRight(endpoint, "/"),
			akId:     akId,
			akSecret: akSecret,
			token:    token,
		},
	}
}

func (c *IamAdminClient) AddSecret(req *iam_api.AdminAddSecretRequest) (*iam_api.AdminAddSecretResponse, error) {
	baseUrl := "/iam/admin/secrets"
	res := &iam_api.AdminAddSecretResponse{}
	err := c.SendWithLog("POST", baseUrl, req, res, "AdminAddSecret")
	return res, err
}

func (c *IamAdminClient) GetSecret(req *iam_api.AdminGetSecretRequest) (*iam_api.AdminGetSecretResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/secrets/%s", req.AccessKeyId)
	res := &iam_api.AdminGetSecretResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "AdminGetSecret")
	return res, err
}

func (c *IamAdminClient) DeleteSecret(req *iam_api.AdminDeleteSecretRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/secrets/%s", req.AccessKeyId)
	err := c.SendWithLog("DELETE", baseUrl, nil, nil, "AdminDeleteSecret")
	return err
}

func (c *IamAdminClient) ListSecret(req *iam_api.AdminListSecretRequest) (*iam_api.AdminListSecretResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/secrets/user/%s", req.UserId)
	res := &iam_api.AdminListSecretResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "AdminListSecret")
	return res, err
}

func (c *IamAdminClient) UpdateTag(req *iam_api.AdminUpdateTagRequest) error {
	// baseUrl := "/iam/admin/secrets"
	baseUrl := fmt.Sprintf("/iam/admin/secrets/%s/tag", req.AccessKeyId)
	err := c.SendWithLog(http.MethodPut, baseUrl, req, nil, "AdminUpdateTag")
	return err
}

func (c *IamAdminClient) ListSecrets(req *iam_api.AdminListSecretsRequest) (*iam_api.AdminListSecretResponse, error) {
	baseUrl := "/iam/admin/secrets"
	res := &iam_api.AdminListSecretResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "AdminListSecrets", fmt.Sprintf("PageOffset=%d&PageSize=%d", req.PageOffset, req.PageSize))
	return res, err
}

func (c *IamAdminClient) AddRole(req *iam_api.AdminAddRoleRequest) error {
	baseUrl := "/iam/admin/roles"
	err := c.SendWithLog("POST", baseUrl, req, nil, "AdminAddRole")
	return err
}

func (c *IamAdminClient) GetRole(req *iam_api.AdminGetRoleRequest) (*iam_api.AdminGetRoleResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/roles/%s/%s", req.UserId, req.RoleName)
	res := &iam_api.AdminGetRoleResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "AdminGetRole")
	return res, err
}

func (c *IamAdminClient) ListRole(req *iam_api.AdminListRoleRequest) (*iam_api.AdminListRoleResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/roles/%s", req.UserId)
	res := &iam_api.AdminListRoleResponse{}
	err := c.SendWithLog(http.MethodGet, baseUrl, nil, res, "AdminListRole")
	return res, err
}

func (c *IamAdminClient) DeleteRole(req *iam_api.AdminDeleteRoleRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/roles/%s/%s", req.UserId, req.RoleName)
	err := c.SendWithLog(http.MethodDelete, baseUrl, nil, nil, "AdminDeleteRole")
	return err
}
func (c *IamAdminClient) UpdateRole(req *iam_api.AdminUpdateRoleRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/roles/%s/%s", req.UserId, req.Role.RoleName)
	err := c.SendWithLog(http.MethodPut, baseUrl, req, nil, "AdminUpdateRole")
	return err
}

func (c *IamAdminClient) AddPolicy(req *iam_api.AdminAddPolicyRequest) error {
	baseUrl := "/iam/admin/policies"
	err := c.SendWithLog("POST", baseUrl, req, nil, "AdminAddPolicy")
	return err
}

func (c *IamAdminClient) GetPolicy(req *iam_api.AdminGetPolicyRequest) (*iam_api.AdminGetPolicyResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/policies/%s/%s", req.UserId, req.PolicyName)
	res := &iam_api.AdminGetPolicyResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "AdminGetPolicy")
	return res, err
}

func (c *IamAdminClient) ListPolicy(req *iam_api.AdminListPolicyRequest) (*iam_api.AdminListPolicyResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/policies/%s", req.UserId)
	res := &iam_api.AdminListPolicyResponse{}
	err := c.SendWithLog(http.MethodGet, baseUrl, nil, res, "AdminListPolicy")
	return res, err
}

func (c *IamAdminClient) UpdatePolicy(req *iam_api.AdminUpdatePolicyRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/policies/%s/%s", req.UserId, req.Policy.PolicyName)
	err := c.SendWithLog(http.MethodPut, baseUrl, req, nil, "AdminUpdatePolicy")
	return err
}

func (c *IamAdminClient) DeletePolicy(req *iam_api.AdminDeletePolicyRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/policies/%s/%s", req.UserId, req.PolicyName)
	err := c.SendWithLog(http.MethodDelete, baseUrl, nil, nil, "AdminDeletePolicy")
	return err
}

func (c *IamAdminClient) AddRolePolicyRelation(req *iam_api.AdminAddRolePolicyRelationRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/roles/%s", req.RoleName)
	err := c.SendWithLog("PATCH", baseUrl, req, nil, "AdminAddRolePolicyRelation")
	return err
}

func (c *IamAdminClient) DeleteRolePolicyRelation(req *iam_api.AdminDeleteRolePolicyRelationRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/roles/%s", req.RoleName)
	err := c.SendWithLog(http.MethodPost, baseUrl, req, nil, "AdminDeleteRolePolicyRelation")
	return err
}

func (c *IamAdminClient) GetUserInfo(req *iam_api.AdminGetUserRequest) (*iam_api.AdminGetUserResponse, error) {
	baseUrl := fmt.Sprintf("/iam/admin/users/%s", req.UserId)
	res := &iam_api.AdminGetUserResponse{}
	err := c.SendWithLog(http.MethodGet, baseUrl, nil, res, "AdminGetUser")
	return res, err
}

func (c *IamAdminClient) AddUser(req *iam_api.AdminAddUserRequest) (*iam_api.AdminAddUserResponse, error) {
	baseUrl := "/iam/admin/users"
	res := &iam_api.AdminAddUserResponse{}
	err := c.SendWithLog(http.MethodPost, baseUrl, req, res, "AdminAddUser")
	return res, err
}

func (c *IamAdminClient) UpdateUser(req *iam_api.AdminUpdateUserRequest) error {
	baseUrl := fmt.Sprintf("/iam/admin/users/%s", req.UserId)
	err := c.SendWithLog(http.MethodPut, baseUrl, req, nil, "AdminUpdateUser")
	return err
}

func (c *IamAdminClient) ListUserByName(req *iam_api.AdminListUserByNameRequest) (*iam_api.AdminListUserByNameResponse, error) {
	baseUrl := "/iam/admin/users"
	res := &iam_api.AdminListUserByNameResponse{}
	err := c.SendWithLog(http.MethodGet, baseUrl, nil, res, "AdminListUserByName", fmt.Sprintf("PageOffset=%d&PageSize=%d&Name=%s", req.PageOffset, req.PageSize, req.Name))
	return res, err
}

package iam_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils/signer"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/common/project-root-iam/iam-client/cache"
)

var (
	once      sync.Once
	iamClient *IamClient
)

type IamClient struct {
	endpoint string
	akId     string
	akSecret string
	token    string
	Proxy    string
	cache    *cache.Cache
	useCache bool
	initDone bool
}

type Option func(*IamClient)

func NewNonSingletonClient(endpoint, akId, akSecert string, options ...Option) *IamClient {
	iamClient = &IamClient{
		endpoint: strings.TrimRight(endpoint, "/"),
		akId:     akId,
		akSecret: akSecert,
		useCache: false,
	}
	for _, option := range options {
		option(iamClient)
	}
	return iamClient
}

func NewClient(endpoint, akId, akSecert string, options ...Option) *IamClient {

	once.Do(func() {
		iamClient = &IamClient{
			endpoint: strings.TrimRight(endpoint, "/"),
			akId:     akId,
			akSecret: akSecert,
			useCache: true,
		}
		for _, option := range options {
			option(iamClient)
		}
		if iamClient.useCache {
			iamClient.cache = cache.NewCache()
			go func() {
				for range time.Tick(time.Second * 5) {
					logging.Default().Debugf("Reload iam cache,iam client pointer address is %p", iamClient)
					secrets, err := iamClient.listAllSecrets()
					if err != nil {
						if res, ok := err.(ErrorResponse); ok {
							switch res.Status {
							case http.StatusForbidden:
								logging.Default().Errorf("current user %s reload iam cache forbidden", iamClient.akId)
								return
							default:
								logging.Default().Errorf("Reload iam cache error: %v", err)
								continue
							}
						} else {
							logging.Default().Warnf("connection lost, reload iam cache error: %v", err)
							continue
						}
					}
					logging.Default().Debug("secrets' length: ", len(secrets.Secrets))
					iamClient.setSecrets(secrets)
					iamClient.initDone = true
				}
			}()
		}
	})
	return iamClient
}

func UseCache(isUse bool) Option {
	return func(c *IamClient) {
		c.useCache = isUse
	}
}

func Token(token string) Option {
	return func(c *IamClient) {
		c.token = token
	}
}

func (c *IamClient) setSecrets(secrets *iam_api.ListAllSecretResponse) {
	c.cache.ClearAndSet(secrets.Secrets)
}

func (c *IamClient) SetProxy(proxy string) {
	c.Proxy = proxy
}

func (c *IamClient) AssumeRole(req *iam_api.AssumeRoleRequest) (*iam_api.AssumeRoleResponse, error) {
	baseUrl := "/iam/v1/AssumeRole"
	res := &iam_api.AssumeRoleResponse{}
	err := c.SendWithLog("POST", baseUrl, req, res, "AssumeRole")
	return res, err
}

func (c *IamClient) AssumeRoleDefault(userId string, roleName string) (*iam_api.AssumeRoleResponse, error) {
	if roleName == "" {
		roleName = "YS_CloudComputeRole"
	}
	req := &iam_api.AssumeRoleRequest{
		RoleYrn:         fmt.Sprintf("yrn:ys:iam::%s:role/%s", userId, roleName),
		RoleSessionName: "noname",
		DurationSeconds: 3600 * 4,
	}
	return c.AssumeRole(req)
}

func (c *IamClient) IsAllow(req *iam_api.IsAllowRequest) (*iam_api.IsAllowResponse, error) {
	baseUrl := "/iam/v1/IsAllow"
	res := &iam_api.IsAllowResponse{}
	err := c.SendWithLog("POST", baseUrl, req, res, "IsAllow")
	return res, err
}

// ysProductName can be: cc(cloudcompute), cs(cloudstorage), csp, cae365 . and more in future
func (c *IamClient) IsAllowDefault(accessKeyId string, path string, ysProductName string) (*iam_api.IsAllowResponse, error) {
	if ysProductName == "" { // cloud storage
		ysProductName = "cs"
	}
	path = strings.TrimLeft(path, "/")
	l := strings.SplitN(path, "/", 2)
	if len(l) != 2 {
		return nil, ErrorResponse{
			Message: "Error path, like: /userid/path1/",
			Code:    "InvalidPath",
			Status:  http.StatusBadRequest,
		}
	}
	userId, resourceId := l[0], l[1]
	// root path, resourceId is empty,replaced by "/"
	if len(resourceId) == 0 {
		resourceId = "/"
	}
	req := &iam_api.IsAllowRequest{
		Action:   "*",
		Resource: fmt.Sprintf("yrn:ys:cs::%s:path/%s", userId, resourceId),
		Subject:  accessKeyId,
	}
	return c.IsAllow(req)
}

func (c *IamClient) GetSecret(req *iam_api.GetSecretRequest) (*iam_api.GetSecretResponse, error) {
	if c.useCache {
		if c.initDone {
			ak := req.AccessKeyId
			if ak == "" {
				return nil, errors.New("AccessKeyId is empty")
			}
			secret, err := c.cache.GetSecret(ak)
			if errors.Is(err, cache.ErrSecretNotFound) {
				// if not found in cache, get from server
				logging.Default().Infof("AccessKeyId %s not found in cache, get from server", ak)
				return c.getSecret(req)
			}
			res := &iam_api.GetSecretResponse{
				AccessKeyId:     secret.AccessKeyId,
				AccessKeySecret: secret.AccessKeySecret,
				YSId:            secret.YSId,
				Expire:          secret.Expire,
			}
			return res, nil
		}
	}
	return c.getSecret(req)
}

func (c *IamClient) getSecret(req *iam_api.GetSecretRequest) (*iam_api.GetSecretResponse, error) {
	baseUrl := fmt.Sprintf("/iam/v1/secrets/%s", req.AccessKeyId)
	res := &iam_api.GetSecretResponse{}
	err := c.SendWithLog(http.MethodGet, baseUrl, nil, res, "GetSecret")
	return res, err
}

func (c *IamClient) IsYSProductAccount(req *iam_api.IsYSProductAccountRequest) (*iam_api.IsYSProductAccountResponse, error) {
	baseUrl := "/iam/v1/IsYSProductAccount"
	res := &iam_api.IsYSProductAccountResponse{}
	queryParam := "UserId=" + req.UserId
	err := c.SendWithLog(http.MethodGet, baseUrl, req, res, "IsYSProductAccount", queryParam)
	return res, err
}

func (c *IamClient) AddSecret() (*iam_api.AddSecretResponse, error) {
	baseUrl := "/iam/v1/secrets"
	res := &iam_api.AddSecretResponse{}
	err := c.SendWithLog("POST", baseUrl, nil, res, "AddSecret")
	return res, err
}

func (c *IamClient) ListSecrets() (*iam_api.ListSecretResponse, error) {
	baseUrl := "/iam/v1/secrets"
	res := &iam_api.ListSecretResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "ListSecrets")
	return res, err
}

func (c *IamClient) DeleteSecret(req *iam_api.DeleteSecretRequest) (*iam_api.DeleteSecretResponse, error) {
	baseUrl := fmt.Sprintf("/iam/v1/secrets/%s", req.AccessKeyId)
	res := &iam_api.DeleteSecretResponse{}
	err := c.SendWithLog(http.MethodDelete, baseUrl, nil, res, "DeleteSecret")
	return res, err
}

func (c *IamClient) AddPolicy(req *iam_api.AddPolicyRequest) error {
	baseUrl := "/iam/v1/policies"
	return c.SendWithLog("POST", baseUrl, req, nil, "AddPolicy")
}

func (c *IamClient) GetPolicy(req *iam_api.GetPolicyRequest) (*iam_api.GetPolicyResponse, error) {
	baseUrl := "/iam/v1/policies/" + req.PolicyName
	res := &iam_api.GetPolicyResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "GetPolicy")
	return res, err
}

func (c *IamClient) ListPolicies() (*iam_api.ListPolicyResponse, error) {
	baseUrl := "/iam/v1/policies"
	res := &iam_api.ListPolicyResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "ListPolicies")
	return res, err
}

func (c *IamClient) UpdatePolicy(req *iam_api.UpdatePolicyRequest) error {
	baseUrl := "/iam/v1/policies/" + req.Policy.PolicyName
	return c.SendWithLog("PUT", baseUrl, req, nil, "UpdatePolicy")
}

func (c *IamClient) DeletePolicy(req *iam_api.DeletePolicyRequest) error {
	baseUrl := "/iam/v1/policies/" + req.PolicyName
	err := c.SendWithLog("DELETE", baseUrl, nil, nil, "DeletePolicy")
	return err
}

func (c *IamClient) GetRole(req *iam_api.GetRoleRequest) (*iam_api.GetRoleResponse, error) {
	baseUrl := "/iam/v1/roles/" + req.RoleName
	res := &iam_api.GetRoleResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "GetRole")
	return res, err
}

func (c *IamClient) ListRoles() (*iam_api.ListRoleResponse, error) {
	baseUrl := "/iam/v1/roles"
	res := &iam_api.ListRoleResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "ListRoles")
	return res, err
}

func (c *IamClient) AddRole(req *iam_api.AddRoleRequest) error {
	baseUrl := "/iam/v1/roles"
	return c.SendWithLog("POST", baseUrl, req, nil, "AddRole")
}

func (c *IamClient) UpdateRole(req *iam_api.UpdateRoleRequest) error {
	baseUrl := "/iam/v1/roles/" + req.Role.RoleName
	return c.SendWithLog(http.MethodPut, baseUrl, req, nil, "UpdateRole")
}

func (c *IamClient) DeleteRole(req *iam_api.DeleteRoleRequest) error {
	baseUrl := "/iam/v1/roles/" + req.RoleName
	return c.SendWithLog("DELETE", baseUrl, nil, nil, "DeleteRole")
}

func (c *IamClient) AddRelation(req *iam_api.AddRolePolicyRelationRequest) error {
	// baseUrl := "/iam/v1/roles/" + req.RoleName
	baseUrl := fmt.Sprintf("/iam/v1/roles/%s/policies/%s", req.RoleName, req.PolicyName)
	return c.SendWithLog("PATCH", baseUrl, req, nil, "AddRelation")
}

func (c *IamClient) DeleteRelation(req *iam_api.DeleteRolePolicyRelationRequest) error {
	baseUrl := fmt.Sprintf("/iam/v1/roles/%s/policies/%s", req.RoleName, req.PolicyName)
	return c.SendWithLog("DELETE", baseUrl, nil, nil, "DeleteRelation")
}

func (c *IamClient) SendWithLog(method, baseUrl string, inputData interface{}, resultData interface{},
	funcName string, queryParam ...string) error {
	baseRes, statusCode, err := c.Send(method, baseUrl, inputData, resultData, queryParam...)
	if err != nil || statusCode <= 0 {
		logging.Default().Warnf("%s() Fail, Error: %s, Status: %d, Response: %v", funcName,
			err.Error(), statusCode, baseRes)
	} else {
		logging.Default().Debugf("%sEnd, Status: %d, Response: %v", funcName,
			statusCode, baseRes)
	}
	return err
}

func (c *IamClient) listAllSecrets() (*iam_api.ListAllSecretResponse, error) {
	baseUrl := "/iam/internal/secrets"
	res := &iam_api.ListAllSecretResponse{}
	err := c.SendWithLog(http.MethodGet, baseUrl, nil, res, "ListAllSecrets")

	return res, err
}

func (c *IamClient) InternalAddSecret(req *iam_api.AdminAddSecretRequest) (*iam_api.AdminAddSecretResponse, error) {
	baseUrl := "/iam/internal/secrets"
	res := &iam_api.AdminAddSecretResponse{}
	err := c.SendWithLog("POST", baseUrl, req, res, "AdminAddSecret")
	return res, err
}

func (c *IamClient) InternalListSecrets(req *iam_api.AdminListSecretRequest) (*iam_api.AdminListSecretResponse, error) {
	baseUrl := fmt.Sprintf("/iam/internal/secrets/user/%s", req.UserId)
	res := &iam_api.AdminListSecretResponse{}
	err := c.SendWithLog("GET", baseUrl, nil, res, "AdminListSecret")
	return res, err
}

func (c *IamClient) AddAccount(req *iam_api.AddAccountRequest) (*iam_api.AddAccountResponse, error) {
	baseUrl := "/iam/v1/api/account"
	res := &iam_api.AddAccountResponse{}
	err := c.SendWithLog(http.MethodPost, baseUrl, req, res, "AddIAMAccount")
	return res, err
}

func (c *IamClient) ExchangeCredential(req *iam_api.ExchangeCredentialsRequest) (*iam_api.ExchangeCredentialsResponse, error) {
	baseUrl := "/iam/v1/api/account/exchange"
	res := &iam_api.ExchangeCredentialsResponse{}
	err := c.SendWithLog(http.MethodPost, baseUrl, req, res, "ChangeIAMCredential")
	return res, err
}

func (c *IamClient) Send(method, baseUrl string, inputData interface{}, resultData interface{},
	queryParam ...string) (*iam_api.BasicResponse, int, error) {
	targetUrl := fmt.Sprintf("%s/%s?AccessKeyId=%s&SessionToken=%s&Timestamp=%s",
		c.endpoint, strings.TrimLeft(baseUrl, "/"), c.akId, c.token, CurrentTimestamp())
	if len(queryParam) > 0 {
		// for each query param , add to targetUrl
		for _, param := range queryParam {
			targetUrl += fmt.Sprintf("&%s", param)
		}
	}
	var bodyReader io.Reader
	if inputData != nil {
		bodyBytes, err := json.Marshal(inputData)
		if err != nil {
			return nil, -1, err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}
	hReq, err := http.NewRequest(method, targetUrl, bodyReader)
	if err != nil {
		return nil, -1, err
	}
	hReq.Header.Set("Content-Type", "application/json")
	cred := credential.NewCredential(c.akId, c.akSecret)
	sig, err := signer.NewSigner(cred)
	signed, err := sig.SignHttp(hReq)
	if err != nil {
		return nil, -1, err
	}
	hReq.URL.RawQuery += fmt.Sprintf("&Signature=%s", signed.Signature)

	// 默认使用环境变量的proxy
	// 1. http.DefaultTransport默认使用http2，会复用TCP连接.
	// 2. 每次http请求会创建一个新的Stream.
	// 3. 当到达Stream数量上限时，http server端会发送GOAWAY信号，并要求重建TCP连接。
	// 4. http.DefaultTransport会自动重建TCP连接
	t := http.DefaultTransport
	t.(*http.Transport).ForceAttemptHTTP2 = false
	if c.Proxy != "" {
		proxyURL, err := url.Parse(c.Proxy)
		if err != nil {
			return nil, -1, err
		}
		t.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	}
	hc := http.Client{
		Transport: t,
	}
	var resp *http.Response
	defer closeResponse(resp)
	var respStr []byte
	resp, err = hc.Do(hReq)
	if err != nil {
		return nil, -1, err
	}
	respStr, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}
	baseRes := iam_api.BasicResponse{}
	if resultData != nil {
		baseRes.Data = resultData
	}
	err = json.Unmarshal(respStr, &baseRes)
	if err != nil {
		err = errors.New(fmt.Sprintf("Response: %s", string(respStr)))
	}
	if resp.StatusCode == 200 {
		return &baseRes, resp.StatusCode, nil
	}
	err = httpRespToErrorResponse(resp.StatusCode, &baseRes)
	return &baseRes, resp.StatusCode, err
}

func CurrentTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// ref minio-go api.go

// closeResponse close non nil response with any response Body.
// convenient wrapper to drain any remaining data on response body.
//
// Subsequently this allows golang http RoundTripper
// to re-use the same connection for future requests.
func closeResponse(resp *http.Response) {
	// Callers should close resp.Body when done reading from it.
	// If resp.Body is not closed, the Client's underlying RoundTripper
	// (typically Transport) may not be able to re-use a persistent TCP
	// connection to the server for a subsequent "keep-alive" request.
	if resp != nil && resp.Body != nil {
		// Drain any remaining Body and then close the connection.
		// Without this closing connection would disallow re-using
		// the same connection for future uses.
		//  - http://stackoverflow.com/a/17961593/4465767
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

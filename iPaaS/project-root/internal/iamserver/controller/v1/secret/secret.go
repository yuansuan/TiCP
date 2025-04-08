package secret

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/snowflake"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	iam_api "github.com/yuansuan/ticp/common/project-root-iam/iam-api"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/config"
	v1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/controller/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/pkg/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/rpc"
	srvv1 "github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/service/v1"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/iamserver/store/dao"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
	"net/http"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	"google.golang.org/grpc/status"
)

type SecretController struct {
	srv srvv1.Svc
}

func NewSecretController(s store.Factory) *SecretController {
	return &SecretController{
		srv: srvv1.NewSvc(s),
	}
}

func (s *SecretController) GetSecret(c *gin.Context) {
	akID := strings.TrimSpace(c.Param("accessKeyId"))
	if len(akID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid ak ID")
		return
	}
	secret, err := s.srv.Secrets().Get(c, akID)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.AccessKeyNotFound, "secret not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	userInfo := common.GetUserInfo(c)
	if secret.ParentUser != userInfo.UserID.String() && !srvv1.IsYuansuanProductAccount(userInfo.Tag) {
		common.ErrorRespWithAbort(c, http.StatusNotFound, v1.AccessKeyNotFound, "secret not found")
		return
	}

	res := &iam_api.GetSecretResponse{
		AccessKeyId:     secret.AccessKeyId,
		AccessKeySecret: secret.AccessKeySecret,
		YSId:            secret.ParentUser,
		Expire:          secret.Expiration,
	}
	common.SuccessResp(c, &res)
}

func (s *SecretController) IsYSProductAccount(c *gin.Context) {
	userID := strings.TrimSpace(c.Query("UserId"))
	if len(userID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, fmt.Sprintf("invalid user ID: %s", userID))
		return
	}
	userInfo := common.GetUserInfo(c)
	if !srvv1.IsYuansuanProductAccount(userInfo.Tag) && !userInfo.IsAdmin() {
		common.ErrorRespWithAbort(c, http.StatusForbidden, v1.PermissionDenied, "not allowed")
		return
	}
	u, err := s.srv.Secrets().GetByUserID(c, userID)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.UserNotFound, "user not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	var res iam_api.IsYSProductAccountResponse

	res.IsYSProductAccount = srvv1.IsYuansuanProductAccount(u.Tag)

	common.SuccessResp(c, &res)
}

func (s *SecretController) ListByParentUserID(c *gin.Context) {
	userInfo := common.GetUserInfo(c)

	secrets, err := s.srv.Secrets().ListByParentUserID(c, userInfo.UserID.String(), 0, 1000)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := make([]*iam_api.GetSecretResponse, 0)
	for _, secret := range secrets {
		res = append(res, &iam_api.GetSecretResponse{
			AccessKeyId:     secret.AccessKeyId,
			AccessKeySecret: secret.AccessKeySecret,
			YSId:            secret.ParentUser,
			Expire:          secret.Expiration,
		})
	}
	common.SuccessResp(c, &res)
}

func (s *SecretController) ListSecrets(c *gin.Context) {
	userInfo := common.GetUserInfo(c)
	if !srvv1.IsYuansuanProductAccount(userInfo.Tag) {
		common.ErrorResp(c, http.StatusForbidden, "Forbidden", "Only yusansuan product account can call this")
		return
	}
	// exclude assume role key
	// TODO: pagination in future
	secrets, err := s.srv.Secrets().List(c, 0, 10000)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}

	var sec []*iam_api.CacheSecret
	for _, v := range secrets {
		s := &iam_api.CacheSecret{
			AccessKeyId:     v.AccessKeyId,
			AccessKeySecret: v.AccessKeySecret,
			YSId:            v.ParentUser,
			Expire:          v.Expiration,
			Tag:             v.Tag,
		}
		sec = append(sec, s)
	}

	secretsRes := &iam_api.ListAllSecretResponse{
		Secrets: sec,
	}
	common.SuccessResp(c, secretsRes)
}

func (s *SecretController) CreateSecret(c *gin.Context) {
	userInfo := common.GetUserInfo(c)
	secret, err := s.srv.Secrets().CreateSecret(c, userInfo.UserID.String(), userInfo.Tag)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := &iam_api.AddSecretResponse{
		AccessKeyId:     secret.AccessKeyId,
		AccessKeySecret: secret.AccessKeySecret,
		YSId:            secret.ParentUser,
		Expire:          secret.Expiration,
	}
	common.SuccessResp(c, &res)
}

func (s *SecretController) DeleteSecretByParentUser(c *gin.Context) {
	userInfo := common.GetUserInfo(c)
	akID := strings.TrimSpace(c.Param("accessKeyId"))
	if len(akID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid ak ID")
		return
	}
	err := s.srv.Secrets().DeleteByParentUser(c, akID, userInfo.UserID.String())
	secretErrHandler(c, err)
}

// admin api
func (s *SecretController) AdminGetSecret(c *gin.Context) {
	akID := strings.TrimSpace(c.Param("accessKeyId"))
	if len(akID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid ak ID")
	}
	secret, err := s.srv.Secrets().Get(c, akID)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.AccessKeyNotFound, "secret not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}

	se := iam_api.GetSecretResponse{
		AccessKeyId:     secret.AccessKeyId,
		AccessKeySecret: secret.AccessKeySecret,
		YSId:            secret.ParentUser,
		Expire:          secret.Expiration,
	}
	res := &iam_api.AdminGetSecretResponse{
		Tag: secret.Tag,
	}
	res.GetSecretResponse = se
	common.SuccessResp(c, &res)
}

func (s *SecretController) InternalCreateSecret(c *gin.Context) {
	userInfo := common.GetUserInfo(c)
	if !srvv1.IsYuansuanProductAccount(userInfo.Tag) {
		common.ErrorResp(c, http.StatusForbidden, "Forbidden", "Only yusansuan product account can call this")
		return
	}
	s.AdminCreateSecret(c)
}

func (s *SecretController) AdminCreateSecret(c *gin.Context) {
	req := iam_api.AdminAddSecretRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	if len(req.UserId) == 0 || req.Tag == common.IamAdminTag {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "error User id or Tag")
		return
	}
	// tag length limit
	if len(req.Tag) > 255 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "tag length limit 255")
		return
	}
	secret, err := s.srv.Secrets().CreateSecret(c, req.UserId, req.Tag)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := &iam_api.AdminAddSecretResponse{
		AccessKeyId:     secret.AccessKeyId,
		AccessKeySecret: secret.AccessKeySecret,
		YSId:            secret.ParentUser,
		Expire:          secret.Expiration,
	}
	common.SuccessResp(c, &res)
}

func (s *SecretController) AdminDeleteSecret(c *gin.Context) {
	akID := strings.TrimSpace(c.Param("accessKeyId"))
	if len(akID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid ak ID")
		return
	}
	err := s.srv.Secrets().AdminDelete(c, akID)
	secretErrHandler(c, err)
}

func (s *SecretController) InvalidAppKey(c *gin.Context) {
	common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "empty app key")
}

func (s *SecretController) InternalListSecret(c *gin.Context) {
	userInfo := common.GetUserInfo(c)
	if !srvv1.IsYuansuanProductAccount(userInfo.Tag) {
		common.ErrorResp(c, http.StatusForbidden, "Forbidden", "Only yusansuan product account can call this")
		return
	}
	s.AdminListSecret(c)
}

func (s *SecretController) AdminListSecret(c *gin.Context) {
	userID := strings.TrimSpace(c.Param("userId"))
	if len(userID) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid user ID")
		return
	}
	secrets, err := s.srv.Secrets().ListByParentUserID(c, userID, 0, 1000)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}
	res := make([]*iam_api.AdminGetSecretResponse, 0, len(secrets))
	for _, secret := range secrets {
		res = append(res, &iam_api.AdminGetSecretResponse{
			Tag: secret.Tag,
			GetSecretResponse: iam_api.GetSecretResponse{
				AccessKeyId:     secret.AccessKeyId,
				AccessKeySecret: secret.AccessKeySecret,
				YSId:            secret.ParentUser,
				Expire:          secret.Expiration,
			},
		})
	}
	common.SuccessResp(c, &res)
}

func (s *SecretController) AdminUpdateTag(c *gin.Context) {
	req := iam_api.AdminUpdateTagRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "invalid argument")
		return
	}
	if len(req.AccessKeyId) == 0 || len(req.UserId) == 0 || req.Tag == common.IamAdminTag || len(req.Tag) == 0 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "error User id or Tag or AccessKeyId")
		return
	}
	// tag length limit
	if len(req.Tag) > 255 {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "tag length limit 255")
		return
	}
	err = s.srv.Secrets().UpdateTag(c, req.AccessKeyId, req.UserId, req.Tag)

	secretErrHandler(c, err)
}

func (s *SecretController) AdminListSecrets(c *gin.Context) {
	// get PageOffSet and PageSize from query
	pageOffset := c.Query("PageOffset")
	pageSize := c.Query("PageSize")
	if len(pageOffset) == 0 {
		pageOffset = "0"
	}
	if len(pageSize) == 0 {
		pageSize = "100"
	}
	// convert string to int
	offset, err := strconv.Atoi(pageOffset)
	if err != nil {
		offset = 0
	}
	limit, err := strconv.Atoi(pageSize)
	if err != nil {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	if limit < 0 || limit > 100 {
		limit = 100
	}

	secrets, err := s.srv.Secrets().List(c, offset, limit)
	if err != nil {
		common.InternalServerError(c, "")
		return
	}

	res := make([]*iam_api.AdminGetSecretResponse, 0, len(secrets))
	for _, secret := range secrets {
		res = append(res, &iam_api.AdminGetSecretResponse{
			Tag: secret.Tag,
			GetSecretResponse: iam_api.GetSecretResponse{
				AccessKeyId:     secret.AccessKeyId,
				AccessKeySecret: secret.AccessKeySecret,
				YSId:            secret.ParentUser,
				Expire:          secret.Expiration,
			},
		})
	}
	common.SuccessResp(c, &res)
}

func (s *SecretController) GetAKByPhone(c *gin.Context) {
	phone := strings.TrimSpace(c.Param("phone"))
	password := strings.TrimSpace(c.Query("password"))

	if phone == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "phone is empty")
		return
	}
	if password == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "password is empty, phone is: "+phone)
		return
	}

	user, err := rpc.GetInstance().QueryInfoByPhoneNumber(c, phone, password)
	if err != nil {
		s, ok := status.FromError(err)
		if ok && s.Code() == consts.ErrHydraLcpDBUserNotExist {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.UserNotFound, "user not exist")
			return
		}
		logging.Default().Warnf("query info by phone number error: %s, phone: %s", err.Error(), phone)
		common.InternalServerError(c, "")
		return
	}

	if user == nil {
		common.ErrorRespWithAbort(c, http.StatusNotFound, v1.UserNotFound, "user not exist")
		return
	}

	userID := user.Ysid

	secret, err := s.srv.Secrets().GetByUserID(c, userID)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			// create secret
			se, err := s.srv.Secrets().CreateSecret(c, userID, "")
			if err != nil {
				common.InternalServerError(c, "")
				return
			}
			secret = se
		}
	}
	resSecret := &iam_api.GetSecretResponse{
		AccessKeyId:     secret.AccessKeyId,
		AccessKeySecret: secret.AccessKeySecret,
		YSId:            secret.ParentUser,
		Expire:          secret.Expiration,
	}

	common.SuccessResp(c, resSecret)
}

func (s *SecretController) AddAccount(c *gin.Context) {
	u := common.GetUserInfo(c)
	users := config.GetConfig().AllowAddUsers
	if stringNotInSlice(u.UserID.String(), users) {
		common.ErrorRespWithAbort(c, http.StatusForbidden, v1.PermissionDenied, "you are not allowed to add ys account")
		return
	}
	req := &iam_api.AddAccountRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, err.Error())
		return
	}
	if req.Phone == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "phone is empty")
		return
	}
	if req.CompanyName == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "company name is empty")
		return
	}
	if req.Name == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "name is empty")
		return
	}
	if req.Password == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "password is empty")
		return
	}
	if req.UnifiedSocialCreditCode == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "unified social credit code is empty")
		return
	}
	if req.UserChannel == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "user channel is empty")
		return
	}
	if req.Email == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "email is empty")
		return
	}
	user := &hydra_lcp.AddUserReq{
		Phone:                   req.Phone,
		Name:                    req.Name,
		CompanyName:             req.CompanyName,
		UserChannel:             req.UserChannel,
		Password:                req.Password,
		UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
		Email:                   req.Email,
	}
	reply, err := rpc.GetInstance().AddUser2(c, user)

	if err != nil {
		state, ok := status.FromError(err)
		if ok && state.Code() == consts.ErrHydraLcpUserExist {
			logging.Default().Infof("user existed: %s", req.Phone)
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.AlreadyExists,
				fmt.Sprintf("account existed, phone: %s, email: %s, company name: %s", req.Phone, req.Email, req.CompanyName))
			return
		}
		logging.Default().Errorf("add user failed, err: %v", err)
		common.InternalServerError(c, "")
		return
	}
	// user secret not exist
	cre, err := s.srv.Secrets().CreateSecret(c, reply.UserId, "")
	if err != nil {
		logging.Default().Errorf("add user create secret failed, err: %v, request: %v", err, req)
		common.InternalServerError(c, "")
		return
	}
	resp := &iam_api.AddAccountResponse{
		YsId:                    reply.UserId,
		AccessKeyId:             cre.AccessKeyId,
		AccessKeySecret:         cre.AccessKeySecret,
		Phone:                   req.Phone,
		Name:                    req.Name,
		CompanyName:             req.CompanyName,
		UserChannel:             req.UserChannel,
		Password:                req.Password,
		UnifiedSocialCreditCode: req.UnifiedSocialCreditCode,
		Email:                   req.Email,
	}
	common.SuccessResp(c, resp)
}

func (s *SecretController) ExchangeCredentials(c *gin.Context) {
	req := &iam_api.ExchangeCredentialsRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, err.Error())
		return
	}

	if req.Email == "" && req.Phone == "" && req.YsId == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "Email or Phone or YsId is empty")
		return
	}
	if req.Password == "" {
		common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, "password is empty")
		return
	}
	if req.YsId != "" {
		if _, err := snowflake.ParseString(req.YsId); err != nil {
			common.ErrorRespWithAbort(c, http.StatusBadRequest, v1.InvalidArgument, fmt.Sprintf("ys id is in bad format: %s", req.YsId))
			return
		}
	}

	reply, err := rpc.GetInstance().CheckPassword2(c, req.YsId, req.Phone, req.Email, req.Password)

	if err != nil {
		status, ok := status.FromError(err)
		if ok && (status.Code() == consts.ErrHydraLcpDBUserNotExist || status.Code() == consts.ErrHydraLcpPwdNotMatch) {
			common.ErrorRespWithAbort(c, http.StatusForbidden, v1.PermissionDenied, "error password")
			return
		}
		logging.Default().Errorf("check password failed, err: %v", err)
		common.InternalServerError(c, "")
		return
	}

	cre, dbErr := getCredentials(c, s.srv, reply.Ysid)
	if dbErr != nil {
		logging.Default().Errorf("get credentials failed, err: %v", dbErr)
		common.InternalServerError(c, "")
		return
	}
	res := &iam_api.ExchangeCredentialsResponse{
		AccessKeyId:     cre.AccessKeyId,
		AccessKeySecret: cre.AccessKeySecret,
		YsId:            reply.Ysid,
		Email:           reply.Email,
		Name:            reply.Name,
		Phone:           reply.Phone,
	}
	common.SuccessResp(c, res)
}
func secretErrHandler(c *gin.Context, err error) {
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			common.ErrorRespWithAbort(c, http.StatusNotFound, v1.AccessKeyNotFound, "secret not found")
			return
		}
		common.InternalServerError(c, "")
		return
	}
	common.SuccessResp(c, nil)
}

func stringNotInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return false
		}
	}
	return true
}

func getCredentials(c *gin.Context, srv srvv1.Svc, ysid string) (*dao.Secret, error) {
	secret, err := srv.Secrets().GetByUserID(c, ysid)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			se, err := srv.Secrets().CreateSecret(c, ysid, "")
			if err != nil {
				return nil, err
			}
			secret = se
		} else {
			return nil, err
		}
	}
	return secret, nil
}

package handler

import (
	"context"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	myModels "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
)

// Signup Signup
// swagger:route POST /api/signup fe SignupReq
//
// signup
//
//	Responses:
//		200: SignupResp
//		90002: ErrHydraLcpBadRequest
//
// @POST /api/signup
func (h *Handler) Signup(c *gin.Context) {
	var req SignupReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}
	var i int
	i++

	re := util.IsPasswordValid(req.Password)
	if !re {
		http.Errf(c, consts.ErrHydraLcpPwdInvalidate, "password is invalid")
		return
	}
	err := h.phoneSrv.VerifyCode(c, req.Phone, req.PhoneCode)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	if req.UserSource == "" {
		req.UserSource = "未知来源"
	}
	if req.UserChannel == "" {
		req.UserChannel = "UNKNOWN"
	}
	//	get userID, if user not exist, add User
	userID, err := h.addUser(c, req.Phone, "", req.Password, req.Name, req.Company, req.Email, req.UserSource, req.UserChannel)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	// encode ysid with base58
	id58 := snowflake.ID(userID).String()

	http.Ok(c, SignupResp{UserID: id58})
}

// SignupByName SignupByName
// swagger:route POST /api/signupbyname fe SignupReq
//
// 柳汽临时使用
//
//	Responses:
//		200: SignupResp
//		90002: ErrHydraLcpBadRequest
//
// @POST /api/signupbyname
func (h *Handler) SignupByName(c *gin.Context) {
	var req SignupByNameReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}

	// check sign key
	key := "6HIQEad22365fJQEw7fgT"
	if req.SignupKey == "" || req.SignupKey != key {
		http.Err(c, consts.ErrHydraLcpBadRequest, "check signup failed")
		return
	}

	// check name and password
	// TODO :

	//	get userID, if user not exist, add User
	ssoModel := &models.SsoUser{Name: req.Name}
	ok, err := h.userSrv.Get(c, ssoModel)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	if ok {
		http.ErrFromGrpc(c, status.Errorf(consts.ErrHydraLcpUserExist, "user(%v) exist", req.Name))
		return
	}

	id, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}

	err = h.userSrv.Add(c, &myModels.SsoUser{Ysid: id.Id, Name: req.Name, WechatUnionId: "", RealName: req.RealName}, req.Password)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}

	// encode ysid with base58
	id58 := id.String()

	http.Ok(c, SignupResp{UserID: id58})
}

// SignupByNameReq ...
type SignupByNameReq struct {
	Name      string `json:"user_name"`
	Password  string `json:"password"`
	SignupKey string `json:"signup_key"`
	RealName  string `json:"real_name"`
}

// SignupResp SignupResp
// swagger:response SignupResp
type SignupResp struct {
	UserID string `json:"user_id"`
}

// SignupReq SignupReq
// swagger:parameters SignupReq
type SignupReq struct {
	Phone       string `json:"phone"`
	PhoneCode   string `json:"phone_code"`
	Password    string `json:"password"`
	UserSource  string `json:"user_source"`
	UserChannel string `json:"user_channel"`
	Company     string `json:"company"`
	Email       string `json:"email"`
	Name        string `json:"name"`
}

// Resend Resend
// @POST /api/resend
func (h *Handler) Resend(c *gin.Context) {

}

// Activate Activate
// swagger:route POST /api/activate fe ActivateReq
//
// activate email
//
//	Responses:
//		200: ActivateResp
//		90002: ErrHydraLcpBadRequest
//
// @POST /api/activate
func (h *Handler) Activate(c *gin.Context) {
	var req ActivateReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}

	userID, err := h.userSrv.VerifyPasswordByEmail(c, req.Email, req.Password)
	if err != nil {
		http.ErrFromGrpc(c, err)
	}

	err = h.emailSrv.Activate(c, req.Token, userID)
	if err != nil {
		http.ErrFromGrpc(c, err)
	}

	http.Ok(c, ActivateResp{UserID: userID})
}

// ActivateReqWrapper ActivateReqWrapper
// swagger:parameters ActivateReq
type ActivateReqWrapper struct {
	// in: body
	Req ActivateReq
}

// ActivateReq ActivateReq
type ActivateReq struct {
	// required: true
	Email string `json:"email"`
	// required: true
	Password string `json:"password"`
	// required: true
	Token string `json:"token"`
}

// ActivateResp ActivateResp
// swagger:response ActivateResp
type ActivateResp struct {
	UserID int64 `json:"user_id"`
}

func (h *Handler) addUser(c context.Context, phone string, wechatUnionID string, pwd string, name, company, email, source, channel string) (int64, error) {
	//	get userID, if user not exist, add User
	_, err := h.userSrv.GetID(c, email, phone, wechatUnionID)
	if err != nil {
		s, _ := status.FromError(err)
		if s.Code() == consts.ErrHydraLcpDBUserNotExist {
			id, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
			if err != nil {
				return 0, err
			}
			err = h.userSrv.Add(c, &myModels.SsoUser{
				Name:          name,
				Ysid:          id.Id,
				Phone:         phone,
				Company:       company,
				Email:         email,
				WechatUnionId: wechatUnionID,
				UserChannel:   channel,
				UserSource:    source,
			}, pwd)
			if err != nil {
				return 0, err
			}
			userID := id.Id

			return userID, err
		}
	}

	return 0, status.Errorf(consts.ErrHydraLcpUserExist, "user(%v) exist", phone)
}

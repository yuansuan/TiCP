package handler

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	myModels "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"
)

// ResetPassword ResetPassword
// swagger:route POST /api/reset_password fe ResetPasswordReq
//
//	 reset password
//
//		Responses:
//	     200: ResetPasswordResp
//	     90002: ErrHydraLcpBadRequest
//	     90021: ErrHydraLcpUserExist
//	     90029: ErrHydraLcpCaptchaVerifyFailed
//	     90036: ErrHydraLcpPwdInvalidate
//
// @POST /api/reset_password
// @Example:
//
//	curl http://127.0.0.1:8899/api/reset_password -d '{"phone":"15910957558", "phone_code":"334370", "password":"12345678ff”}'
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}

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

	// update user passwd
	userID, err := h.resetUserPwd(c, req.Phone, "", req.Password)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	// encode ysid with base58
	id58 := snowflake.ID(userID)
	http.Ok(c, ResetPasswordResp{UserID: id58.String()})
}

func (h *Handler) ResetEmailPassword(c *gin.Context) {
	var req ResetEmailPasswordReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}

	re := util.IsPasswordValid(req.Password)
	if !re {
		http.Errf(c, consts.ErrHydraLcpPwdInvalidate, "password is invalid")
		return
	}

	err := h.emailSrv.VerifyCode(c, req.Email, req.EmailCode)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}

	// update user passwd
	userID, err := h.resetEmailPwd(c, req.Email, req.Password)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}

	// encode ysid with base58
	id58 := snowflake.ID(userID)
	http.Ok(c, ResetPasswordResp{UserID: id58.String()})
}

type ResetEmailPasswordReq struct {
	Email     string `json:"email"`
	EmailCode string `json:"email_code"`
	Password  string `json:"password"`
}

// ResetPasswordResp ResetPasswordResp
// swagger:response ResetPasswordResp
type ResetPasswordResp struct {
	UserID string `json:"user_id"`
}

// ResetPasswordReq ResetPasswordReq
// swagger:parameters ResetPasswordReq
type ResetPasswordReq struct {
	Phone     string `json:"phone"`
	PhoneCode string `json:"phone_code"`
	Password  string `json:"password"`
}

func (h *Handler) resetUserPwd(c *gin.Context, phone string, wechatUnionID string, pwd string) (int64, error) {
	// get userID
	userID, err := h.userSrv.GetID(c, "", phone, wechatUnionID)
	// if userID doesn't exist, add user with phone number automatically
	// if userID exists, change user pwd
	if err != nil {
		s, _ := status.FromError(err)
		if s.Code() == consts.ErrHydraLcpDBUserNotExist {
			id, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
			if err != nil {
				return 0, err
			}
			err = h.userSrv.Add(c, &myModels.SsoUser{Ysid: id.Id, Phone: phone, WechatUnionId: wechatUnionID}, pwd)
			if err != nil {
				return 0, err
			}
			userID := id.Id
			return userID, nil
		}
	} else {
		// get user's old sso info
		ssouserinfo := myModels.SsoUser{
			Ysid: userID,
		}
		ok, err := h.userSrv.Get(c, &ssouserinfo)
		if nil != err {
			return 0, err
		}
		if !ok {
			return 0, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
		}
		// update user pwd
		return userID, h.userSrv.UpdatePwd(c, ssouserinfo, pwd)
	}
	return 0, err
}

func (h *Handler) resetEmailPwd(c *gin.Context, email, pwd string) (int64, error) {

	// get user's old sso info
	ssouserinfo := myModels.SsoUser{
		Email: email,
	}
	ok, err := h.userSrv.Get(c, &ssouserinfo)
	if nil != err {
		return 0, err
	}
	if !ok {
		return 0, status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}

	// update user pwd
	return ssouserinfo.Ysid, h.userSrv.UpdatePwd(c, ssouserinfo, pwd)
}

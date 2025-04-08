package handler

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common/ssojwt"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/company"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	myModels "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// 由于微信H5登陆仅做查询，不涉及到SSO及其他系统，甚至没有加密必要，采用固定serect (UUID)
const secret = "32eab94d-f0d8-4c70-8184-657411dea7b1"

// MarketEventLoginReq MarketEventLoginReq
// swagger:parameters MarketEventLoginReq
type MarketEventLoginReq struct {
	// in: query
	Phone     string `json:"phone"`
	PhoneCode string `json:"phone_code"`
	Password  string `json:"password"`
	//additional info
	WechatNicName string `json:"wechat_nick_name"`
	OpenID        string `json:"openid"`
	UnionID       string `json:"unionid"`
	Headimg       string `json:"headimg_url"`
	UserChannel   string `json:"user_channel"`
	UserSource    string `json:"user_source"`
	// need ysid
	UserReferer string `json:"user_referer"`
}

// MarketEventLoginResp MarketEventLoginResp
// swagger:response MarketEventLoginResp
type MarketEventLoginResp struct {
	UserToken     string      `json:"user_token"`
	ExpiresAt     int64       `json:"expires_at"`
	UserInfo      *UserInfo   `json:"user_info"`
	ReferedUsers  []*UserInfo `json:"refered_users"`
	TotalReferers int         `json:"total_referers"`
	IsRegister    bool        `json:"is_register"`
}

// UserInfo UserInfo
type UserInfo struct {
	UserID         string `json:"user_id"`
	Name           string `json:"name"`
	RealName       string `json:"real_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	WechatUnionId  string `json:"unionid"`
	WechatOpenId   string `json:"openid"`
	WechatNickName string `json:"wechat_nick_name"`
	HeadimgUrl     string `json:"headimg_url"`
	UserChannel    string `json:"user_channel"`
	UserSource     string `json:"user_source"`
	UserReferer    string `json:"user_referer"`
}

// MarketEventLogin MarketEventLogin
// swagger:route POST /api/mkeventlogin expand_login clear
//
// login api used by market event, may merge to MarketService micro service later.
//
// Responses:
//
//		     200: MarketEventLoginResp
//	      90001: ErrHydraLcpFailedToReqHydra
//
// @POST /api/mkevent/login
func (h *Handler) MarketEventLogin(c *gin.Context) {
	var req MarketEventLoginReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "err: %v", err)
		return
	}
	logger := logging.GetLogger(c)
	logger.Infof("[Market Event login] start expand login for user with phone %v", req.Phone)

	if err := h.phoneSrv.VerifyCode(c, req.Phone, req.PhoneCode); err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	var resp MarketEventLoginResp
	resp.IsRegister = false
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	userID, err := h.userSrv.GetID(ctx, "", req.Phone, "")
	logger.Debug("template string, args ...interface{}")
	//新用户注册
	if err != nil {
		s, _ := status.FromError(err)
		if s.Code() == consts.ErrHydraLcpDBUserNotExist {

			reply, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
			id := snowflake.ID(reply.Id)
			if err != nil {
				http.ErrFromGrpc(c, err)
				return
			}

			if req.UserReferer != "" {
				refererInfo, err := h.CompanyUserService.GetUserInfo(ctx, &company.GetUserInfoRequest{
					UserId: req.UserReferer,
				})
				if err != nil {
					// referer填错了，不处理，不阻止用户注册 referer置零(base 58为“1”)，便于运维明确问题
					logger.Errorf("[Market Event Register], wrong referer id: %v, set to 0", req.UserReferer)
					req.UserReferer = "1"
				} else {
					if refererInfo.RealName == "" {
						refererInfo.RealName = refererInfo.Phone
					}
					req.UserSource = fmt.Sprintf("%v(%v)的分享", refererInfo.RealName, refererInfo.UserId)
				}
			}
			//允许已绑定的微信用户用新手机注册，但不允许绑定微信
			ok, err := h.userSrv.Get(ctx, &myModels.SsoUser{
				WechatUnionId: req.UnionID,
			})
			if err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
			if ok {
				req.UnionID = ""
			}
			if req.Password == "" {
				req.Password = randPwd()
			}
			err = h.userSrv.Add(c, &myModels.SsoUser{
				Ysid:           id.Int64(),
				Phone:          req.Phone,
				WechatUnionId:  req.UnionID,
				WechatOpenId:   req.OpenID,
				WechatNickName: req.WechatNicName,
				HeadimgUrl:     req.Headimg,
				UserChannel:    req.UserChannel,
				UserSource:     req.UserSource,
				UserReferer:    snowflake.MustParseString(req.UserReferer).Int64(),
			}, req.Password)

			if err != nil {
				http.ErrFromGrpc(c, err)
				return
			}

			_, err = h.CompanyUserService.UserInit(ctx, &company.UserInitRequest{
				UserId:   id.String(),
				Phone:    req.Phone,
				RealName: req.WechatNicName,
			})
			if err != nil {
				http.ErrFromGrpc(c, err)
				return
			}
			userID = id.Int64()
			resp.IsRegister = true
		}
	}
	id := snowflake.ID(userID)
	resp.UserToken, err = ssojwt.GenerateJwtToken(&ssojwt.UserInfo{
		UserID:          id.String(),
		CookieExpiredAt: time.Now().Add(168 * time.Hour).Unix(),
	}, secret, []byte(secret), time.Now().Add(168*time.Hour).Unix())
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	resp.ExpiresAt = time.Now().Add(168 * time.Hour).Unix()
	err = h.getUserInfo(ctx, id.Int64(), &resp)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	//老用户微信绑定
	if resp.UserInfo.WechatUnionId == "" {
		//更新数据，更新失败数据无需回退，忽略此处err
		h.userSrv.UpdateWechatInfo(ctx, myModels.SsoUser{
			Ysid:           id.Int64(),
			WechatNickName: req.WechatNicName,
			WechatOpenId:   req.OpenID,
			WechatUnionId:  req.UnionID,
			HeadimgUrl:     req.Headimg,
		})
	}
	http.Ok(c, resp)
}

// GetUserInfo GetUserInfo
// swagger:route GET /api/mkevent/getuserinfo expand_login clear
//
// get user info api used by market event, may merge to MarketService micro service later.
//
//	    Responses:
//		     200: MarketEventLoginResp
//	      90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/mkevent/getuserinfo
func (h *Handler) GetUserInfo(c *gin.Context) {
	jwt := c.Query(common.HydraLoginJwt)
	if jwt == "" {
		http.Err(c, consts.ErrHydraLcpJwtAuthError, "token is null")
		return
	}

	info, err := ssojwt.ParseJwtToken(jwt, func(secretID string) ([]byte, error) {
		return []byte(secret), nil
	})
	if err != nil {
		http.Err(c, consts.ErrHydraLcpJwtAuthError, "error parse token")
		return
	}

	if info.CookieExpiredAt < time.Now().Unix() {
		http.Err(c, consts.ErrHydraLcpJwtAuthError, "token expeired")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := snowflake.MustParseString(info.UserID)
	userInfo := myModels.SsoUser{
		Ysid: id.Int64(),
	}
	err = h.userSrv.CheckUserExists(ctx, userInfo)
	if err != nil {
		http.Err(c, consts.ErrHydraLcpUserExist, "user do not exist")
		return
	}

	var resp MarketEventLoginResp
	resp.UserToken, err = ssojwt.GenerateJwtToken(&ssojwt.UserInfo{
		UserID:          id.String(),
		CookieExpiredAt: time.Now().Add(168 * time.Hour).Unix(),
	}, secret, []byte(secret), time.Now().Add(168*time.Hour).Unix())
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	resp.ExpiresAt = time.Now().Add(168 * time.Hour).Unix()
	err = h.getUserInfo(ctx, id.Int64(), &resp)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}
	resp.IsRegister = false
	http.Ok(c, resp)
}

func (h *Handler) getUserInfo(ctx context.Context, userID int64, resp *MarketEventLoginResp) error {
	var userInfo myModels.SsoUser
	userInfo.Ysid = userID
	ok, err := h.userSrv.Get(ctx, &userInfo)
	if err != nil {
		return err
	}
	if !ok {
		return status.Errorf(consts.ErrHydraLcpDBUserNotExist, "user does not exist")
	}
	refererInfo, total, err := h.userSrv.GetReferedUsers(ctx, userID)
	resp.UserInfo = modelToRespUserInfo(&userInfo)
	for _, referer := range refererInfo {
		resp.ReferedUsers = append(resp.ReferedUsers, modelToRespUserInfo(referer))
	}
	resp.TotalReferers = total
	return nil
}

func modelToRespUserInfo(model *models.SsoUser) *UserInfo {
	userInfo := &UserInfo{
		UserID:         snowflake.ID(model.Ysid).String(),
		Name:           model.Name,
		RealName:       model.RealName,
		Email:          model.Email,
		Phone:          model.Phone,
		WechatUnionId:  model.WechatUnionId,
		WechatOpenId:   model.WechatOpenId,
		WechatNickName: model.WechatNickName,
		HeadimgUrl:     model.HeadimgUrl,
		UserChannel:    model.UserChannel,
		UserSource:     model.UserSource,
		UserReferer:    snowflake.ID(model.UserReferer).String(),
	}
	return userInfo
}

func randPwd() string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < 20; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

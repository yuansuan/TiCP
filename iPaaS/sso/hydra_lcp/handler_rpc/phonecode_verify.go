package handler_rpc

import (
	"context"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common/ssojwt"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	hydra_lcp "github.com/yuansuan/ticp/iPaaS/sso/protos"
)

// 仅做状态验证，使用固定serect
const secret = "fc78a29f99ded736526af21c4f09f49b870032a9ef6254cb9ad8df96fcdde258"

// VerifyPhoneCode  验证手机验证码
func (h *HydraLcpService) VerifyPhoneCode(ctx context.Context, req *hydra_lcp.VerifyPhoneCodeReq) (*hydra_lcp.VerifyPhoneCodeResp, error) {
	phone := strings.TrimSpace(req.Phone)
	code := strings.TrimSpace(req.Code)
	// 手机号验证
	if m, _ := regexp.MatchString("^1[0-9]{10}$", phone); !m {
		return nil, status.Error(consts.ErrHydraLcpPhoneInvalidate, "phone invalidate")
	}
	err := h.phoneSvr.VerifyCode(ctx, phone, code)
	if err != nil {
		return nil, err
	}
	var resp hydra_lcp.VerifyPhoneCodeResp
	resp.Token, err = ssojwt.GenerateJwtToken(&ssojwt.UserInfo{
		UserID:          phone,
		CookieExpiredAt: time.Now().Add(10 * time.Minute).Unix(),
	}, secret, []byte(secret), time.Now().Add(10*time.Minute).Unix())
	if err != nil {
		return nil, err
	}
	resp.IsSucceed = true
	return &resp, nil
}

// VerifyJwtToken 验证手机验证码Token
func (h *HydraLcpService) VerifyJwtToken(ctx context.Context, req *hydra_lcp.VerifyJwtTokenReq) (*hydra_lcp.VerifyPhoneCodeResp, error) {
	phone := strings.TrimSpace(req.Phone)
	token := strings.TrimSpace(req.Token)
	// 手机号验证
	if m, _ := regexp.MatchString("^1[0-9]{10}$", phone); !m {
		return nil, status.Error(consts.ErrHydraLcpPhoneInvalidate, "phone invalidate")
	}
	// token验证
	info, err := ssojwt.ParseJwtToken(token, func(secretID string) ([]byte, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	tokenPhone := strings.TrimSpace(info.UserID)
	if strings.Compare(phone, tokenPhone) != 0 {
		return nil, status.Errorf(consts.ErrHydraLcpPhoneJwtAuthError, "token and phone do not match, please check")
	}
	var resp hydra_lcp.VerifyPhoneCodeResp
	//新token交换
	resp.Token, err = ssojwt.GenerateJwtToken(&ssojwt.UserInfo{
		UserID:          phone,
		CookieExpiredAt: time.Now().Add(10 * time.Minute).Unix(),
	}, secret, []byte(secret), time.Now().Add(10*time.Minute).Unix())
	if err != nil {
		return nil, err
	}
	resp.IsSucceed = true
	return &resp, nil
}

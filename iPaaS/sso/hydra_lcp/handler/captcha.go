package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
)

// CreateCaptchaReqWrapper CreateCaptchaReqWrapper
// swagger:parameters CreateCaptchaReq
type CreateCaptchaReqWrapper struct {
	// in: body
	Req CreateCaptchaReq
}

// CreateCaptchaReq CreateCaptchaReq
type CreateCaptchaReq struct {
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	ImageCaptchaID string `json:"image_captcha_id"`
}

// CreateCaptchaRespWrapper CreateCaptchaRespWrapper
// swagger:response CreateCaptchaResp
type CreateCaptchaRespWrapper struct {
	// in: body
	Resp CreateCaptchaResp
}

// CreateCaptchaResp CreateCaptchaResp
type CreateCaptchaResp struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

// CreateCaptcha CreateCaptcha
// swagger:route POST /api/captcha fe CreateCaptchaReq
//
// create captcha image, get captcha id
//
//	    Responses:
//		     90002: CreateCaptchaResp
//
// @POST /api/captcha
func (h *Handler) CreateCaptcha(c *gin.Context) {
	var req CreateCaptchaReq
	if err := c.BindJSON(&req); err != nil {
		http.Errf(c, consts.ErrHydraLcpBadRequest, "failed to bind request, err: %v", err)
		return
	}

	_, image, err := h.captchaSrv.CreateImageCaptcha(c, req.Width, req.Height, req.ImageCaptchaID)
	if err != nil {
		http.ErrFromGrpc(c, err)
		return
	}

	http.Ok(c, CreateCaptchaResp{
		ID:    req.ImageCaptchaID,
		Image: image,
	})
}

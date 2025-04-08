package handler

import "github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

// HydraConsentRespWrapper HydraConsentRespWrapper
// 情况1: request的with_hydra为true但没有带login_challenge.
// swagger:response ErrHydraLcpFailedToReqHydra
type HydraConsentRespWrapper struct {
	// in: body
	Resp struct {
		http.SwaggerErrorMeta
		// example: 90001
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}

// ErrHydraLcpBadRequestWrapper ErrHydraLcpBadRequestWrapper
// swagger:response ErrHydraLcpBadRequest
type ErrHydraLcpBadRequestWrapper struct {
	// in: body
	Resp struct {
		http.SwaggerErrorMeta
		// example: 90002
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}

// ErrHydraLcpWechatGetOpenIDWrapper ErrHydraLcpWechatGetOpenIDWrapper
// swagger:response ErrHydraLcpWechatGetOpenID
type ErrHydraLcpWechatGetOpenIDWrapper struct {
	// in: body
	Resp struct {
		http.SwaggerErrorMeta
		// example: 90019
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}

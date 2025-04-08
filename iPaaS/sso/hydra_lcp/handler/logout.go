package handler

import (
	http2 "net/http"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
)

// HydraLogout HydraLogout
// swagger:route GET /api/hydra/logout hydra HydraLogoutReq
//
//		logout api used by hydra, it's behavior is strictly defined by hydra.
//
//	  Responses:
//			302
//			90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/hydra/logout
func (h *Handler) HydraLogout(c *gin.Context) {
	challenge := c.Query(common.HydraLogoutChallenge)
	logger := logging.GetLogger(c)
	logger.Infof("[hydra logout] start hydra logout for challenge %v", challenge)
	logger.Info(">>>>>>>>>>>>>>> challenge: ", challenge)

	res, err := h.HydraClient.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().WithLogoutChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra logout exception] unable to send request for logout to hydra: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for logout to hydra: %v", err)
		return
	}

	logger.Info(">>>>>>>>>>>>>>> response: ", res)

	acceptRes, err := h.HydraClient.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().
		WithLogoutChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra logout exception] failed to send request to hydra for accepting logout, err: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "failed to send request to hydra for accepting logout, err: %v", err)
		return
	}

	logger.Infof("[hydra logout] redirect to %v", acceptRes.Payload.RedirectTo)
	c.Redirect(http2.StatusFound, acceptRes.Payload.RedirectTo)
	logger.Infof("[hydra logout] hydra logout successful for challenge %v", challenge)
}

// HydraPortalLogout HydraPortalLogout
// swagger:route GET /api/hydra_portal/logout hydra_portal HydraLogoutReq
//
//		logout api used by hydra-portal, it's behavior is strictly defined by hydra-portal.
//
//	  Responses:
//			302
//			90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/hydra_portal/logout
func (h *Handler) HydraPortalLogout(c *gin.Context) {
	challenge := c.Query(common.HydraLogoutChallenge)
	logger := logging.GetLogger(c)
	logger.Infof("[hydra-portal logout] start hydra-portal logout for challenge %v", challenge)
	logger.Info(">>>>>>>>>>>>>>> challenge: ", challenge)

	res, err := h.HydraClientPortal.Admin.GetLogoutRequest(admin.NewGetLogoutRequestParams().WithLogoutChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra-portal logout exception] unable to send request for logout to hydra-portal: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to send request for logout to hydra-portal: %v", err)
		return
	}

	logger.Info(">>>>>>>>>>>>>>> response: ", res)

	acceptRes, err := h.HydraClientPortal.Admin.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParams().
		WithLogoutChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra-portal logout exception] failed to send request to hydra-portal for accepting logout, err: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "failed to send request to hydra-portal for accepting logout, err: %v", err)
		return
	}

	logger.Infof("[hydra-portal logout] redirect to %v", acceptRes.Payload.RedirectTo)
	c.Redirect(http2.StatusFound, acceptRes.Payload.RedirectTo)
	logger.Infof("[hydra-portal logout] hydra-portal logout successful for challenge %v", challenge)
}

// HydraLogoutReq HydraLogoutReq
// swagger:parameters HydraLogoutReq
type HydraLogoutReq struct {
	// required: true
	// in: param
	LogoutChallenge string `json:"logout_challenge"`
}

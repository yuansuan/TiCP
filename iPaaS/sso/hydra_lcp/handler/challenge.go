package handler

import (
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
)

// ToLoginURLResp ...
type ToLoginURLResp struct {
	URL string `json:"url"`
}

// HydraToLoginURL HydraToLoginURL
// swagger:route GET /api/hydra/challengetologinurl   ToLoginURLResp
// Get login url by challenge
//
//	    Responses:
//		     200
//	      90037: ErrHydraLcpChallengeExpired
//
// @GET /api/hydra/challengetologinurl
func (h *Handler) HydraToLoginURL(c *gin.Context) {
	challenge := c.Query(common.HydraLoginChallenge)
	// get login request from hydra
	res, err := h.HydraClient.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(challenge))
	if err != nil {
		http.Errf(c, consts.ErrHydraLcpChallengeExpired, "login challenge not found: %v", err)
		return
	}
	http.Ok(c, ToLoginURLResp{URL: res.Payload.RequestURL})
}

// HydraPortalToLoginURL HydraPortalToLoginURL
// swagger:route GET /api/hydra_portal/challengetologinurl   ToLoginURLResp
// Get login url by challenge
//
//	    Responses:
//		     200
//	      90037: ErrHydraLcpChallengeExpired
//
// @GET /api/hydra_portal/challengetologinurl
func (h *Handler) HydraPortalToLoginURL(c *gin.Context) {
	challenge := c.Query(common.HydraLoginChallenge)
	// get login request from hydra
	res, err := h.HydraClientPortal.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(challenge))
	if err != nil {
		http.Errf(c, consts.ErrHydraLcpChallengeExpired, "login challenge not found: %v", err)
		return
	}
	http.Ok(c, ToLoginURLResp{URL: res.Payload.RequestURL})
}

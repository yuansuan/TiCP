package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"
	"github.com/ory/hydra/sdk/go/hydra/models"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"
	modelsYs "github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/dao/models"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"
	http2 "net/http"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// HydraConsent HydraConsent
// swagger:route GET /api/hydra/consent  hydra HydraConsentReq
//
// consent api used by hydra, it's behavior is strictly defined by hydra.
//
//	    Responses:
//			 302
//			 90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/hydra/consent
func (h *Handler) HydraConsent(c *gin.Context) {
	challenge := c.Query(common.HydraConsentChallenge)
	logger := logging.GetLogger(c)
	logger.Infof("[hydra consent] start hydra consent for challenge %v", challenge)
	logger.Info(">>>>>>>>>>>>>>> challenge: ", challenge)

	res, err := h.HydraClient.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra consent exception] unable to fetch consent request: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "Unable to fetch consent request: %v", err)
		return
	}

	sid := snowflake.MustParseString(res.Payload.Subject).Int64()
	user := modelsYs.SsoUser{Ysid: sid}
	_, err = h.userSrv.Get(c, &user)
	if err != nil {
		logger.Warnf("[hydra consent exception] fail to get user: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "fail to get user: %v", err)
		return
	}

	//Scope内资源解析
	requestScope := res.Payload.RequestedScope
	scopeInfo := map[string]interface{}{}
	for _, str := range requestScope {
		switch str {
		case "email":
			scopeInfo["email"] = user.Email
		case "name":
			scopeInfo["name"] = user.Name
		case "phone":
			scopeInfo["phone"] = user.Phone
		case "avatar":
			scopeInfo["avatar"] = user.HeadimgUrl
		case "company":
			scopeInfo["company"] = user.Company
		default:
			logger.Infof("scope_info %v no exist", str)
		}
	}

	logger.Info(">>>>>>>>>>>>>>> response: ", res)
	acceptRes, err := h.HydraClient.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
		WithConsentChallenge(challenge).
		WithBody(&models.HandledConsentRequest{
			GrantedScope: res.Payload.RequestedScope,
			Remember:     true,
			Session: &models.ConsentRequestSessionData{
				IDToken: scopeInfo,
			},
			RememberFor: h.TokenExpireTime,
		},
		))
	if err != nil {
		logger.Warnf("[hydra consent exception] unable to accept hydra consent request: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to accept hydra consent request: %v", err)
		return
	}

	logger.Infof("[hydra consent] redirect to %v", acceptRes.Payload.RedirectTo)
	c.Redirect(http2.StatusFound, acceptRes.Payload.RedirectTo)
	logger.Infof("[hydra consent] hydra consent successful for challenge %v", challenge)
}

// HydraPortalConsent HydraPortalConsent
// swagger:route GET /api/hydra_portal/consent  hydra_portal HydraConsentReq
//
// consent api used by hydra-portal, it's behavior is strictly defined by hydra-portal.
//
//	    Responses:
//			 302
//			 90001: ErrHydraLcpFailedToReqHydra
//
// @GET /api/hydra_portal/consent
func (h *Handler) HydraPortalConsent(c *gin.Context) {
	challenge := c.Query(common.HydraConsentChallenge)
	logger := logging.GetLogger(c)
	logger.Infof("[hydra-portal consent] start hydra-portal consent for challenge %v", challenge)
	logger.Info(">>>>>>>>>>>>>>> challenge: ", challenge)

	res, err := h.HydraClientPortal.Admin.GetConsentRequest(admin.NewGetConsentRequestParams().WithConsentChallenge(challenge))
	if err != nil {
		logger.Warnf("[hydra-portal consent exception] unable to fetch consent request: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "Unable to fetch consent request: %v", err)
		return
	}
	logger.Info(">>>>>>>>>>>>>>> response: ", res)
	acceptRes, err := h.HydraClientPortal.Admin.AcceptConsentRequest(admin.NewAcceptConsentRequestParams().
		WithConsentChallenge(challenge).
		WithBody(&models.HandledConsentRequest{
			GrantedScope: res.Payload.RequestedScope,
			Remember:     true,
			Session: &models.ConsentRequestSessionData{
				IDToken: res.Payload.Context,
			},
			RememberFor: h.TokenExpireTimePortal,
		},
		))
	if err != nil {
		logger.Warnf("[hydra-portal consent exception] unable to accept hydra-portal consent request: %v", err)
		http.Errf(c, consts.ErrHydraLcpFailedToReqHydra, "unable to accept hydra-portal consent request: %v", err)
		return
	}

	logger.Infof("[hydra-portal consent] redirect to %v", acceptRes.Payload.RedirectTo)
	c.Redirect(http2.StatusFound, acceptRes.Payload.RedirectTo)
	logger.Infof("[hydra-portal consent] hydra-portal consent successful for challenge %v", challenge)
}

// HydraConsentReqWrapper HydraConsentReqWrapper
// swagger:parameters HydraConsentReq
type HydraConsentReqWrapper struct {
	// required: true
	// in: query
	ConsentChallenge string `json:"consent_challenge"`
}

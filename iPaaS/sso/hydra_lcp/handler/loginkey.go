package handler

import (
	"time"

	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"

	"github.com/gin-gonic/gin"

	//"github.com/ory/hydra/sdk/go/hydra/client/admin"

	"github.com/yuansuan/ticp/iPaaS/sso/protos/platform/idgen"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/pkg/snowflake"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/consts"

	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/common"
	"github.com/yuansuan/ticp/iPaaS/sso/hydra_lcp/util"

	"github.com/yuansuan/ticp/common/go-kit/logging"
)

// GetLoginKeyResp ...
type GetLoginKeyResp struct {
	Key   string `json:"key"`
	KeyId string `json:"key_id"`
}

// GetLoginKey GetLoginKey
// swagger:route GET /api/login/getkey   GetLoginKeyResp
//
//	    Responses:
//		     200
//	      90042: ErrHydraLcpGetLoginKeyFail
//
// @GET /api/login/getkey
func (h *Handler) GetLoginKey(c *gin.Context) {
	logger := logging.GetLogger(c)

	defaultCache := boot.MW.DefaultCache()

	// 生成ID
	id, err := h.Idgen.GenerateID(c, &idgen.GenRequest{})
	if err != nil {
		logger.Error("failed to generate id")
		http.Err(c, consts.ErrHydraLcpGetLoginKeyFail, "failed to generate Key ID")
		return
	}

	keyID := snowflake.ParseInt64(id.Id).String()

	key := util.RandomString(24)

	err = defaultCache.PutWithExpire(common.HydraLcpLoginKey, keyID, key, time.Second*20)
	if err != nil {
		logger.Error("failed to storage login key")
		http.Err(c, consts.ErrHydraLcpGetLoginKeyFail, "failed to storage key")
		return
	}

	http.Ok(c, GetLoginKeyResp{Key: key, KeyId: keyID})
	return
}

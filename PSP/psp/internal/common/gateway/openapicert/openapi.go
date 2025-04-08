package openapicert

import (
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/yuansuan/ticp/common/go-kit/logging"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/internal/common/gateway/service"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/ginutil"
)

func CertCheck(ctx *gin.Context) bool {
	logger := logging.GetLogger(ctx)
	url := ctx.Request.URL.Path
	logger.Infof("start basic auth, url:%s", url)

	cert := ctx.Request.Header.Get(common.HttpHeaderOpenapiCertificate)

	// 非openapi不允许调用
	if !strings.HasPrefix(url, common.HttpOpenapiUrl) {
		return false
	}

	user, exist, err := service.RouteSrv.UserService.CheckOpenapiCertificate(ctx, cert)
	if err != nil {
		logger.Errorf("check openapi certificate err:[%v]", err)
		return false
	}

	if !exist {
		return false
	}

	ginutil.SetUser(ctx, user.Id, user.Name)
	v4, _ := uuid.NewV4()
	ginutil.SetTraceID(ctx, v4.String())
	return true
}

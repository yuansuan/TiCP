package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
	iamclient "github.com/yuansuan/ticp/common/project-root-iam/iam-client"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/response"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/errorcode"
)

func UrlValidator(c *gin.Context) {
	logger := trace.GetLogger(c)
	conf, err := getConfig(c)
	if err = response.InternalErrorIfError(c, err, errorcode.InternalServerError); err != nil {
		logger.Errorf("get config from gin ctx failed, %v", err)
		return
	}

	iamclient.ValidUserIDMiddleware(conf.Iam.YsID)(c)
}

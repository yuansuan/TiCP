package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging/middleware"
)

var mids []gin.HandlerFunc

func init() {
	mids = append(mids, ingressLogger())
}

func Middlewares() []gin.HandlerFunc {
	return mids
}

func ingressLogger() gin.HandlerFunc {
	return middleware.IngressLogger(middleware.IngressLoggerConfig{
		IsLogRequestHeader: true,
		IsLogRequestBody:   true,

		IsLogResponseHeader: true,
		IsLogResponseBody:   true,
	})
}

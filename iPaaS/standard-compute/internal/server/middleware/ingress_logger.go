package middleware

import (
	"github.com/gin-gonic/gin"
	logmiddleware "github.com/yuansuan/ticp/common/go-kit/logging/middleware"
)

func IngressLogger(c *gin.Context) {
	logmiddleware.IngressLogger(logmiddleware.IngressLoggerConfig{
		IsLogRequestHeader:  true,
		IsLogRequestBody:    true,
		IsLogResponseHeader: true,
		IsLogResponseBody:   true,
	})(c)
}

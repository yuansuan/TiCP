package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
)

// RequestIDMiddleware add request id to context and response header
func RequestIDMiddleware(c *gin.Context) {
	requestID := trace.GetRequestId(c)
	c.Set(common.RequestIDKey, requestID)
	// set request id to logger
	c.Set(logging.LoggerName, logging.Default().With(common.RequestIDKey, requestID))
	c.Writer.Header().Set(common.RequestIDKey, requestID)
	c.Next()
}

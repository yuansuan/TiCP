package middleware

// RequestIDMiddleware add request id to context and response header
// func RequestIDMiddleware(c *gin.Context) {
// 	requestID := uuid.New().String()
// 	c.Set(common.RequestIDKey, requestID)
// 	// set request id to logger
// 	c.Set(logging.LoggerName, logging.Default().With(common.RequestIDKey, requestID))
// 	c.Writer.Header().Set(common.RequestIDKey, requestID)
// 	c.Next()
// }

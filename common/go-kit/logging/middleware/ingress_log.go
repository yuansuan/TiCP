package middleware

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
)

const (
	ClientIpKey         = "client-ip"
	LatencyKey          = "latency"
	ResponseBodySizeKey = "response-body-size"
	StatusKey           = "status"
)

type IngressLoggerConfig struct {
	IsLogRequestHeader bool
	IsLogRequestBody   bool

	IsLogResponseHeader bool
	IsLogResponseBody   bool
}

func IngressLogger(config IngressLoggerConfig) func(c *gin.Context) {
	return func(c *gin.Context) {
		requestId := trace.GetRequestId(c)
		logger := trace.GetLogger(c).
			Base().
			With(zap.String(trace.PathKey, c.Request.URL.String()))
		startTime := time.Now()

		requestOpts := make([]interface{}, 0)
		if config.IsLogRequestHeader {
			requestOpts = append(requestOpts, zap.Any(trace.RequestHeaderKey, c.Request.Header))
		}
		if config.IsLogRequestBody {
			// body类型若为Nobody，经过下面处理会改变类型，从而影响到hash计算，故body类型不为Nobody时才执行打印。
			if reflect.TypeOf(c.Request.Body) != reflect.TypeOf(http.NoBody) {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(c.Request.Body)
				if err != nil {
					logger.Warnf("read from request body failed")
					c.AbortWithStatusJSON(http.StatusInternalServerError, response{
						ErrorCode: "InternalServerError",
						RequestID: requestId,
						ErrorMsg:  "InternalServerError",
					})
					return
				}

				c.Request.Body = io.NopCloser(buf)
				requestOpts = append(requestOpts, zap.Any(trace.RequestBodyKey, buf.String()))
			}
		}
		requestOpts = append(requestOpts, zap.Any(ClientIpKey, c.ClientIP()))

		logger.With(requestOpts...).Info("request infos")

		if config.IsLogResponseBody {
			c.Writer = &responseBodyWriterHijacker{
				ResponseWriter: c.Writer,
				bodyCache:      new(bytes.Buffer),
			}
		}
		c.Next()

		responseOpts := make([]interface{}, 0)
		if config.IsLogResponseHeader {
			responseOpts = append(responseOpts, zap.Any(trace.ResponseHeaderKey, c.Writer.Header()))
		}
		if config.IsLogResponseBody {
			responseOpts = append(responseOpts, zap.Any(trace.ResponseBodyKey, c.Writer.(*responseBodyWriterHijacker).bodyCache.String()))
		}
		responseOpts = append(responseOpts, zap.Any(LatencyKey, time.Now().Sub(startTime).String()))
		responseOpts = append(responseOpts, zap.Any(ResponseBodySizeKey, c.Writer.Size()))
		responseOpts = append(responseOpts, zap.Any(StatusKey, c.Writer.Status()))

		logger.With(responseOpts...).Info("response infos")
	}
}

type responseBodyWriterHijacker struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

func (w *responseBodyWriterHijacker) Write(d []byte) (int, error) {
	w.bodyCache.Write(d)
	return w.ResponseWriter.Write(d)
}

type response struct {
	ErrorCode string `json:"ErrorCode"`
	ErrorMsg  string `json:"ErrorMsg"`
	RequestID string `json:"RequestID"`
}

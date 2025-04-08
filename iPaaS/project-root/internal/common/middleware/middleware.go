package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	boot "github.com/yuansuan/ticp/common/go-kit/gin-boot"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/monitor"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"go.uber.org/zap/zapcore"
	"io"
	"net"
	"strings"
	"time"
)

const (
	userIDKey                    = "x-ys-user-id"
	metricRequestBodyBytesTotal  = "request_body_bytes_total"
	metricResponseBodyBytesTotal = "response_body_bytes_total"
	labelAppName                 = "app_name"
	labelPath                    = "path"
	labelYSID                    = "ys_id"
	bodyUnit                     = 1024 * 1024 // MB为单位
)

// 自定义requestBody的Decorator，用于在真正写入data之前把data长度写入prometheus
type requestBodyDecorator struct {
	readCloser  io.ReadCloser
	totalLength int
	appName     string
	userID      string
	path        string
	labels      []*monitor.Label
}

// 自定义responseWriter的Decorator，用于在真正写入data之前把data长度写入prometheus
type responseWriterDecorator struct {
	gin.ResponseWriter
	totalLength int
	appName     string
	userID      string
	path        string
	labels      []*monitor.Label
}

// HttpStatMiddleware 用于记录请求中request/response中body长度统计(bytes)
func HttpStatMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 没有ys_id不需要进行统计
		userID := c.Request.Header.Get(userIDKey)
		if userID == "" {
			logging.Default().Warnf("ysid can't be empty")
			c.Next()
			return
		}

		appName := config.Conf.App.Name
		path := c.Request.URL.Path
		labels := []*monitor.Label{
			{
				Name:  labelAppName,
				Value: appName,
			},
			{
				Name:  labelYSID,
				Value: userID,
			},
			{
				Name:  labelPath,
				Value: path,
			},
		}

		// 设置RequestBody的装饰器
		// 在调用Request.Body.Read()时会先调用装饰器的Read()方法，写入data长度到Prometheus
		c.Request.Body = &requestBodyDecorator{
			readCloser: c.Request.Body,
			appName:    appName,
			userID:     userID,
			path:       path,
			labels:     labels,
		}

		// 设置responseWriter的装饰器
		// 在调用responseWriter.Write()时会先调用装饰器的Write()方法，写入data长度到Prometheus
		c.Writer = &responseWriterDecorator{
			ResponseWriter: c.Writer,
			appName:        appName,
			userID:         userID,
			path:           path,
			labels:         labels,
		}

		c.Next()
	}
}

func (r *requestBodyDecorator) Read(data []byte) (n int, err error) {
	startTime := time.Now()
	bytesLen, err := r.readCloser.Read(data)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			logging.Default().Errorf("read request body err, err: %v", err)
		}
		return bytesLen, err
	}

	r.totalLength += bytesLen
	addMetricErr := boot.Monitor.Add(metricRequestBodyBytesTotal, float64(bytesLen)/bodyUnit, r.labels)
	if addMetricErr != nil {
		logging.Default().Errorf("add request_bytes_total metric err, app name: %s, userID: %s, "+
			"path: %s, err: %v", r.appName, r.userID, r.path, err)
	}

	if isDebugLevel() {
		ts := time.Now()
		formattedTime := ts.Format(time.RFC3339)
		elapsedTime := time.Since(startTime)
		logging.Default().Debugf("stat metric request_bytes_total, elapsedTime: %d ms, ts: %s, app name: %s, "+
			"userID: %s, path: %s, request body length: %dB, total length: %dB",
			elapsedTime.Milliseconds(), formattedTime, r.appName, r.userID, r.path, bytesLen, r.totalLength)
	}
	return bytesLen, err
}

func (r *requestBodyDecorator) Close() error {
	return r.readCloser.Close()
}

func (w *responseWriterDecorator) Write(data []byte) (int, error) {
	startTime := time.Now()
	bytesLen, err := w.ResponseWriter.Write(data)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			logging.Default().Warnf("write response data timeout, err: %v", err)
		} else if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
			logging.Default().Warnf("write response data broken pipe, err: %v", err)
		} else {
			logging.Default().Errorf("write response data err, err: %v", err)
		}
		return bytesLen, err
	}

	w.totalLength += bytesLen
	addMetricErr := boot.Monitor.Add(metricResponseBodyBytesTotal, float64(len(data))/bodyUnit, w.labels)
	if addMetricErr != nil {
		logging.Default().Errorf("add response_bytes_total metric err, app name: %s, userID: %s, "+
			"path: %s, err: %v", w.appName, w.userID, w.path, err)
	}

	if isDebugLevel() {
		ts := time.Now()
		formattedTime := ts.Format(time.RFC3339)
		elapsedTime := time.Since(startTime)
		logging.Default().Debugf("stat metric response_bytes_total, elapsedTime: %d ms, ts: %s, "+
			"app name: %s, userID: %s, path: %s, request body length: %dB, total length: %dB",
			elapsedTime.Milliseconds(), formattedTime, w.appName, w.userID, w.path, bytesLen, w.totalLength)
	}
	return bytesLen, err
}

func isDebugLevel() bool {
	return env.Env.LogLevel == int(zapcore.DebugLevel)
}

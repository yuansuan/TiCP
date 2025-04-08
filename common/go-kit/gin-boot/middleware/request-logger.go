/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/yuansuan/ticp/common/go-kit/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/config"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/env"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/http"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

// TraceID TraceID
const TraceID = "trace_id"

// GRPCPrintReqResp GRPCPrintReqResp
func GRPCPrintReqResp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)

	logger := logging.GetLogger(ctx)

	if env.Env.LogLevel == env.LevelDebug && needRequestInfo(info) {
		sReq := MarshalPB(req, false)
		sResp := MarshalPB(resp, true)

		logger.Debugf("%s | req: %s | resp: %v | err: %+v", info.FullMethod, sReq, sResp, err)
	}
	return
}

// GRPCLogger is grpc interceptor for inject logger, log request info
func GRPCLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestid, err := util.GetInMetadata(ctx, TraceID)
	if err != nil {
		requestid = util.RandomString(30)
		ctx = metadata.AppendToOutgoingContext(ctx, TraceID, requestid)
	}

	logger := logging.GetLogger(ctx).With(TraceID, requestid)
	ctx = context.WithValue(ctx, logging.LoggerName, logger)

	_start := time.Now()
	_peer, _peerok := peer.FromContext(ctx)

	// real invoke
	result, err := handler(ctx, req)

	if !needRequestInfo(info) {
		return result, err
	}

	_end := time.Now()
	_latency := _end.Sub(_start)
	_address := "unknown"
	if _peerok {
		_address = _peer.Addr.String()
	}

	_status := status.Convert(err)
	_status.Code()

	appName := "unknown"
	if config.Conf != nil {
		appName = config.Conf.App.Name
	}
	logger.Infof("[GRPC][%s] %v | %v | %13v | %15v | %v %v",
		appName,
		_end.Format("2006/01/02 - 15:04:05"),
		_status.Code(),
		_latency,
		_address,
		info.FullMethod,
		_status.Message())

	return result, err
}

// GinLogger is gin interceptor for inject logger, log request info
func GinLogger(notlogged ...string) http.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		requestid := util.RandomString(30)
		logger := logging.GetLogger(c).With(TraceID, requestid)
		c.Set(logging.LoggerName, logger)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor, resetColor string
			if logging.IsTerminal() {
				statusColor = logging.ColorForStatus(statusCode)
				methodColor = logging.ColorForMethod(method)
				resetColor = logging.ResetColor
			}
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			logger.Infof("[GIN][%s] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				config.Conf.App.Name,
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment,
			)
		}
	}
}

func needRequestInfo(info *grpc.UnaryServerInfo) bool {
	return !strings.HasPrefix(info.FullMethod, "/hpcmonitor.")
}

// MarshalPB : marshal iface(proto.Message) to string
// if some field is too long, will replace with XXXXXXX...(len=123)
// if is not proto.Message, return "UNKNOWN"
func MarshalPB(iface interface{}, skip bool) string {
	result := "UNKNOWN"
	marshaler := jsonpb.Marshaler{}
	pbMessage, ok := iface.(proto.Message)
	if ok {
		result, _ = marshaler.MarshalToString(pbMessage)
		if skip {
			result = string(jsonFilter._json([]byte(result)))
		} else {
			result = string(jsonFilter._jsonNoSkip([]byte(result)))
		}
	}

	return result
}

const messageFieldSkipLen = 5120

var messageTransValueByKey = map[string]string{
	"password": "******",
}

type tJSONFilter struct{}

var jsonFilter tJSONFilter

func (tJSONFilter) skipTooLong(input string) string {
	if len(input) > messageFieldSkipLen {
		return fmt.Sprintf("%s...(len=%d, b64~=%d)", input[:messageFieldSkipLen], len(input), len(input)/4*3)
	}
	return input
}

// if no key, leave empty
func (t tJSONFilter) filter(key string, v interface{}) interface{} {
	switch vv := v.(type) {
	case string:
		if key != "" && messageTransValueByKey[key] != "" {
			return messageTransValueByKey[key]
		}
		return t.skipTooLong(vv)
	}
	return v
}

func (t tJSONFilter) _slice(s []interface{}) {
	for i, v := range s {
		s[i] = jsonFilter._field("", v)
	}
}

func (t tJSONFilter) _map(m map[string]interface{}) {
	for k, v := range m {
		m[k] = jsonFilter._field(k, v)
	}
}

// if no key, leave empty
func (t tJSONFilter) _field(key string, v interface{}) interface{} {
	switch vv := v.(type) {
	case map[string]interface{}:
		jsonFilter._map(vv)
	case []interface{}:
		jsonFilter._slice(vv)
	}
	return jsonFilter.filter(key, v)
}

func (t tJSONFilter) _json(bs []byte) []byte {
	place := map[string]interface{}{}
	_ = json.Unmarshal(bs, &place)

	jsonFilter._field("", place)

	bs, _ = json.Marshal(place)

	return []byte(t.skipTooLong(string(bs)))
}

func (t tJSONFilter) _jsonNoSkip(bs []byte) []byte {
	place := map[string]interface{}{}
	_ = json.Unmarshal(bs, &place)

	jsonFilter._field("", place)

	bs, _ = json.Marshal(place)

	return bs
}

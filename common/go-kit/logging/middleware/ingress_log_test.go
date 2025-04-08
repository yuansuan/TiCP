package middleware

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/yuansuan/ticp/common/go-kit/logging"
	"github.com/yuansuan/ticp/common/go-kit/logging/trace"
)

type Req struct {
	A int `json:"a,omitempty"`
}

func TestIngressLogger(t *testing.T) {
	engine := gin.Default()

	engine.Use(IngressLogger(IngressLoggerConfig{
		IsLogRequestHeader:  true,
		IsLogRequestBody:    true,
		IsLogResponseHeader: true,
		IsLogResponseBody:   true,
	}))

	engine.POST("/a/:id", func(context *gin.Context) {
		logger := logging.GetLogger(context)
		logger.Info("custom")

		traceLogger := trace.GetLogger(context)
		traceLogger.Info("trace")

		context.Header("response-header-key1", "response-header-value1")
		context.JSON(http.StatusOK, response{
			ErrorCode: "errorCode",
			ErrorMsg:  "errorMsg",
			RequestID: "requestId",
		})
	})

	listener, port, err := getFreeListener()
	assert.NoError(t, err)
	defer listener.Close()
	t.Logf("port: %d", port)

	go func() {
		s := &http.Server{
			Handler: engine,
		}

		_ = s.Serve(listener)
	}()

	// wait server started
	time.Sleep(1 * time.Second)

	hc := resty.New()
	_, err = hc.R().
		SetHeader(trace.RequestIdKey, "requestId1").
		SetQueryParam("queryKey1", "queryValue1").
		SetBody(&Req{
			A: 1,
		}).
		Post(fmt.Sprintf("http://127.0.0.1:%d/a/1234", port))
	assert.NoError(t, err)
}

func getFreeListener() (net.Listener, int, error) {
	// Listen on a random port with ":0"
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, 0, err
	}

	// Get the port from the listener's address
	addr := listener.Addr().(*net.TCPAddr)
	return listener, addr.Port, nil
}

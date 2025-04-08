package mock

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

// GinContext mock gin context
// use httptest.NewRecorder() to get mock w
func GinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// HTTPRequest mock http request
// 模拟各种http请求的上下文，用于测试http接口接到不同参数的请求时的处理逻辑
// Example usage:
// HTTPRequest(c, "GET", nil, params, u)
// HTTPRequest(c, "POST", content, params, nil)
// HTTPRequest(c, "PUT", content, params, nil)
// HTTPRequest(c, "PATCH", content, params, nil)
// HTTPRequest(c, "DELETE", content, params, nil)
func HTTPRequest(c *gin.Context, method string, content interface{}, params gin.Params, u url.Values) {
	c.Request.Method = method
	if method != http.MethodGet { //? maybe not need
		c.Request.Header.Set("Content-Type", "application/json")
	}

	// set path params like /api/v1/example/:id
	if params != nil {
		c.Params = params
	}

	// set query params like /api/v1/example?name=xx
	if u != nil {
		c.Request.URL.RawQuery = u.Encode()
	}

	// marshal content to json
	if content != nil {
		jsonbytes, err := json.Marshal(content)
		if err != nil {
			panic(err)
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
)

func ParseResp(resp *http.Response, data interface{}) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		err := json.Unmarshal(body, data)
		if err != nil {
			return err
		}
		return errors.Errorf("http: %v, body: %v", resp.Status, string(body))
	}
	return json.NewDecoder(resp.Body).Decode(data)
}

// GetTestGinContext returns a gin context for testing
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)

	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func BindRequest(ctx *gin.Context, req interface{}) {

	// 获取原始查询参数
	originalQuery := ctx.Request.URL.Query()

	// 使用反射获取结构体字段的名称和值，并将其设置为查询参数
	valueOfReq := reflect.ValueOf(req)
	typeOfReq := reflect.TypeOf(req)

	for i := 0; i < typeOfReq.NumField(); i++ {
		field := typeOfReq.Field(i)
		fieldValue := valueOfReq.Field(i)

		// 跳过未导出字段或空值字段
		if field.PkgPath != "" || fieldValue.IsZero() {
			continue
		}

		// 将字段名称和值转换为字符串，并设置为查询参数
		paramName := field.Name
		paramValue := fmt.Sprintf("%v", fieldValue.Interface())
		originalQuery.Set(paramName, paramValue)
	}

	// 生成新的 URL 对象并替换上下文中的 URL
	ctx.Request.URL.RawQuery = originalQuery.Encode()
}

func BindJsonRequest(ctx *gin.Context, req interface{}) error {
	bodyData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyData))
	ctx.Request.ContentLength = int64(len(bodyData))
	ctx.Request.Header.Set("Content-Type", "application/json")
	return nil

}

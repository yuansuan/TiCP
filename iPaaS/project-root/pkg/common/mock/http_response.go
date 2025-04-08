package mock

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

// HTTPResponse mock http response
// 模拟第三方接口的返回值, 用于测试代码内部调用第三方接口时的逻辑
func HTTPResponse(client *http.Client, url, method string, statusCode int, headers map[string][]string, responseBody string, cookies []*http.Cookie) func() {
	// Activate HTTP mock for the client
	httpmock.ActivateNonDefault(client)

	// Register a mock response for the URL
	mockResponse := httpmock.NewStringResponse(statusCode, responseBody)
	mockResponse.Status = fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	mockResponse.Proto = "HTTP/1.1"
	mockResponse.Header = headers

	// Set cookies in the response
	for _, cookie := range cookies {
		mockResponse.Header.Add("Set-Cookie", cookie.String())
	}

	httpmock.RegisterResponder(method, url, httpmock.ResponderFromResponse(mockResponse))

	// Return a cleanup function
	return func() {
		// Deactivate HTTP mock for the client and reset
		httpmock.DeactivateAndReset()
	}
}

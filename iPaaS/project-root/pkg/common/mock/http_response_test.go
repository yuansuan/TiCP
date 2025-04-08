package mock

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/go-resty/resty/v2"
)

// 参考示例

const testUrl = "https://mytest.com"

// httpResponseExample 示例
type httpResponseExample struct {
	hc   *resty.Client
	resp *resty.Response
}

// sendRequest 示例
func (e *httpResponseExample) sendRequest() {
	resp, err := e.hc.R().Get(testUrl)
	if err != nil {
		log.Fatal(err)
	}

	e.resp = resp
}

func (e *httpResponseExample) printResponse() {
	if e.resp == nil {
		log.Fatal("Response is nil")
	}

	fmt.Println("Response Info:")
	fmt.Println("Status Code:", e.resp.StatusCode())
	fmt.Println("Status:", e.resp.Status())
	fmt.Println("Proto:", e.resp.Proto())
	fmt.Println("Time:", e.resp.Time())
	fmt.Println("Received At:", e.resp.ReceivedAt())
	fmt.Println("Size:", e.resp.Size())
	fmt.Println("Headers:")
	for key, value := range e.resp.Header() {
		fmt.Println(key, "=", value)
	}
	fmt.Println("Cookies:")
	for i, cookie := range e.resp.Cookies() {
		fmt.Printf("cookie%d: name:%s value:%s\n", i, cookie.Name, cookie.Value)
	}
}

func newCustomCookieJar(cookies []*http.Cookie) http.CookieJar {
	jar, _ := cookiejar.New(nil)
	jar.SetCookies(&url.URL{Scheme: "http", Host: "example.com"}, cookies)
	return jar
}

// TestSendRequest 示例
func TestSendRequest(t *testing.T) {

	// Create a new Resty client
	client := resty.New()
	cookies := []*http.Cookie{
		{Name: "cookie1", Value: "value1"},
		{Name: "cookie2", Value: "value2"},
	}
	client.SetCookieJar(newCustomCookieJar(cookies))

	// Create YourStruct instance with the Resty client
	e := httpResponseExample{hc: client}

	// Mock HTTP setup and cleanup
	defer HTTPResponse(client.GetClient(), testUrl, "GET", 200, map[string][]string{"Content-Type": {"application/json"}}, "Mock response", cookies)()

	// Call the method under test
	e.sendRequest()

	// Print response info
	e.printResponse()
}

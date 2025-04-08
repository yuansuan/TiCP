package openapi_go

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/yuansuan/ticp/common/openapi-go/apiv1"
	"github.com/yuansuan/ticp/common/openapi-go/credential"
	"github.com/yuansuan/ticp/common/openapi-go/utils/signer"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xhttp"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xtime"
	"github.com/yuansuan/ticp/common/openapi-go/utils/xurl"
)

const (
	// header
	V1Header = "X-Ys-Version"
	V1Value  = "2023-05-30"
	// DefaultBaseURL defaults gateway for Open-API
	DefaultBaseURL = "http://localhost:8899/"
)

// Client represents the Open-API client
type Client struct {
	*apiv1.API

	hc     *xhttp.Client
	cred   *credential.Credential
	signer *signer.Signer

	baseURL        string
	timeout        time.Duration
	proxy          string
	retryTimes     int
	retryInterval  time.Duration
	retryCondition func(*http.Response, error) bool
}

func (c *Client) interceptor(req *http.Request, next xhttp.MiddlewareHandler) (*http.Response, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	xurl.Merge(req.URL, u)
	xhttp.AppendQuery(req, "AppKey", c.cred.GetAppKey())
	xhttp.AppendQuery(req, "Timestamp", xtime.CurrentTimestamp())
	xhttp.AddHeader(req, V1Header, V1Value)
	f, err := c.signer.SignHttp(req)
	if err != nil {
		return nil, err
	}

	return next(f.AsQuery(req))
}

type Option func(c *Client) error

func WithBaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.timeout = timeout
		return nil
	}
}

func WithProxy(proxy string) Option {
	return func(c *Client) error {
		c.proxy = proxy
		return nil
	}
}

func WithRetryTimes(retryTimes int) Option {
	return func(c *Client) error {
		if retryTimes <= 0 {
			return fmt.Errorf("retry times cannot less or equal than 0")
		}

		c.retryTimes = retryTimes
		return nil
	}
}

func WithRetryInterval(retryInterval time.Duration) Option {
	return func(c *Client) error {
		if retryInterval <= 0 {
			return fmt.Errorf("retry interval cannot less or equal than 0")
		}

		c.retryInterval = retryInterval
		return nil
	}
}

func WithRetryCondition(retryCondition func(resp *http.Response, err error) bool) Option {
	return func(c *Client) error {
		if retryCondition == nil {
			return fmt.Errorf("retry condition cannot be nil")
		}
		c.retryCondition = retryCondition
		return nil
	}
}

// NewClient creates new client for Open-YS
func NewClient(cred *credential.Credential, options ...Option) (c *Client, err error) {
	c = &Client{
		cred:           cred,
		baseURL:        DefaultBaseURL,
		retryInterval:  xhttp.DefaultRetryInterval,
		retryCondition: xhttp.DefaultRetryCondition,
	}
	for _, option := range options {
		if err = option(c); err != nil {
			return nil, err
		}
	}

	if c.signer, err = signer.NewSigner(c.cred); err != nil {
		return nil, err
	}

	var httpClientOpts []xhttp.ClientOption
	httpClientOpts = append(httpClientOpts,
		xhttp.WithMiddleware(xhttp.MiddlewareFunc(c.interceptor)),
		xhttp.WithTimeout(c.timeout),
		xhttp.WithProxy(c.proxy),
		xhttp.WithRetryTimes(c.retryTimes),
		xhttp.WithRetryInterval(c.retryInterval),
		xhttp.WithRetryCondition(c.retryCondition),
	)

	if c.hc, err = xhttp.NewClient(httpClientOpts...); err != nil {
		return nil, err
	}

	if c.API, err = apiv1.NewAPI(c.hc); err != nil {
		return nil, err
	}

	return
}

func (c *Client) SetBaseUrl(baseURL string) {
	c.baseURL = baseURL
}

func (c *Client) GetBaseUrl() string {
	return c.baseURL
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.hc.SetTimeout(timeout)
}

func (c *Client) SetProxy(proxy string) error {
	err := c.hc.SetProxy(proxy)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SetCredential(cred *credential.Credential) error {
	c.cred = cred

	sign, err := signer.NewSigner(c.cred)
	if err != nil {
		return fmt.Errorf("new signer failed, %w", err)
	}
	c.signer = sign

	return nil
}

func (c *Client) PrintSign(rb *xhttp.RequestBuilder) {
	req, err := rb.Build()
	if err != nil {
		fmt.Printf("build request exception: %v", err)
		return
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		fmt.Printf("parse url exception: %v", err)
		return
	}

	xurl.Merge(req.URL, u)
	xhttp.AppendQuery(req, "AppKey", c.cred.GetAppKey())
	ts := xtime.CurrentTimestamp()
	xhttp.AppendQuery(req, "Timestamp", ts)

	xhttp.AddHeader(req, V1Header, V1Value)

	f, err := c.signer.SignHttp(req)
	if err != nil {
		fmt.Printf("sign exception: %v", err)
		return
	}

	fmt.Printf("?AppKey=%v&Timestamp=%v&Signature=%v\n", f.AppKey, ts, f.Signature)

	fmt.Printf("AppKey: %v\n", f.AppKey)
	fmt.Printf("Timestamp: %v\n", ts)
	fmt.Printf("Signature: %v\n", f.Signature)
	fmt.Println("")
}

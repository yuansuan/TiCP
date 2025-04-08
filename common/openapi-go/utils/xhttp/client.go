package xhttp

import (
	"net/http"
	"time"
)

type Client struct {
	hc *http.Client
	tr RoundTripperWithMiddleware

	retryTimes     int
	retryInterval  time.Duration
	retryCondition func(*http.Response, error) bool
}

func (c *Client) Prepare(rb *RequestBuilder) ResponseFuture {
	return ResponseFutureFunc(func(resolve ResponseResolver) error {
		req, err := rb.Build()
		if err != nil {
			return err
		}

		resp, err := c.doWithRetry(req)
		if err != nil {
			return err
		}

		return resolve(resp)
	})
}

func (c *Client) PrepareDownload(rb *RequestBuilder) ResponseFuture {
	return ResponseFutureFunc(func(resolve ResponseResolver) error {
		req, err := rb.Build()
		if err != nil {
			return err
		}

		resp, err := c.doWithRetry(req)
		if err != nil {
			return err
		}

		return resolve(resp)
	})
}

func (c *Client) Do(rb *RequestBuilder) (*http.Response, error) {
	req, err := rb.Build()
	if err != nil {
		return nil, err
	}

	resp, err := c.hc.Do(req)
	return resp, err
}

func (c *Client) doWithRetry(req *http.Request) (resp *http.Response, err error) {
	for i := 0; ; i++ {
		resp, err = c.hc.Do(req)
		if !c.retryCondition(resp, err) || i >= c.retryTimes {
			return
		}
		<-time.After(c.retryInterval)
	}
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.hc.Timeout = timeout
}

func (c *Client) SetProxy(proxy string) error {
	err := c.tr.SetProxy(proxy)
	if err != nil {
		return err
	}
	return nil
}

type ClientOption func(hc *Client) error

func NewClient(options ...ClientOption) (*Client, error) {
	tr := newTransport()
	c := &Client{
		tr: tr,
		hc: &http.Client{
			Transport: tr,
		},
		retryCondition: DefaultRetryCondition,
		retryInterval:  DefaultRetryInterval,
	}
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type ResponseResolver func(resp *http.Response) error

type ResponseFuture interface {
	Resolve(resolver ResponseResolver) error
}

type ResponseFutureFunc func(resolver ResponseResolver) error

func (f ResponseFutureFunc) Resolve(resolver ResponseResolver) error {
	return f(resolver)
}

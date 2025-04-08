package xhttp

import (
	"net/http"
	"net/url"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xurl"
)

type hijacker struct {
	*http.Transport

	request  []RequestInterceptor
	response []ResponseInterceptor
}

func (h *hijacker) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, interceptor := range h.request {
		if err := interceptor(req); err != nil {
			return nil, err
		}
	}

	resp, err := h.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	for _, interceptor := range h.response {
		if err = interceptor(resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func newHijacker() *hijacker {
	return &hijacker{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment}}
}

type RequestInterceptor = func(req *http.Request) error

func WithRequestInterceptor(ri RequestInterceptor) ClientOption {
	return func(c *Client) error {
		c.hijacker.request = append(c.hijacker.request, ri)
		return nil
	}
}

type ResponseInterceptor = func(resp *http.Response) error

func WithBaseUrl(base string) ClientOption {
	raw, err := url.Parse(base)
	return WithRequestInterceptor(func(req *http.Request) error {
		if err != nil {
			return err
		}

		req.URL = xurl.Merge(req.URL, raw)
		return nil
	})
}

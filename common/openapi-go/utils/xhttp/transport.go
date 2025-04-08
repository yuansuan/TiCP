package xhttp

import (
	"fmt"
	"net/http"
	"net/url"
)

type RoundTripperWithMiddleware interface {
	http.RoundTripper

	AddMiddleware(middleware Middleware, priority ...int)
	SetProxy(proxy string) error
}

type transport struct {
	*http.Transport

	middlewares Middlewares
}

func (tr *transport) AddMiddleware(middleware Middleware, priority ...int) {
	pm := newPriorityMiddleware(middleware, _DefaultMiddlewarePriority)
	if len(priority) != 0 {
		pm.priority = priority[0]
	}

	tr.middlewares = append(tr.middlewares, pm)
}

func (tr *transport) SetProxy(proxy string) error {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return fmt.Errorf("Error parsing proxy URL: %w", err)
	}
	tr.Transport.Proxy = http.ProxyURL(proxyURL)
	return nil
}

func (tr *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	next := tr.Transport.RoundTrip
	for _, middleware := range tr.middlewares {
		next = func(next MiddlewareHandler) MiddlewareHandler {
			return func(req *http.Request) (*http.Response, error) {
				return middleware.Process(req, next)
			}
		}(next)
	}
	return next(req)
}

func newTransport() RoundTripperWithMiddleware {
	return &transport{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}}
}

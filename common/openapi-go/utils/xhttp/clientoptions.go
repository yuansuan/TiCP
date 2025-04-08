package xhttp

import (
	"fmt"
	"net/http"
	"time"
)

func WithMiddleware(middleware Middleware, priority ...int) ClientOption {
	return func(hc *Client) error {
		hc.tr.AddMiddleware(middleware, priority...)
		return nil
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(hc *Client) error {
		hc.SetTimeout(timeout)
		return nil
	}
}

func WithProxy(proxy string) ClientOption {
	return func(hc *Client) error {
		if proxy == "" {
			return nil
		}
		err := hc.SetProxy(proxy)
		if err != nil {
			return err
		}
		return nil
	}
}

func WithRetryTimes(retryTimes int) ClientOption {
	return func(hc *Client) error {
		hc.retryTimes = retryTimes
		return nil
	}
}

func WithRetryInterval(retryInterval time.Duration) ClientOption {
	return func(hc *Client) error {
		if retryInterval == 0 {
			return fmt.Errorf("retryInterval cannot be 0")
		}

		hc.retryInterval = retryInterval
		return nil
	}
}

func WithRetryCondition(retryCondition func(*http.Response, error) bool) ClientOption {
	return func(hc *Client) error {
		if retryCondition != nil {
			hc.retryCondition = retryCondition
		}

		return nil
	}
}

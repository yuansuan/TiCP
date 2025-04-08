package xhttp

import (
	_errors "errors"
	"fmt"
	"net/http"
)

var (
	// ErrHttpNotFound 404: 找不到资源
	ErrHttpNotFound = _errors.New("http: not found")
)

// Client 基于官方 http.Client 封装的Http客户端
type Client struct {
	*http.Client

	hijacker *hijacker
}

// ClientOption 创建客户端的额外选项
type ClientOption func(c *Client) error

// New 创建一个增强版的Http客户端
func New(options ...ClientOption) (*Client, error) {
	c := &Client{Client: &http.Client{}, hijacker: newHijacker()}
	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	c.Client.Transport = c.hijacker

	return c, nil
}

// HttpError 是一个Http错误
type HttpError struct {
	Code    int
	Message string
}

// Error 返回Http的错误消息
func (e *HttpError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

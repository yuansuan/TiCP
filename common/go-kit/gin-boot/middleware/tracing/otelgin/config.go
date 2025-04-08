package otelgin

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type SpanNaming func(c *gin.Context) string

func defaultSpanNaming(c *gin.Context) string {
	spanName := c.FullPath()
	if spanName == "" {
		spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
	}
	return spanName
}

var (
	_DefaultSpanNaming    = defaultSpanNaming
	_DefaultIgnoreMatcher IgnoreMatcher
)

func SetDefaultSpanNaming(sn SpanNaming) {
	_DefaultSpanNaming = sn
}

func SetDefaultIgnoreMatcher(matcher IgnoreMatcher) {
	_DefaultIgnoreMatcher = matcher
}

type Option func(*config)

func newConfig(options ...Option) *config {
	cfg := &config{
		tp: otel.GetTracerProvider(),
		pg: otel.GetTextMapPropagator(),
		sn: _DefaultSpanNaming,

		ignoreMatcher: _DefaultIgnoreMatcher,
	}

	for _, opt := range options {
		opt(cfg)
	}

	return cfg
}

type config struct {
	tp oteltrace.TracerProvider
	pg propagation.TextMapPropagator

	sn            SpanNaming
	excludes      excludes
	ignoreMatcher IgnoreMatcher

	enabled bool
}

func WithExcludes(paths ...string) Option {
	return func(c *config) {
		c.excludes = append(c.excludes, paths...)
	}
}

func WithSpanNaming(sn SpanNaming) Option {
	return func(cfg *config) {
		cfg.sn = sn
	}
}

type IgnoreMatcher func(c *gin.Context) bool

func WithIgnoreMatcher(matcher IgnoreMatcher) Option {
	return func(c *config) {
		c.ignoreMatcher = matcher
	}
}

func WithEnabled(enabled bool) Option {
	return func(c *config) {
		c.enabled = enabled
	}
}

func WithIgnorePaths(paths ...string) Option {
	m := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		m[path] = struct{}{}
	}

	return WithIgnoreMatcher(func(c *gin.Context) bool {
		_, ok := m[c.Request.URL.Path]
		return !ok
	})
}

type excludes []string

func (e excludes) Match(path string) bool {
	if e != nil {
		for _, item := range e {
			if len(item) == len(path) && path == item {
				return true
			}
			if len(item) < len(path) && strings.HasPrefix(path, item) {
				return true
			}
		}
	}

	return false
}

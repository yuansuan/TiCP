package middleware

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type lzHandler struct {
	*Options
}

func NewLzHandler(options ...Option) *lzHandler {
	handler := &lzHandler{
		Options: DefaultOptions,
	}
	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (g *lzHandler) Handle(c *gin.Context) {
	if fn := g.DecompressFn; fn != nil && c.Request.Header.Get("Content-Encoding") == "lz" {
		fn(c)
	}

	if !g.shouldCompress(c.Request) {
		return
	}

	w := newLzWriter(c.Writer)

	c.Header("Content-Encoding", "lz")
	c.Header("Vary", "Accept-Encoding")
	c.Writer = w
	defer func() {
		if err := w.Close(); err != nil {
			c.AbortWithStatus(400)
			return
		}

		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}

func (g *lzHandler) shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "lz") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Accept"), "text/event-stream") {
		return false
	}

	extension := filepath.Ext(req.URL.Path)
	if g.ExcludedExtensions.Contains(extension) {
		return false
	}

	if g.ExcludedPaths.Contains(req.URL.Path) {
		return false
	}
	if g.ExcludedPathesRegexs.Contains(req.URL.Path) {
		return false
	}

	return true
}

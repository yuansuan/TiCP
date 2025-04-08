package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	registerEndpoint(endPoint{
		Method:       http.MethodGet,
		RelativePath: "/ready",
		Handler:      ready,
	})
}

func ready(c *gin.Context) {
	c.Status(http.StatusOK)
}

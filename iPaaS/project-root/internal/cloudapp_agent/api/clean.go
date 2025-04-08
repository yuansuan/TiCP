package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	registerEndpoint(endPoint{
		Method:       http.MethodPost,
		RelativePath: "/clean",
		Handler:      clean,
	})
}

func clean(c *gin.Context) {
	//onclose.Clean()

	c.Status(http.StatusOK)
}

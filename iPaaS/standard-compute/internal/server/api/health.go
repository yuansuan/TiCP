package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, nil)
}

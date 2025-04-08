package api

import (
	"github.com/gin-gonic/gin"
)

type endPoint struct {
	Method       string
	RelativePath string
	Handler      gin.HandlerFunc
}

var Endpoints []endPoint

func registerEndpoint(e endPoint) {
	Endpoints = append(Endpoints, e)
}

func GetEndpoints() []endPoint {
	return Endpoints
}

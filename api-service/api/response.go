package api

import (
	"github.com/gin-gonic/gin"
)

// ResponseGin wrap gin.Context for router/api used
type ResponseGin struct {
	Context *gin.Context
}

// Response struct
type Response struct {
	Data interface{} `json:"data"`
}

// Response function used by ResponseGin
func (g *ResponseGin) Response(httpCode int, data interface{}) {
	g.Context.JSON(httpCode, Response{
		Data: data,
	})

	return
}

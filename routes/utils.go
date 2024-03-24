package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UnImplemented() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{
			"message": "endpoint is not yet implemented",
		})
	}
}

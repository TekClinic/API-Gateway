package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(router *gin.Engine) {
	router.GET("/patient", UnImplemented())
	router.POST("/patient", UnImplemented())
	router.GET("/patient/:id", UnImplemented())
}

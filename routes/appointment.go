package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterAppointmentRoutes(router *gin.Engine) {
	router.GET("/appointment", UnImplemented())
	router.POST("/appointment", UnImplemented())
	router.GET("/appointment/:id", UnImplemented())
	router.PUT("/appointment/:id/patient", UnImplemented())
	router.DELETE("/appointment/:id/patient", UnImplemented())
}

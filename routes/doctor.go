package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterDoctorRoutes(router *gin.Engine) {
	router.GET("/doctor", UnImplemented())
	router.GET("/doctor/:id", UnImplemented())
}

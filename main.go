package main

import (
	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	router := gin.New()

	// enable logging
	router.Use(gin.Logger())
	// recover in case of panic
	router.Use(gin.Recovery())
	// setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization"},
	}))
	// require authorization on all endpoints
	router.Use(middlewares.AuthRequired())

	routes.RegisterPatientRoutes(router)
	routes.RegisterDoctorRoutes(router)
	routes.RegisterAppointmentRoutes(router)

	err := router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal(err)
	}
}

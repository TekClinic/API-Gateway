package main

import (
	"log"

	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/gin-contrib/location"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	envURIScheme = "URI_SCHEME"
	envURIHost   = "URI_HOST"

	defaultURIScheme = "http"
	defaultURIHost   = "localhost"
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
	// setup middleware to discover hostname
	router.Use(location.New(location.Config{
		Scheme: ms.GetOptionalEnv(envURIScheme, defaultURIScheme),
		Host:   ms.GetOptionalEnv(envURIHost, defaultURIHost),
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

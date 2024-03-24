package main

import (
	"github.com/TekClinic/API-Gateway/middlewares"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"

	patients "github.com/TekClinic/Patients-MicroService/patients_protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func fetchPatientData(patientsService *ms.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authToken, err := extractBearerToken(ctx.GetHeader("Authorization"))

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if authToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "bearer token is missing",
			})
			return
		}

		// Authenticate auth_token with auth-microservice
		conn, err := grpc.Dial(patientsService.GetAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "unknown error occurred",
			})
			return
		}
		defer conn.Close()
		client := patients.NewPatientsServiceClient(conn)

		patient, err := client.GetMe(ctx, &patients.MeRequest{Token: authToken})
		if err != nil {
			if status.Code(err) == codes.Unauthenticated {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "invalid authentication token",
				})
			} else if status.Code(err) == codes.PermissionDenied {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "you are not allowed to do this",
				})
			} else if status.Code(err) == codes.NotFound {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "data is missing",
				})
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":   "unknown error occurred",
					"details": err.Error(),
				})
			}
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user_id":  patient.GetUserId(),
			"username": patient.GetName(),
		})
	}
}

func main() {
	patientsService, err := ms.FetchServiceParameters("patients")
	if err != nil {
		log.Fatal(err)
	}

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

	router.GET("/patients/me", fetchPatientData(patientsService))

	err = router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	patientpb "github.com/TekClinic/Patients-MicroService/patients_protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	host string
	port string
}

func (s Service) getAddr() string {
	return s.host + ":" + s.port
}

func getRequiredEnv(key string) (string, error) {
	value, set := os.LookupEnv(key)
	if !set {
		return "", errors.New(key + " environment variable is missing")
	}
	return value, nil
}

func getOptionalEnv(key string, def string) string {
	value, set := os.LookupEnv(key)
	if set {
		return value
	}
	return def
}

func fetchServiceParameters(serviceName string) (*Service, error) {
	host, err := getRequiredEnv(fmt.Sprintf("MS_%s_HOST", strings.ToUpper(serviceName)))
	if err != nil {
		return nil, err
	}

	port := getOptionalEnv(fmt.Sprintf("MS_%s_PORT", strings.ToUpper(serviceName)), "9090")
	return &Service{host: host, port: port}, nil
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 || jwtToken[0] != "Bearer" {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func fetchPatientData(patientsService *Service) gin.HandlerFunc {
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
		conn, err := grpc.Dial(patientsService.getAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "unknown error occurred",
			})
			return
		}
		defer conn.Close()
		client := patientpb.NewPatientsServiceClient(conn)

		patient, err := client.GetMe(ctx, &patientpb.MeRequest{Token: authToken})
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
	router := gin.Default()
	patientsService, err := fetchServiceParameters("patients")
	if err != nil {
		log.Fatal(err)
	}

	// setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Authorization"},
	}))

	router.GET("/patients/me", fetchPatientData(patientsService))

	err = router.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal(err)
	}
}

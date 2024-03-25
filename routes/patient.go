package routes

import (
	"log"
	"net/http"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/schemas"
	ms "github.com/TekClinic/MicroService-Lib"
	patients "github.com/TekClinic/Patients-MicroService/patients_protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const resourceName = "patient"

type PatientsParams struct {
	Skip  int32 `form:"skip,default=0"`
	Limit int32 `form:"limit,default=20"`
}

func getPatients(service patients.PatientsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the query
		var params PatientsParams
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.GetPatientsIds(ctx, &patients.PatientsRequest{
			Token:  ctx.GetString(middlewares.TokenKey),
			Limit:  params.Limit,
			Offset: params.Skip,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			CreateNamedAPIResourceList(ctx, resourceName,
				params.Skip, params.Limit, response.GetCount(), response.GetResults()))
	}
}

func RegisterPatientRoutes(router *gin.Engine) {
	patientsService, err := ms.FetchServiceParameters(resourceName)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial(patientsService.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := patients.NewPatientsServiceClient(conn)

	router.GET("/patient", getPatients(client))
	router.POST("/patient", UnImplemented())
	router.GET("/patient/:id", UnImplemented())
}

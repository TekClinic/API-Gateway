package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/schemas"
	doctors "github.com/TekClinic/Doctors-MicroService/doctors_protobuf"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const d_resourceName = "doctor"

type DoctorsParams struct {
	Skip  int32 `form:"skip,default=0"`
	Limit int32 `form:"limit,default=20"`
}

func getDoctors(service doctors.DoctorServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the query
		var params DoctorsParams
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.GetDoctorsIDs(ctx, &doctors.DoctorsRequest{
			Token:  ctx.GetString(middlewares.TokenKey),
			Limit:  params.Limit,
			Offset: params.Skip,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			CreateNamedAPIResourceList(ctx, d_resourceName,
				params.Skip, params.Limit, response.GetCount(), response.GetResults()))
	}
}

type DoctorParams struct {
	ID int64 `uri:"id" binding:"required"`
}

func getDoctor(service doctors.DoctorServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var params DoctorParams
		err := ctx.ShouldBindUri(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.GetDoctor(ctx, &doctors.DoctorRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			schemas.Doctor{
				DoctorBase: schemas.DoctorBase{
					Name:         response.GetName(),
					Gender:       strings.ToLower(response.GetGender().String()),
					PhoneNumber:  response.GetPhoneNumber(),
					Specialities: response.GetSpecialities(),
					SpecialNote:  response.GetSpecialNote(),
				},
				ID:     response.GetId(),
				Active: response.GetActive(),
			})
	}
}

func RegisterDoctorRoutes(router *gin.Engine) {
	doctorsService, err := ms.FetchServiceParameters(d_resourceName)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial(doctorsService.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := doctors.NewDoctorServiceClient(conn)
	router.GET("/doctor", getDoctors(client))
	router.GET("/doctor/:id", getDoctor(client))
}

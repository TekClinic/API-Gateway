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
)

const resourceNameDoctor = "doctor"

type DoctorsParams struct {
	Skip  int32 `form:"skip,default=0"`
	Limit int32 `form:"limit,default=20"`
}

func getDoctors(service doctors.DoctorsServiceClient) gin.HandlerFunc {
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
		response, err := service.GetDoctorsIDs(ctx, &doctors.GetDoctorsIDsRequest{
			Token:  ctx.GetString(middlewares.TokenKey),
			Limit:  params.Limit,
			Offset: params.Skip,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			CreateNamedAPIResourceList(ctx, resourceNameDoctor,
				params.Skip, params.Limit, response.GetCount(), response.GetResults()))
	}
}

type DoctorParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getDoctor(service doctors.DoctorsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var uriParams DoctorParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.GetDoctor(ctx, &doctors.GetDoctorRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		if response.GetDoctor() == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Message: "Invalid response from the server.",
			})
			return
		}

		doctor := response.GetDoctor()

		specialities := doctor.GetSpecialities()
		if specialities == nil {
			specialities = []string{}
		}

		ctx.JSON(http.StatusOK, schemas.Doctor{
			ID:           doctor.GetId(),
			Active:       doctor.GetActive(),
			Name:         doctor.GetName(),
			Gender:       strings.ToLower(doctor.GetGender().String()),
			PhoneNumber:  doctor.GetPhoneNumber(),
			Specialities: specialities,
			SpecialNote:  doctor.GetSpecialNote(),
		})
	}
}

func RegisterDoctorRoutes(router *gin.Engine) {
	doctorsService, err := ms.FetchServiceParameters(resourceNameDoctor)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.NewClient(doctorsService.GetAddr(), grpc.WithTransportCredentials(GetTransportCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := doctors.NewDoctorsServiceClient(conn)
	router.GET("/doctor", getDoctors(client))
	router.GET("/doctor/:id", getDoctor(client))
}

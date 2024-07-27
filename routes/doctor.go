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
	Skip   int32  `form:"skip,default=0"`
	Limit  int32  `form:"limit,default=20"`
	Search string `form:"search" binding:"omitempty,min=1,max=100"`
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

		// call doctor microservice
		response, err := service.GetDoctorsIDs(ctx, &doctors.GetDoctorsIDsRequest{
			Token:  ctx.GetString(middlewares.TokenKey),
			Limit:  params.Limit,
			Offset: params.Skip,
			Search: params.Search,
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

		// call doctor microservice
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
			DoctorBase: schemas.DoctorBase{
				Name:         doctor.GetName(),
				Gender:       strings.ToLower(doctor.GetGender().String()),
				PhoneNumber:  doctor.GetPhoneNumber(),
				Specialities: specialities,
				SpecialNote:  doctor.GetSpecialNote(),
			},
			ID:     doctor.GetId(),
			Active: doctor.GetActive(),
		})
	}
}

func createDoctor(service doctors.DoctorsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the body
		var bodyParams schemas.DoctorBase
		err := ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call doctor microservice
		response, err := service.CreateDoctor(ctx, &doctors.CreateDoctorRequest{
			Token:        ctx.GetString(middlewares.TokenKey),
			Name:         bodyParams.Name,
			Gender:       doctors.Doctor_Gender(doctors.Doctor_Gender_value[strings.ToUpper(bodyParams.Gender)]),
			PhoneNumber:  bodyParams.PhoneNumber,
			Specialities: bodyParams.Specialities,
			SpecialNote:  bodyParams.SpecialNote,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusCreated, schemas.IDHolder{
			ID: response.GetId(),
		})
	}
}

func deleteDoctor(service doctors.DoctorsServiceClient) gin.HandlerFunc {
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

		// call doctor microservice
		_, err = service.DeleteDoctor(ctx, &doctors.DeleteDoctorRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.Status(http.StatusOK)
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
	router.POST("/doctor", createDoctor(client))
	router.GET("/doctor/:id", getDoctor(client))
	router.DELETE("/doctor/:id", deleteDoctor(client))
}

package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/schemas"
	appointments "github.com/TekClinic/Appointments-MicroService/appointments_protobuf"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/gin-gonic/gin"
	sf "github.com/sa-/slicefunk"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AppointmentIDHolder implements AppointmentIDHolder schema.
type AppointmentIDHolder struct {
	ID int32 `json:"id"`
}

// PatientIDHolder implements PatientIDHolder schema.
type PatientIDHolder struct {
	ID int32 `json:"id"`
}

// AssignPatientIDHolder PatientIDHolder implements PatientIDHolder schema.
type AssignPatientIDHolder struct {
	ID int32 `json:"patient_id"`
}

// DeletedMessageHolder implements DeletedMessageHolder schema.
type DeletedMessageHolder struct {
	Message string `json:"message"`
}

// CreateAppointmentAPIResourceList creates AppointmentAPIResourceList for the given request.
func CreateAppointmentAPIResourceList(ctx *gin.Context, resourceName string,
	count int32, ids []int32) schemas.NamedAPIResourceList {
	var previous, next *string
	previousString := "previous"
	previous = &previousString
	nextString := "next"
	next = &nextString
	return schemas.NamedAPIResourceList{
		Count:    count,
		Next:     next,
		Previous: previous,
		Results: sf.Map(ids, func(id int32) schemas.NamedAPIResource {
			return CreateAppointmentAPIResource(ctx, resourceName, id)
		}),
	}
}

// CreateAppointmentAPIResource creates AppointmentAPIResource for resourceName with given id.
func CreateAppointmentAPIResource(ctx *gin.Context, resourceName string, id int32) schemas.NamedAPIResource {
	requestURL := retrieveRequestURL(ctx)
	requestURL.RawQuery = ""
	requestURL.Path = fmt.Sprintf("/%s/%d", resourceName, id)
	return schemas.NamedAPIResource{
		Name: resourceName,
		URL:  requestURL.String(),
	}
}

const resourceNameAppointment = "appointment"

type AppointmentsParams struct {
	Date      string `form:"date"`
	DoctorID  int32  `form:"doctor_id"`
	PatientID int32  `form:"patient_id"`
	Skip      int32  `form:"skip,default=0"`
	Limit     int32  `form:"limit,default=20"`
}

func getAppointments(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params AppointmentsParams
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// cal appointment microservice
		response, err := service.GetAppointments(ctx, &appointments.GetAppointmentsRequest{
			Token:     ctx.GetString(middlewares.TokenKey),
			Skip:      params.Skip,
			Limit:     params.Limit,
			Date:      params.Date,
			DoctorId:  params.DoctorID,
			PatientId: params.PatientID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			CreateAppointmentAPIResourceList(ctx, resourceNameAppointment,
				response.GetCount(), response.GetResults()))
	}
}

type AppointmentParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params AppointmentParams
		err := ctx.ShouldBindUri(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.GetAppointment(ctx, &appointments.GetAppointmentRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			schemas.Appointment{
				ID:                response.GetId(),
				PatientID:         response.GetPatientId(),
				DoctorID:          response.GetDoctorId(),
				StartTime:         response.GetStartTime(),
				EndTime:           response.GetEndTime(),
				ApprovedByPatient: response.GetApprovedByPatient(),
				Visited:           response.GetVisited(),
			})
	}
}

func createAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params schemas.AppointmentBase
		err := ctx.ShouldBindJSON(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.CreateAppointment(ctx, &appointments.CreateAppointmentRequest{
			Token:     ctx.GetString(middlewares.TokenKey),
			PatientId: params.PatientID,
			DoctorId:  params.DoctorID,
			StartTime: params.StartTime,
			EndTime:   params.EndTime,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusCreated,
			AppointmentIDHolder{
				ID: response.GetId(),
			})
	}
}

type AssignPatientParams struct {
	AppointmentID int32 `uri:"id" binding:"required"`
}

func assignPatient(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams AssignPatientParams
		uriErr := ctx.ShouldBindUri(&uriParams)
		if uriErr != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: uriErr.Error(),
			})
			return
		}
		// Not sure about the binding here
		var params AssignPatientIDHolder
		err := ctx.ShouldBindJSON(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.AssignPatient(ctx, &appointments.AssignPatientRequest{
			Token:         ctx.GetString(middlewares.TokenKey),
			AppointmentId: uriParams.AppointmentID,
			PatientId:     params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			PatientIDHolder{
				ID: response.GetPatientId(),
			})
	}
}

type RemovePatientParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func removePatient(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params RemovePatientParams
		err := ctx.ShouldBindUri(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.RemovePatient(ctx, &appointments.RemovePatientRequest{
			Token:         ctx.GetString(middlewares.TokenKey),
			AppointmentId: params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			PatientIDHolder{
				ID: response.GetPatientId(),
			})
	}
}

type DeleteAppointmentParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func deleteAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params DeleteAppointmentParams
		err := ctx.ShouldBindUri(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		_, err = service.DeleteAppointment(ctx, &appointments.DeleteAppointmentRequest{
			Token:         ctx.GetString(middlewares.TokenKey),
			AppointmentId: params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}

func RegisterAppointmentRoutes(router *gin.Engine) {
	appointmentService, err := ms.FetchServiceParameters(resourceNameAppointment)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.NewClient(appointmentService.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := appointments.NewAppointmentsServiceClient(conn)

	router.GET("/appointment/:id", getAppointment(client))
	router.POST("/appointment", createAppointment(client))
	router.GET("/appointment", getAppointments(client))
	router.PUT("/appointment/:id/patient", assignPatient(client))
	router.DELETE("/appointment/:id/patient", removePatient(client))
	router.DELETE("/appointment/:id", deleteAppointment(client))
}

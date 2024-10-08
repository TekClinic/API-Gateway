package routes

import (
	"net/http"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/schemas"
	appointments "github.com/TekClinic/Appointments-MicroService/appointments_protobuf"
	"github.com/gin-gonic/gin"
)

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
			CreateNamedAPIResourceList(ctx, resourceNameAppointment,
				params.Skip, params.Limit, response.GetCount(), response.GetResults()))
	}
}

type AppointmentParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams AppointmentParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.GetAppointment(ctx, &appointments.GetAppointmentRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, schemas.Appointment{
			AppointmentBase: schemas.AppointmentBase{
				PatientID: response.GetPatientId(),
				DoctorID:  response.GetDoctorId(),
				StartTime: response.GetStartTime(),
				EndTime:   response.GetEndTime(),
			},
			ID:                response.GetId(),
			ApprovedByPatient: response.GetApprovedByPatient(),
			Visited:           response.GetVisited(),
		})
	}
}

func createAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var bodyParams schemas.AppointmentBase
		err := ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.CreateAppointment(ctx, &appointments.CreateAppointmentRequest{
			Token:     ctx.GetString(middlewares.TokenKey),
			PatientId: bodyParams.PatientID,
			DoctorId:  bodyParams.DoctorID,
			StartTime: bodyParams.StartTime,
			EndTime:   bodyParams.EndTime,
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

type AssignPatientParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func assignPatient(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams AssignPatientParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var bodyParams schemas.PatientIDHolder
		err = ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.AssignPatient(ctx, &appointments.AssignPatientRequest{
			Token:     ctx.GetString(middlewares.TokenKey),
			Id:        uriParams.ID,
			PatientId: bodyParams.PatientID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, schemas.PatientIDHolder{
			PatientID: response.GetPatientId(),
		})
	}
}

type RemovePatientParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func removePatient(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams RemovePatientParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.RemovePatient(ctx, &appointments.RemovePatientRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, schemas.PatientIDHolder{
			PatientID: response.GetPatientId(),
		})
	}
}

type DeleteAppointmentParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func deleteAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams DeleteAppointmentParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		_, err = service.DeleteAppointment(ctx, &appointments.DeleteAppointmentRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}

type UpdateAppointmentParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func updateAppointment(service appointments.AppointmentsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var uriParams UpdateAppointmentParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var bodyParams schemas.AppointmentUpdate
		err = ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call appointment microservice
		response, err := service.UpdateAppointment(ctx, &appointments.UpdateAppointmentRequest{
			Token:             ctx.GetString(middlewares.TokenKey),
			Id:                uriParams.ID,
			PatientId:         bodyParams.PatientID,
			DoctorId:          bodyParams.DoctorID,
			StartTime:         bodyParams.StartTime,
			EndTime:           bodyParams.EndTime,
			ApprovedByPatient: bodyParams.ApprovedByPatient,
			Visited:           bodyParams.Visited,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK, schemas.IDHolder{
			ID: response.GetId(),
		})
	}
}

func RegisterAppointmentRoutes(router *gin.Engine) {
	client := InitiateClient(resourceNameAppointment, appointments.NewAppointmentsServiceClient)

	// deprecated
	router.GET("/appointment/:id", getAppointment(client))
	router.POST("/appointment", createAppointment(client))
	router.GET("/appointment", getAppointments(client))
	router.PUT("/appointment/:id/patient", assignPatient(client))
	router.DELETE("/appointment/:id/patient", removePatient(client))
	router.DELETE("/appointment/:id", deleteAppointment(client))
	router.PUT("/appointment/:id", updateAppointment(client))
	// end deprecated

	router.GET("/appointments/:id", getAppointment(client))
	router.POST("/appointments", createAppointment(client))
	router.GET("/appointments", getAppointments(client))
	router.PUT("/appointments/:id/patient", assignPatient(client))
	router.DELETE("/appointments/:id/patient", removePatient(client))
	router.DELETE("/appointments/:id", deleteAppointment(client))
	router.PUT("/appointments/:id", updateAppointment(client))
}

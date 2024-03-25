package routes

import (
	"log"
	"net/http"
	"strings"

	sf "github.com/sa-/slicefunk"

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
		response, err := service.GetPatientsIDs(ctx, &patients.PatientsRequest{
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

type PatientParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getPatient(service patients.PatientsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var params PatientParams
		err := ctx.ShouldBindUri(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.GetPatient(ctx, &patients.PatientRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			schemas.Patient{
				PatientBase: schemas.PatientBase{
					Name: response.GetName(),
					PersonalID: schemas.PersonalID{
						ID:   response.GetPersonalId().GetId(),
						Type: response.GetPersonalId().GetType(),
					},
					Gender:      strings.ToLower(response.GetGender().String()),
					PhoneNumber: response.GetPhoneNumber(),
					Languages:   response.GetLanguages(),
					BirthDate:   response.GetBirthDate(),
					EmergencyContacts: sf.Map(response.GetEmergencyContacts(),
						func(contact *patients.Patient_EmergencyContact) schemas.EmergencyContact {
							return schemas.EmergencyContact{
								Name:      contact.GetName(),
								Closeness: contact.GetCloseness(),
								Phone:     contact.GetPhone(),
							}
						}),
					ReferredBy:  response.GetReferredBy(),
					SpecialNote: response.GetSpecialNote(),
				},
				ID:     response.GetId(),
				Active: response.GetActive(),
				Age:    response.GetAge(),
			})
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
	router.GET("/patient/:id", getPatient(client))
}

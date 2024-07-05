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
)

const resourceNamePatient = "patient"

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
		response, err := service.GetPatientsIDs(ctx, &patients.GetPatientsIDsRequest{
			Token:  ctx.GetString(middlewares.TokenKey),
			Limit:  params.Limit,
			Offset: params.Skip,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			CreateNamedAPIResourceList(ctx, resourceNamePatient,
				params.Skip, params.Limit, response.GetCount(), response.GetResults()))
	}
}

type PatientParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getPatient(service patients.PatientsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var uriParams PatientParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.GetPatient(ctx, &patients.GetPatientRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    uriParams.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		if response.GetPatient() == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, schemas.ErrorResponse{
				Message: "Invalid response from the server.",
			})
			return
		}

		patient := response.GetPatient()

		languages := patient.GetLanguages()
		if languages == nil {
			languages = []string{}
		}

		ctx.JSON(http.StatusOK,
			schemas.Patient{
				PatientBase: schemas.PatientBase{
					Name: patient.GetName(),
					PersonalID: schemas.PersonalID{
						ID:   patient.GetPersonalId().GetId(),
						Type: patient.GetPersonalId().GetType(),
					},
					Gender:      strings.ToLower(patient.GetGender().String()),
					PhoneNumber: patient.GetPhoneNumber(),
					Languages:   languages,
					BirthDate:   patient.GetBirthDate(),
					EmergencyContacts: sf.Map(patient.GetEmergencyContacts(),
						func(contact *patients.Patient_EmergencyContact) schemas.EmergencyContact {
							return schemas.EmergencyContact{
								Name:      contact.GetName(),
								Closeness: contact.GetCloseness(),
								Phone:     contact.GetPhone(),
							}
						}),
					ReferredBy:  patient.GetReferredBy(),
					SpecialNote: patient.GetSpecialNote(),
				},
				ID:     patient.GetId(),
				Active: patient.GetActive(),
				Age:    patient.GetAge(),
			})
	}
}

func createPatient(service patients.PatientsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the body
		var bodyParams schemas.PatientBase
		err := ctx.ShouldBindJSON(&bodyParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		response, err := service.CreatePatient(ctx, &patients.CreatePatientRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Name:  bodyParams.Name,
			PersonalId: &patients.Patient_PersonalID{
				Id:   bodyParams.PersonalID.ID,
				Type: bodyParams.PersonalID.Type,
			},
			Gender:      patients.Patient_Gender(patients.Patient_Gender_value[strings.ToUpper(bodyParams.Gender)]),
			PhoneNumber: bodyParams.PhoneNumber,
			Languages:   bodyParams.Languages,
			BirthDate:   bodyParams.BirthDate,
			EmergencyContacts: sf.Map(bodyParams.EmergencyContacts,
				func(contact schemas.EmergencyContact) *patients.Patient_EmergencyContact {
					return &patients.Patient_EmergencyContact{
						Name:      contact.Name,
						Closeness: contact.Closeness,
						Phone:     contact.Phone,
					}
				}),
			ReferredBy:  bodyParams.ReferredBy,
			SpecialNote: bodyParams.SpecialNote,
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

func deletePatient(service patients.PatientsServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var uriParams PatientParams
		err := ctx.ShouldBindUri(&uriParams)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call patient microservice
		_, err = service.DeletePatient(ctx, &patients.DeletePatientRequest{
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

func RegisterPatientRoutes(router *gin.Engine) {
	patientsService, err := ms.FetchServiceParameters(resourceNamePatient)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.NewClient(patientsService.GetAddr(), grpc.WithTransportCredentials(GetTransportCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := patients.NewPatientsServiceClient(conn)
	router.GET("/patient", getPatients(client))
	router.POST("/patient", createPatient(client))
	router.GET("/patient/:id", getPatient(client))
	router.DELETE("/patient/:id", deletePatient(client))
}

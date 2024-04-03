package routes

import (
	"log"
	"net/http"
	"strings"

	sf "github.com/sa-/slicefunk"

	"github.com/TekClinic/API-Gateway/middlewares"
	"github.com/TekClinic/API-Gateway/schemas"
	medRecords "github.com/TekClinic/Medical-Records-Microservice/MR_protobuf"
	ms "github.com/TekClinic/MicroService-Lib"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const resourceName = "medical_record"

type RecordParams struct {
	Skip  int32 `form:"skip,default=0"`
	Limit int32 `form:"limit,default=20"`
}

func geRecords(service medRecords.RecordServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the query
		var params RecordParams
		err := ctx.ShouldBindQuery(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call medical records microservice
		response, err := service.GetRecordsIDs(ctx, &medRecords.RecordsRequest{
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

type RecordParams struct {
	ID int32 `uri:"id" binding:"required"`
}

func getRecord(service medRecords.RecordServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the path
		var params RecordParams
		err := ctx.ShouldBindUri(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call medical record microservice
		response, err := service.GetRecord(ctx, &medRecords.RecordsRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Id:    params.ID,
		})
		if err != nil {
			HandleGRPCError(err, ctx)
			return
		}

		ctx.JSON(http.StatusOK,
			schemas.Patient{
				MedicalRecordBase: schemas.MedicalRecordBase{
					PatientID: response.GetPatientID(),
					DoctorID:  response.GetDoctorID(),
					Date:      response.GetDate(),
					Note:      response.GetNote(),
			},
			ID:     response.GetId(),
	}
}

func addPatient(service medRecords.RecordServiceClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// fetch params from the body
		var params schemas.RecordBase
		err := ctx.ShouldBindJSON(&params)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, schemas.ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		// call medical records microservice
		response, err := service.CreateMedicalRecord(ctx, &medRecords.CreateRecordRequest{
			Token: ctx.GetString(middlewares.TokenKey),
			Record: &medRecords.Record{
				PatientID: params.PatientID,
				DoctorID:  params.DoctorID,
				Date:      params.Date,
				Note:      params.Note,
			},

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

func RegisterMedicalRecordRoutes(router *gin.Engine) {
	RecordsService, err := ms.FetchServiceParameters(resourceName)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial(RecordsService.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := medRecords.NewRecordsServiceClient(conn)
	router.GET("/record", getRecords(client))
	router.GET("/record/:id", getRecord(client))
	router.POST("/record", addRecord(client))
}

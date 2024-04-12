package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yusufwib/blockchain-medical-record/models/dappointment"
	"github.com/yusufwib/blockchain-medical-record/models/dmedicalrecord"
	"github.com/yusufwib/blockchain-medical-record/service"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"
	"github.com/yusufwib/blockchain-medical-record/utils/trace_id"

	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"

	"github.com/labstack/echo/v4"
)

type (
	IAppointmentHandler interface {
		CreateAppointment(ctx echo.Context) error
		FindAppointmentByPatientID(ctx echo.Context) error
		UpdateAppointmentStatus(ctx echo.Context) error
		UploadFile(ctx echo.Context) error
		WriteMedicalRecord(ctx echo.Context) error
		FindMedicalRecordByID(ctx echo.Context) error
	}

	AppointmentHandler struct {
		Context            context.Context
		Logger             mlog.Logger
		Validator          mvalidator.Validator
		AppointmentService service.AppointmentService
	}
)

func NewAppointmentHandler(
	context context.Context,
	logger mlog.Logger,
	validator mvalidator.Validator,
	appointmentService service.AppointmentService,
) IAppointmentHandler {
	return &AppointmentHandler{
		Context:            context,
		Logger:             logger,
		Validator:          validator,
		AppointmentService: appointmentService,
	}
}

func (i *AppointmentHandler) UploadFile(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	req := dmedicalrecord.UploadFileRequest{
		File: file,
	}

	path, err := i.AppointmentService.UploadFile(usecaseContext, req)
	if err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return SuccessResponse(ctx, http.StatusOK, map[string]string{
		"path": path,
	})
}

func (i *AppointmentHandler) FindAppointmentByPatientID(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	var ID uint64
	if ctx.Get("patient_id").(uint64) != 0 {
		ID, _ = ctx.Get("patient_id").(uint64)
	} else {
		ID, _ = ctx.Get("doctor_id").(uint64)
	}

	ID = 3

	i.Logger.InfoT(traceID, "get appointment by patient id", mlog.Any("id(p/d)", ID))

	healthServiceID, _ := strconv.ParseUint(ctx.QueryParam("health_service_id"), 0, 64)
	patientID, _ := strconv.ParseUint(ctx.QueryParam("patient_id"), 0, 64)
	appointmentID, _ := strconv.ParseUint(ctx.QueryParam("appointment_id"), 0, 64)
	filter := dappointment.AppointmentFilter{
		ScheduleDate:    ctx.QueryParam("schedule_date"),
		Status:          ctx.QueryParam("status"),
		HealthServiceID: healthServiceID,
		IsDoctor:        ctx.QueryParam("is_doctor") == "true",
		PatientID:       patientID,
		AppointmentID:   appointmentID,
	}

	if user, err := i.AppointmentService.FindAppointmentByPatientID(usecaseContext, ID, filter); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else if len(user) == 0 {
		return ErrorResponse(ctx, http.StatusNotFound, "No appointments found", nil)
	} else {
		return SuccessResponse(ctx, http.StatusOK, user)
	}
}

func (i *AppointmentHandler) CreateAppointment(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	var req dappointment.AppointmentCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "bad request", nil)
	}

	if mapErr, err := i.Validator.Struct(req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "invalid payload", mapErr)
	}

	ID, _ := ctx.Get("patient_id").(uint64)
	i.Logger.InfoT(traceID, "create appointment", mlog.Any("payload", req))

	if err := i.AppointmentService.CreateAppointment(usecaseContext, ID, req); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else {
		return SuccessResponse(ctx, http.StatusCreated, nil)
	}
}

func (i *AppointmentHandler) UpdateAppointmentStatus(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	ID, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	var req dappointment.AppointmentUpdateStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "bad request", nil)
	}

	if mapErr, err := i.Validator.Struct(req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "invalid payload", mapErr)
	}

	i.Logger.InfoT(traceID, "update status appointment", mlog.Any("payload", req))

	if err := i.AppointmentService.UpdateAppointmentStatus(usecaseContext, ID, req); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else {
		return SuccessResponse(ctx, http.StatusCreated, nil)
	}
}

func (i *AppointmentHandler) WriteMedicalRecord(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	ID, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	var req dmedicalrecord.MedicalRecordRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "bad request", nil)
	}

	req.AppointmentID = ID
	if mapErr, err := i.Validator.Struct(req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "invalid payload", mapErr)
	}

	i.Logger.InfoT(traceID, "write medical record", mlog.Any("payload", req))

	if err := i.AppointmentService.WriteMedicalRecord(usecaseContext, req); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else {
		return SuccessResponse(ctx, http.StatusCreated, nil)
	}
}

func (i *AppointmentHandler) FindMedicalRecordByID(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	ID, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	i.Logger.InfoT(traceID, "get medical record by appointment id", mlog.Any("id", ID))

	res := i.AppointmentService.FindMedicalRecordByID(usecaseContext, ID)
	return SuccessResponse(ctx, http.StatusOK, res)
}

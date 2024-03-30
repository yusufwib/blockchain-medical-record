package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yusufwib/blockchain-medical-record/service"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"
	"github.com/yusufwib/blockchain-medical-record/utils/trace_id"

	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"

	"github.com/labstack/echo/v4"
)

type (
	IHealthHandler interface {
		FindAll(ctx echo.Context) error
		FindDoctorByHealthID(ctx echo.Context) error
	}

	HealthHandler struct {
		Context       context.Context
		Logger        mlog.Logger
		Validator     mvalidator.Validator
		HealthService service.HealthService
	}
)

func NewHealthHandler(
	context context.Context,
	logger mlog.Logger,
	validator mvalidator.Validator,
	healthService service.HealthService,
) IHealthHandler {
	return &HealthHandler{
		Context:       context,
		Logger:        logger,
		Validator:     validator,
		HealthService: healthService,
	}
}

func (i *HealthHandler) FindAll(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	i.Logger.InfoT(traceID, "get health services")

	if res, err := i.HealthService.FindAll(usecaseContext); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else if len(res) == 0 {
		return ErrorResponse(ctx, http.StatusNotFound, "No health service found", nil)
	} else {
		return SuccessResponse(ctx, http.StatusOK, res)
	}
}

func (i *HealthHandler) FindDoctorByHealthID(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	ID, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	i.Logger.InfoT(traceID, "get doctor by health service id", mlog.Any("id", ID))

	if user, err := i.HealthService.FindDoctorByHealthID(usecaseContext, ID); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else if len(user) == 0 {
		return ErrorResponse(ctx, http.StatusNotFound, "No doctors found", nil)
	} else {
		return SuccessResponse(ctx, http.StatusOK, user)
	}
}

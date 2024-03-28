package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/yusufwib/blockchain-medical-record/models/duser"
	"github.com/yusufwib/blockchain-medical-record/service"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"
	"github.com/yusufwib/blockchain-medical-record/utils/trace_id"

	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"

	"github.com/labstack/echo/v4"
)

type (
	IUserHandler interface {
		Register(ctx echo.Context) error
		Login(ctx echo.Context) error
		FindByID(ctx echo.Context) error
	}

	UserHandler struct {
		Context     context.Context
		Logger      mlog.Logger
		Validator   mvalidator.Validator
		UserService service.UserService
	}
)

func NewUserHandler(
	context context.Context,
	logger mlog.Logger,
	validator mvalidator.Validator,
	userService service.UserService,
) IUserHandler {
	return &UserHandler{
		Context:     context,
		Logger:      logger,
		Validator:   validator,
		UserService: userService,
	}
}

// FindByID swagger
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} duser.UserResponse
// @Router /users/{id} [get]
func (i *UserHandler) FindByID(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	ID, err := strconv.ParseUint(ctx.Param("id"), 0, 64)
	if err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, err.Error(), nil)
	}

	i.Logger.InfoT(traceID, "get employee by id", mlog.Any("id", ID))

	if employee, err := i.UserService.FindByID(usecaseContext, ID); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else if employee.IsEmpty() {
		return ErrorResponse(ctx, http.StatusNotFound, "No employees found", nil)
	} else {
		return SuccessResponse(ctx, http.StatusOK, employee)
	}
}

func (i *UserHandler) Login(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	var req duser.UserLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "bad request", nil)
	}

	if mapErr, err := i.Validator.Struct(req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "invalid payload", mapErr)
	}

	i.Logger.InfoT(traceID, "login", mlog.Any("payload", map[string]interface{}{
		"email": req.Email,
		"type":  req.Type,
	}))

	if res, err := i.UserService.Login(usecaseContext, req); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else {
		return SuccessResponse(ctx, http.StatusOK, res)
	}
}

func (i *UserHandler) Register(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	usecaseContext := trace_id.SetIDx(ctx.Request().Context(), traceID)

	var req duser.UserRegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "bad request", nil)
	}

	if mapErr, err := i.Validator.Struct(req); err != nil {
		return ErrorResponse(ctx, http.StatusBadRequest, "invalid payload", mapErr)
	}

	i.Logger.InfoT(traceID, "register", mlog.Any("payload", req))

	if res, err := i.UserService.Register(usecaseContext, req); err != nil {
		return ErrorResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
	} else {
		return SuccessResponse(ctx, http.StatusCreated, res)
	}
}

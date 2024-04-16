package handler

import (
	"context"
	"net/http"

	"github.com/yusufwib/blockchain-medical-record/service"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"
	"github.com/yusufwib/blockchain-medical-record/utils/trace_id"

	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"

	"github.com/labstack/echo/v4"
)

type (
	IBlockchainHandler interface {
		MineAll(ctx echo.Context) error
	}

	BlockchainHandler struct {
		Context           context.Context
		Logger            mlog.Logger
		Validator         mvalidator.Validator
		BlockchainService service.BlockchainService
	}
)

func NewBlockchainHandler(
	context context.Context,
	logger mlog.Logger,
	validator mvalidator.Validator,
	blockchainService service.BlockchainService,
) IBlockchainHandler {
	return &BlockchainHandler{
		Context:           context,
		Logger:            logger,
		Validator:         validator,
		BlockchainService: blockchainService,
	}
}

func (i *BlockchainHandler) MineAll(ctx echo.Context) error {
	traceID := trace_id.GetID(ctx)
	i.Logger.InfoT(traceID, "mine all medical records")

	res := i.BlockchainService.MineAll()
	if len(res) == 0 {
		return ErrorResponse(ctx, http.StatusNotFound, "No medical records found", nil)
	} else {
		return ctx.JSON(http.StatusOK, res)
	}
}

package handler

import "github.com/labstack/echo/v4"

type BaseResponse struct {
	HttpCode        int                    `json:"-"`
	Status          bool                   `json:"status"`
	Data            interface{}            `json:"data"`
	Msg             string                 `json:"msg,omitempty"`
	Message         string                 `json:"message,omitempty"`
	Error           interface{}            `json:"error,omitempty"`
	ErrorValidation map[string]interface{} `json:"error_validation,omitempty"`
}

func SuccessResponse(c echo.Context, statusCode int, data interface{}) error {
	resp := &BaseResponse{
		HttpCode: statusCode,
		Data:     data,
		Status:   true,
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(statusCode)
	return c.JSON(statusCode, resp)
}

func ErrorResponse(c echo.Context, errorCode int, messages string, errorValidation map[string]interface{}) error {
	resp := &BaseResponse{
		HttpCode:        errorCode,
		Message:         messages,
		ErrorValidation: errorValidation,
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(errorCode)
	return c.JSON(errorCode, resp)
}

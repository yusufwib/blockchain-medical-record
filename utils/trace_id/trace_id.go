package trace_id

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const requestIDKey = "request-id"

type keyCtx string

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(requestIDKey, uuid.NewString())
		return next(c)
	}
}

func GetID(ctx echo.Context) string {
	requestID, ok := ctx.Get(requestIDKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

func SetIDx(ctx context.Context, requsetID string) context.Context {
	return context.WithValue(ctx, keyCtx(requestIDKey), requsetID)
}

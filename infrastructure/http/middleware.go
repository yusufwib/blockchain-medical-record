package infrastructure_http

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/yusufwib/blockchain-medical-record/handler"
	"github.com/yusufwib/blockchain-medical-record/infrastructure"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (httpServer *HttpServer) PrepareMiddleware(app *infrastructure.App) {
	httpServer.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	httpServer.Echo.Use(middleware.Recover())

	httpServer.Echo.HTTPErrorHandler = func(err error, ctx echo.Context) {
		if ctx.Response().Committed {
			return
		}
		code := http.StatusInternalServerError
		he, ok := err.(*echo.HTTPError)
		if ok {
			code = he.Code
		}
		errMSg := map[string]interface{}{
			"data":    nil,
			"error":   err,
			"message": err.Error(),
		}
		ctx.JSON(code, errMSg)
	}
}

func JWTMiddleware(jwtKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return handler.ErrorResponse(c, http.StatusUnauthorized, "missing JWT token", nil)
			}

			// Extract the token from the Authorization header
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return handler.ErrorResponse(c, http.StatusUnauthorized, "invalid token format", nil)
			}

			token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Check the token signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, handler.ErrorResponse(c, http.StatusUnauthorized, "invalid token", nil)
				}
				return jwtKey, nil
			})

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return handler.ErrorResponse(c, http.StatusUnauthorized, "invalid token claims", nil)
			}

			userID, _ := claims["id"].(float64)
			if !ok {
				return handler.ErrorResponse(c, http.StatusUnauthorized, "invalid user ID in token", nil)
			}

			userType, ok := claims["type"].(string)
			if !ok {
				return handler.ErrorResponse(c, http.StatusUnauthorized, "invalid user type in token", nil)
			}

			c.Set("id", uint64(userID))
			c.Set("type", userType)
			return next(c)
		}
	}
}

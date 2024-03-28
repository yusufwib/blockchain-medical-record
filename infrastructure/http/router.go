package infrastructure_http

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/yusufwib/blockchain-medical-record/docs"
	"github.com/yusufwib/blockchain-medical-record/infrastructure"
)

// @title Employee API
// @version 1.0
func (httpServer *HttpServer) PrepareRoute(app *infrastructure.App) {
	dependency := infrastructure.NewDependency(app.Context, app.Logger, app.Validator, app.Cfg, app.Database)

	v1 := httpServer.Echo.Group("/v1")
	v1User := v1.Group("/users")

	v1User.POST("/login", dependency.UserHandler.Login)
	v1User.POST("/register", dependency.UserHandler.Register)

	v1User.Use(JWTMiddleware(app.Cfg.Server.JWTSecretKey))
	v1User.GET("/:id", dependency.UserHandler.FindByID)

	httpServer.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

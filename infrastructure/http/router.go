package infrastructure_http

import (
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/yusufwib/blockchain-medical-record/docs"
	"github.com/yusufwib/blockchain-medical-record/infrastructure"
)

// @title Medical Records API
// @version 1.0
func (httpServer *HttpServer) PrepareRoute(app *infrastructure.App) {
	dependency := infrastructure.NewDependency(app.Context, app.Logger, app.Validator, app.Cfg, app.Database)

	v1 := httpServer.Echo.Group("/v1")
	v1User := v1.Group("/users")

	v1User.POST("/login", dependency.UserHandler.Login)
	v1User.POST("/register", dependency.UserHandler.Register)

	v1User.Use(JWTMiddleware(app.Cfg.Server.JWTSecretKey))
	v1User.GET("/details", dependency.UserHandler.GetDetails)
	v1User.GET("/:id", dependency.UserHandler.FindByID)

	v1Service := v1.Group("/services")
	v1Service.Use(JWTMiddleware(app.Cfg.Server.JWTSecretKey))

	v1Service.GET("", dependency.HealthHandler.FindAll)
	v1Service.GET("/:id/doctors", dependency.HealthHandler.FindDoctorByHealthID)

	v1Appointment := v1.Group("/appointments")
	v1Appointment.Use(JWTMiddleware(app.Cfg.Server.JWTSecretKey))
	v1Appointment.GET("", dependency.AppointmentHandler.FindAppointmentByPatientID)
	v1Appointment.POST("", dependency.AppointmentHandler.CreateAppointment)
	v1Appointment.PATCH("/:id", dependency.AppointmentHandler.UpdateAppointmentStatus)

	httpServer.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

package infrastructure_http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/yusufwib/blockchain-medical-record/docs"
	"github.com/yusufwib/blockchain-medical-record/infrastructure"
)

// @title Medical Records API
// @version 1.0
func (httpServer *HttpServer) PrepareRoute(app *infrastructure.App) {
	dependency := infrastructure.NewDependency(app.Context, app.Logger, app.Validator, app.Cfg, app.Database, app.LevelDB)

	httpServer.Echo.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Blockchain Medical Record API",
		})
	})

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
	v1Appointment.PUT("/:id", dependency.AppointmentHandler.WriteMedicalRecord)
	v1Appointment.GET("/:id", dependency.AppointmentHandler.FindMedicalRecordByID)
	v1Appointment.GET("/details/:id", dependency.AppointmentHandler.FindAppointmentDetailByID)

	v1Appointment.GET("/export/:id", dependency.AppointmentHandler.ExportMedicalRecordByID)
	v1Appointment.POST("/upload", dependency.AppointmentHandler.UploadFile)

	httpServer.Echo.GET("/swagger/*", echoSwagger.WrapHandler)
}

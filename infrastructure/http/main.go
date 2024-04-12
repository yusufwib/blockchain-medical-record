package infrastructure_http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yusufwib/blockchain-medical-record/infrastructure"
	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
	"github.com/yusufwib/blockchain-medical-record/utils/trace_id"

	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	Echo   *echo.Echo
	Logger mlog.Logger
}

func NewHttpServer(app infrastructure.App) *HttpServer {
	httpServer := &HttpServer{
		Logger: app.Logger,
	}

	httpServer.Echo = echo.New()
	httpServer.Echo.HideBanner = true

	// Setup Middleware
	httpServer.PrepareMiddleware(&app)

	// Setup Router
	httpServer.PrepareRoute(&app)

	// Use Trace ID
	httpServer.Echo.Use(trace_id.RequestID)
	httpServer.Echo.Static("/public", "public")

	// Start echo server on goroutine
	go func() {
		if err := httpServer.Echo.Start(fmt.Sprintf("0.0.0.0:%v", app.Cfg.Server.HTTPPort)); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server, failed to start web service: %v\n", err)
		}
	}()

	<-app.TerminalHandler
	app.Logger.Info("Running cleanup tasks...")

	// Gracefull shutdown stage
	app.Logger.Info("Shutdown echo server...")

	ctx, cancel := context.WithTimeout(app.Context, time.Minute)
	defer cancel()

	if err := httpServer.Echo.Shutdown(ctx); err != nil {
		httpServer.Echo.Logger.Fatal(err)
	}

	return httpServer
}

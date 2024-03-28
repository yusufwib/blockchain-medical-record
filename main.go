package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/datasource"
	"github.com/yusufwib/blockchain-medical-record/infrastructure"
	infrastructure_http "github.com/yusufwib/blockchain-medical-record/infrastructure/http"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"

	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
)

const (
	appName = "api"
	version = "1.0.0"
)

func main() {
	cfg := config.LoadConfig("config/.env")
	logger := mlog.New("info", "stdout")

	database := datasource.NewPostgreSQLSession(&cfg.PostgreSQLConfig)
	instDB, _ := (*database).DB()
	defer instDB.Close()

	terminalHandler := make(chan os.Signal, 1)
	signal.Notify(terminalHandler, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	application := infrastructure.App{
		TerminalHandler: terminalHandler,
		Cfg:             cfg,
		Logger:          logger,
		Database:        database,
		Validator:       mvalidator.New(),
		Context:         context.Background(),
	}

	infrastructure_http.NewHttpServer(application)
}

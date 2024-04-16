package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/datasource"
	"github.com/yusufwib/blockchain-medical-record/infrastructure"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"

	infrastructure_http "github.com/yusufwib/blockchain-medical-record/infrastructure/http"
	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
)

const (
	appName = "blockchain-medical-record"
	version = "1.0.1"
)

func main() {
	cfg := config.LoadConfig("config/.env")
	logger := mlog.New("info", "stdout")

	logger.Info(fmt.Sprintf("%s version: %s", appName, version))

	database := datasource.NewPostgreSQLSession(&cfg.PostgreSQLConfig)
	instDB, _ := (*database).DB()
	defer instDB.Close()

	levelDB, err := leveldb.OpenFile("blockchain_db", nil)
	if err != nil {
		logger.Error("err while open levelDB", err)
	}

	terminalHandler := make(chan os.Signal, 1)
	signal.Notify(terminalHandler, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	application := infrastructure.App{
		TerminalHandler: terminalHandler,
		Cfg:             cfg,
		Logger:          logger,
		Database:        database,
		LevelDB:         levelDB,
		Validator:       mvalidator.New(),
		Context:         context.Background(),
	}

	// run nodes

	infrastructure_http.NewHttpServer(application)
}

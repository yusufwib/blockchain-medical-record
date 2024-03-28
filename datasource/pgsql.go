package datasource

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgreSQLConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	User             string `json:"user"`
	Password         string `json:"-"`
	Database         string `json:"database"`
	MaxIdleConns     int    `json:"max_idle_conns"`
	MaxOpenConns     int    `json:"max_open_conns"`
	ConnMaxLifetime  int    `json:"conn_max_lifetime"`
	Environment      string `json:"environment"`
	SlowLogThreshold int    `json:"slow_log_threshold"`
}

func NewPostgreSQLSession(dbConfig *PostgreSQLConfig) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.Port,
	)

	if dbConfig.SlowLogThreshold == 0 {
		dbConfig.SlowLogThreshold = 200 // Default: 200 Milliseconds
	}

	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel:                  logger.Warn,
				SlowThreshold:             time.Duration(dbConfig.SlowLogThreshold) * time.Millisecond,
				IgnoreRecordNotFoundError: true,
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		log.Panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Panic(err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetime) * time.Second)

	return db
}

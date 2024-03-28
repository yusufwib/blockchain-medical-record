package config

import (
	"log"

	"github.com/yusufwib/blockchain-medical-record/datasource"
	"github.com/yusufwib/blockchain-medical-record/models/server"
	"github.com/yusufwib/blockchain-medical-record/utils/envar"

	"github.com/joho/godotenv"
)

type ConfigGroup struct {
	Server           server.Server
	PostgreSQLConfig datasource.PostgreSQLConfig
}

func LoadConfig(envFileLocation string) *ConfigGroup {
	if err := godotenv.Load(envFileLocation); err != nil {
		log.Printf("%s notfound", envFileLocation)
	}

	return &ConfigGroup{
		Server: server.Server{
			AppName:      envar.GetEnv("APP_NAME", ""),
			HTTPPort:     envar.GetEnv("APP_HTTP_PORT", 9009),
			JWTSecretKey: envar.GetEnv("JWT_SECRET_KEY", "secret"),
		},
		PostgreSQLConfig: datasource.PostgreSQLConfig{
			Host:             envar.GetEnv("DATABASE_HOST", "localhost"),
			Port:             envar.GetEnv("DATABASE_PORT", 3306),
			User:             envar.GetEnv("DATABASE_USERNAME", "root"),
			Password:         envar.GetEnv("DATABASE_PASSWORD", ""),
			Database:         envar.GetEnv("DATABASE_NAME", ""),
			MaxIdleConns:     envar.GetEnv("DATABASE_MAX_IDLE", 20),
			MaxOpenConns:     envar.GetEnv("DATABASE_MAX_CONN", 100),
			ConnMaxLifetime:  envar.GetEnv("DATABASE_CONN_LIFETIME", 180),
			SlowLogThreshold: envar.GetEnv("DATABASE_SLOW_LOG_THRESHOLD", 6000),
		},
	}
}

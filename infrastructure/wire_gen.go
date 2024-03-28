// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package infrastructure

import (
	"context"
	"github.com/google/wire"
	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/handler"
	"github.com/yusufwib/blockchain-medical-record/repository"
	"github.com/yusufwib/blockchain-medical-record/service"
	"github.com/yusufwib/blockchain-medical-record/utils/logger"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"
	"gorm.io/gorm"
)

// Injectors from dependency.go:

func NewDependency(context2 context.Context, logger mlog.Logger, validator mvalidator.Validator, Config *config.ConfigGroup, Database *gorm.DB) *Dependency {
	userRepository := repository.NewUserRepository(Database, Config)
	userService := service.NewUserService(userRepository)
	iUserHandler := handler.NewUserHandler(context2, logger, validator, userService)
	dependency := &Dependency{
		UserHandler: iUserHandler,
	}
	return dependency
}

// dependency.go:

type Dependency struct {
	UserHandler handler.IUserHandler
}

var setUserHandler = wire.NewSet(repository.NewUserRepository, service.NewUserService, handler.NewUserHandler)

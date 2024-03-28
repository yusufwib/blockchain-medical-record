//go:build wireinject
// +build wireinject

package infrastructure

import (
	"context"

	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/handler"
	"github.com/yusufwib/blockchain-medical-record/repository"
	"github.com/yusufwib/blockchain-medical-record/service"
	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"

	"github.com/google/wire"
	"gorm.io/gorm"
)

type Dependency struct {
	UserHandler handler.IUserHandler
}

func NewDependency(
	context context.Context,
	logger mlog.Logger,
	validator mvalidator.Validator,
	config *config.ConfigGroup,
	database *gorm.DB,
) *Dependency {
	wire.Build(
		setUserHandler,
		wire.Struct(new(Dependency), "*"),
	)
	return nil
}

var setUserHandler = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
)

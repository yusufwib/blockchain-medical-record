//go:build wireinject
// +build wireinject

package infrastructure

import (
	"context"

	"github.com/syndtr/goleveldb/leveldb"
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
	UserHandler        handler.IUserHandler
	HealthHandler      handler.IHealthHandler
	AppointmentHandler handler.IAppointmentHandler
	BlockchainHandler  handler.IBlockchainHandler
}

func NewDependency(
	context context.Context,
	logger mlog.Logger,
	validator mvalidator.Validator,
	config *config.ConfigGroup,
	database *gorm.DB,
	levelDB *leveldb.DB,
) *Dependency {
	wire.Build(
		setUserHandler,
		setHealthHandler,
		setAppointmentHandler,
		setBlockchainHandler,
		wire.Struct(new(Dependency), "*"),
	)
	return nil
}

var setUserHandler = wire.NewSet(
	repository.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
)

var setHealthHandler = wire.NewSet(
	repository.NewHealthRepository,
	service.NewHealthService,
	handler.NewHealthHandler,
)

var setAppointmentHandler = wire.NewSet(
	repository.NewAppointmentRepository,
	repository.NewBlockchainRepository,
	service.NewAppointmentService,
	handler.NewAppointmentHandler,
)

var setBlockchainHandler = wire.NewSet(
	service.NewBlockchainService,
	handler.NewBlockchainHandler,
)

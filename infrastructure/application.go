package infrastructure

import (
	"context"
	"os"

	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/utils/mvalidator"

	"gorm.io/gorm"

	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
)

type App struct {
	TerminalHandler chan os.Signal
	Cfg             *config.ConfigGroup
	Logger          mlog.Logger
	Database        *gorm.DB
	Validator       mvalidator.Validator
	Context         context.Context
}

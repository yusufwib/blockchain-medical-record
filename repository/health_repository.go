package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/models/dhealthservice"
	"github.com/yusufwib/blockchain-medical-record/models/duser"
	mlog "github.com/yusufwib/blockchain-medical-record/utils/logger"
	"gorm.io/gorm"
)

type HealthRepository struct {
	DB     *gorm.DB
	Config *config.ConfigGroup
	Logger mlog.Logger
}

func NewHealthRepository(DB *gorm.DB, cfg *config.ConfigGroup, log mlog.Logger) HealthRepository {
	return HealthRepository{DB, cfg, log}
}

func (r *HealthRepository) session(ctx context.Context) *gorm.DB {
	trx, ok := ctx.Value("pg").(*gorm.DB)
	if !ok {
		return r.DB
	}
	return trx
}

func (r *HealthRepository) FindAll(ctx context.Context) (res []dhealthservice.HealthService, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err = trx.WithContext(ctxWT).Table(dhealthservice.TableName()).Find(&res).Error; err != nil {
		return nil, fmt.Errorf("error while retrieving health service: %w", err)
	}

	return
}

func (r *HealthRepository) FindDoctorByHealthID(ctx context.Context, ID uint64) (res []duser.UserResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := trx.WithContext(ctxWT).Table(duser.TableName()).Select("doctors.*, users.*, health_services.name AS health_service_name, doctors.id AS doctor_id").
		Joins("LEFT JOIN doctors ON users.id = doctors.user_id").
		Joins("JOIN health_services ON doctors.health_service_id = health_services.id").
		Where("doctors.health_service_id = ?", ID).
		Find(&res).Error; err != nil {
		return res, fmt.Errorf("err while get user by id: %w", err)
	}

	return
}

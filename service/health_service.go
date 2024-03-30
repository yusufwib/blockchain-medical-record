package service

import (
	"context"

	"github.com/yusufwib/blockchain-medical-record/models/dhealthservice"
	"github.com/yusufwib/blockchain-medical-record/models/duser"
	"github.com/yusufwib/blockchain-medical-record/repository"
)

type HealthService struct {
	HealthRepository repository.HealthRepository
}

func NewHealthService(r repository.HealthRepository) HealthService {
	return HealthService{
		HealthRepository: r,
	}
}

func (s HealthService) FindAll(ctx context.Context) ([]dhealthservice.HealthService, error) {
	return s.HealthRepository.FindAll(ctx)
}

func (s HealthService) FindDoctorByHealthID(ctx context.Context, ID uint64) ([]duser.UserResponse, error) {
	return s.HealthRepository.FindDoctorByHealthID(ctx, ID)
}

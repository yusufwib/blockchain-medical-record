package service

import (
	"context"

	"github.com/yusufwib/blockchain-medical-record/models/duser"
	"github.com/yusufwib/blockchain-medical-record/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return UserService{
		UserRepository: r,
	}
}

func (s UserService) FindByID(ctx context.Context, ID uint64, userType string) (duser.UserResponse, error) {
	return s.UserRepository.FindByID(ctx, ID, userType)
}

func (s UserService) Login(ctx context.Context, req duser.UserLoginRequest) (duser.UserLoginResponse, error) {
	return s.UserRepository.Login(ctx, req)
}

func (s UserService) Register(ctx context.Context, req duser.UserRegisterRequest) (duser.UserResponse, error) {
	return s.UserRepository.Register(ctx, req)
}

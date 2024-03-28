package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/models/duser"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB     *gorm.DB
	Config *config.ConfigGroup
}

func NewUserRepository(DB *gorm.DB, cfg *config.ConfigGroup) UserRepository {
	return UserRepository{DB, cfg}
}

func (r *UserRepository) session(ctx context.Context) *gorm.DB {
	trx, ok := ctx.Value("pg").(*gorm.DB)
	if !ok {
		return r.DB
	}
	return trx
}

func (r *UserRepository) FindByID(ctx context.Context, ID uint64) (user duser.UserResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err = trx.WithContext(ctxWT).Table(duser.TableName()).First(&user, ID).Error; err != nil {
		return user, fmt.Errorf("err while get user by id: %w", err)
	}

	return user, nil
}

// TODO: validate using email & password
func (r *UserRepository) Login(ctx context.Context, req duser.UserLoginRequest) (res duser.UserLoginResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var user duser.User
	if err = trx.WithContext(ctxWT).Table(duser.TableName()).Where("email = ?", req.Email).First(&user).Error; err != nil {
		return res, fmt.Errorf("error while retrieving user: %w", err)
	}

	token, err := r.generateJWTToken(user)
	if err != nil {
		return res, fmt.Errorf("error while generating JWT token: %w", err)
	}

	return duser.UserLoginResponse{
		ID:    user.ID,
		Type:  user.Type,
		Token: token,
	}, nil
}

// TODO: hash password
func (r *UserRepository) Register(ctx context.Context, req duser.UserRegisterRequest) (user duser.UserResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err = trx.WithContext(ctxWT).Table(duser.TableName()).Create(&req).Error; err != nil {
		return user, fmt.Errorf("err while register user: %w", err)
	}

	if err = trx.WithContext(ctxWT).Table(duser.TableName()).First(&user, req.ID).Error; err != nil {
		return user, fmt.Errorf("err while get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) generateJWTToken(user duser.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"type":  user.Type,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(r.Config.Server.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

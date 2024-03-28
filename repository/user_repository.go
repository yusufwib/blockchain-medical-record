package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusufwib/blockchain-medical-record/config"
	"github.com/yusufwib/blockchain-medical-record/models/dpatient"
	"github.com/yusufwib/blockchain-medical-record/models/duser"
	"golang.org/x/crypto/bcrypt"
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

func (r *UserRepository) FindByID(ctx context.Context, ID uint64, userType string) (user duser.UserResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return r.getUserByID(ctxWT, trx, ID, userType)
}

func (r *UserRepository) Login(ctx context.Context, req duser.UserLoginRequest) (res duser.UserLoginResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var user duser.User
	if err = trx.WithContext(ctxWT).Table(duser.TableName()).
		Where("email = ?", req.Email).
		First(&user).Error; err != nil {
		return res, fmt.Errorf("error while retrieving user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return res, fmt.Errorf("invalid password")
	}

	if user.Type != req.Type {
		return res, fmt.Errorf("invalid user type")
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

func (r *UserRepository) Register(ctx context.Context, req duser.UserRegisterRequest) (user duser.UserResponse, err error) {
	trx := r.session(ctx)
	ctxWT, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, fmt.Errorf("error while hashing password: %w", err)
	}

	req.Password = string(hashedPassword)
	userData := req.ToUser()
	if err = trx.WithContext(ctxWT).Table(duser.TableName()).Create(&userData).Error; err != nil {
		return user, fmt.Errorf("err while register user: %w", err)
	}

	if req.Type == duser.Patient {
		patientData := req.ToPatient(userData.ID)
		if err = trx.WithContext(ctxWT).Table(dpatient.TableName()).Create(&patientData).Error; err != nil {
			return user, fmt.Errorf("err while register user: %w", err)
		}
	}

	return r.getUserByID(ctx, trx, userData.ID, string(req.Type))
}

func (r *UserRepository) getUserByID(ctx context.Context, trx *gorm.DB, ID uint64, userType string) (user duser.UserResponse, err error) {
	query := trx.WithContext(ctx).Table(duser.TableName())

	fmt.Println(userType, "userType")
	if userType == string(duser.Patient) {
		query = query.Select("patients.*, users.*").
			Joins("LEFT JOIN patients ON users.id = patients.user_id")
	} else if userType == string(duser.Doctor) {
		query = query.Select("doctors.*, users.*").
			Joins("LEFT JOIN doctors ON users.id = doctors.user_id")
	}

	if err = query.First(&user, ID).Error; err != nil {
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
		"exp":   time.Now().Add(time.Hour * 24 * 30 * 365).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(r.Config.Server.JWTSecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

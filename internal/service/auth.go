package service

import (
	"bytes"
	"errors"

	"github.com/kiricle/file-uploader/internal/lib/hash"
	"github.com/kiricle/file-uploader/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserStorage interface {
	GetUserByEmail(email string) (models.User, error)
	SaveUser(signUpDto models.SignUpDTO) (int64, error)
}

type AuthService struct {
	userStorage UserStorage
	jwtService  *JWTService
}

func NewAuthService(userStorage UserStorage, jwtService *JWTService) *AuthService {
	return &AuthService{
		userStorage: userStorage,
		jwtService:  jwtService,
	}
}

func (as *AuthService) SignUp(dto models.SignUpDTO) (int64, error) {
	foundUser, err := as.userStorage.GetUserByEmail(dto.Email)
	if err != nil {
		return 0, err
	}

	if foundUser.ID != 0 {
		return 0, ErrUserAlreadyExists
	}

	userId, err := as.userStorage.SaveUser(dto)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (as *AuthService) SignIn(dto models.SignInDTO) (string, error) {
	user, err := as.userStorage.GetUserByEmail(dto.Email)
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", ErrUserNotFound
	}

	dtoHash, err := hash.HashPassword(dto.Password)
	if err != nil {
		return "", err
	}

	if bytes.Equal(user.Password, dtoHash) {
		return "", errors.New("invalid password")
	}

	token, err := as.jwtService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

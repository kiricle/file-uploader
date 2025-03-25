package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret string
}

const exp = time.Hour * 3

var method = jwt.SigningMethodHS256

func NewJWTService(secret string) *JWTService {
	return &JWTService{secret: secret}
}

func (j *JWTService) GenerateToken(userId int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     time.Now().Add(exp).Unix(),
	}

	token := jwt.NewWithClaims(method, claims)
	signedToken, err := token.SignedString([]byte(j.secret))

	return signedToken, err
}

func (j *JWTService) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid token payload")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

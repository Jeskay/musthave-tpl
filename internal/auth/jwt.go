package auth

import (
	"errors"
	"musthave_tpl/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

type AuthService struct {
	SecretKey []byte
	ExpiresAt time.Duration
}

func NewAuthService(conf *config.GophermartConfig) *AuthService {
	return &AuthService{
		SecretKey: []byte(conf.TokenKey),
		ExpiresAt: time.Hour * 48,
	}
}

func (s *AuthService) CreateToken(login string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ExpiresAt)),
		},
		Login: login,
	})

	tokenString, err := claims.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (s *AuthService) VerifyToken(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return s.SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Login, nil
}

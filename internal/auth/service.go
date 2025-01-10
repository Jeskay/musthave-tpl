package auth

import (
	"musthave_tpl/config"
	"time"
)

type Authentication interface {
	CreateToken(login string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type AuthService struct {
	SecretKey []byte
	ExpiresAt time.Duration
}

func NewAuthService(conf *config.Config) *AuthService {
	return &AuthService{
		SecretKey: []byte(conf.TokenKey),
		ExpiresAt: time.Duration(conf.TokenExpire),
	}
}

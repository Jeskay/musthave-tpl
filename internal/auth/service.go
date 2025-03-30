// Module auth provides jwt authentication functionality to the server.
package auth

import (
	"musthave_tpl/config"
	"time"
)

type Authentication interface {
	CreateToken(login string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

// AuthService - service for creation and verification of jwt.
// It stores a secret key and expiration time.
type AuthService struct {
	SecretKey []byte
	ExpiresAt time.Duration
}

// NewAuthService creates a new instance of AuthServices from given configuration.
func NewAuthService(conf *config.Config) *AuthService {
	return &AuthService{
		SecretKey: []byte(conf.TokenKey),
		ExpiresAt: time.Duration(conf.TokenExpire),
	}
}

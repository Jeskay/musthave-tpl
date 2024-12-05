package internal

import (
	"bytes"
	"encoding/hex"
	"errors"
	"log/slog"
	"musthave_tpl/config"
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/utils"
)

type LoyaltyService struct {
	storage     Storage
	authService *auth.AuthService
	logger      *slog.Logger
	config      *config.LoyaltyConfig
}

func NewLoyaltyService(config *config.LoyaltyConfig, logger slog.Handler, storage Storage) *LoyaltyService {
	return &LoyaltyService{
		logger:      slog.New(logger),
		storage:     storage,
		config:      config,
		authService: auth.NewAuthService(config),
	}
}

func (s *LoyaltyService) Login(login string, password string) (string, error) {
	user, err := s.storage.UserByLogin(login)
	if err != nil {
		return "", err
	}
	hash, err := utils.HashBytes([]byte(password), s.config.HashKey)
	if err != nil {
		return "", err
	}
	passHash, err := hex.DecodeString(user.Password)
	if err != nil {
		return "", err
	}
	if bytes.Equal(hash, passHash) {
		return s.authService.CreateToken(login)
	}
	return "", errors.New("invalid password")
}

func (s *LoyaltyService) Authenticate(tokenString string) (*User, error) {
	login, err := s.authService.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return s.storage.UserByLogin(login)
}

func (s *LoyaltyService) Register(login string, password string) (string, error) {
	hash, err := utils.HashBytes([]byte(password), s.config.HashKey)
	if err != nil {
		return "", err
	}
	err = s.storage.AddUser(User{Login: login, Password: hex.EncodeToString(hash)})
	if err != nil {
		return "", err
	}
	return s.authService.CreateToken(login)
}

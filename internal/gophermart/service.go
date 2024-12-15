package gophermart

import (
	"bytes"
	"encoding/hex"
	"errors"
	"log/slog"
	"musthave_tpl/config"
	"musthave_tpl/internal"
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/loyalty"
	"musthave_tpl/internal/utils"
	"time"
)

type GophermartService struct {
	storage        internal.Storage
	authService    *auth.AuthService
	loyaltyService *loyalty.LoyaltyService
	logger         *slog.Logger
	config         *config.LoyaltyConfig
}

func NewGophermartService(config *config.LoyaltyConfig, logger slog.Handler, storage internal.Storage) *GophermartService {
	return &GophermartService{
		logger:         slog.New(logger),
		storage:        storage,
		config:         config,
		authService:    auth.NewAuthService(config),
		loyaltyService: loyalty.NewLoyaltyService(config, logger),
	}
}

func (s *GophermartService) Login(login string, password string) (string, error) {
	user, err := s.storage.UserByLogin(login)
	if err != nil {
		return "", err
	}
	hash, err := utils.HashBytes([]byte(password), s.config.HashKey)
	if err != nil {
		return "", &IncorrectPassword{}
	}
	passHash, err := hex.DecodeString(user.Password)
	if err != nil {
		return "", &IncorrectPassword{}
	}
	if bytes.Equal(hash, passHash) {
		return s.authService.CreateToken(login)
	}
	return "", &IncorrectPassword{}
}

func (s *GophermartService) Authenticate(tokenString string) (*internal.User, error) {
	login, err := s.authService.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return s.storage.UserByLogin(login)
}

func (s *GophermartService) Register(login string, password string) (string, error) {
	hash, err := utils.HashBytes([]byte(password), s.config.HashKey)
	if err != nil {
		return "", err
	}
	err = s.storage.AddUser(internal.User{Login: login, Password: hex.EncodeToString(hash)})
	if err != nil {
		return "", err
	}
	return s.authService.CreateToken(login)
}

func (s *GophermartService) AddOrder(login string, orderId int64) error {
	order, err := s.loyaltyService.LoyaltyAccrual(orderId)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("response conversion failed")
	}

	order.User = internal.User{Login: login}
	err = s.storage.AddOrder(*order)
	if err != nil {
		return err
	}

	return nil
}

func (s *GophermartService) Orders(login string) ([]internal.Order, error) {
	orders, err := s.storage.OrdersByUser(login)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *GophermartService) Withdrawals(login string) ([]internal.Transaction, error) {
	transactions, err := s.storage.TransactionsByUser(login)
	return transactions, err
}

func (s *GophermartService) MakeWithdrawal(login string, order int64, amount float64) error {
	user, err := s.storage.UserByLogin(login)
	if err != nil {
		return nil
	}
	if user.Balance < amount {
		return &NotEnoughFunds{}
	}
	_, err = s.storage.AddTransaction(internal.Transaction{User: login, Amount: amount, Id: order, Date: time.Now()})
	if err != nil {
		return err
	}
	return nil
}

func (s *GophermartService) GetUser(login string) (*internal.User, error) {
	user, err := s.storage.UserByLogin(login)
	return user, err
}

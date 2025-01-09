package gophermart

import (
	"bytes"
	"encoding/hex"
	"log/slog"
	"musthave_tpl/config"
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/gophermart/db"
	"musthave_tpl/internal/loyalty"
	"musthave_tpl/internal/models"
	"musthave_tpl/internal/utils"
	"time"
)

type GophermartService struct {
	storage        db.GeneralRepository
	authService    *auth.AuthService
	loyaltyService *loyalty.LoyaltyService
	logger         *slog.Logger
	config         *config.GophermartConfig
}

func NewGophermartService(config *config.GophermartConfig, logger slog.Handler, storage db.GeneralRepository) *GophermartService {
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
		return "", &models.IncorrectPassword{}
	}
	passHash, err := hex.DecodeString(user.Password)
	if err != nil {
		return "", &models.IncorrectPassword{}
	}
	if bytes.Equal(hash, passHash) {
		return s.authService.CreateToken(login)
	}
	return "", &models.IncorrectPassword{}
}

func (s *GophermartService) Authenticate(tokenString string) (*models.User, error) {
	login, err := s.authService.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return s.storage.UserByLogin(login)
}

func (s *GophermartService) Register(login string, password string) (string, error) {
	hash, err := utils.HashBytes([]byte(password), s.config.HashKey)
	if err != nil {
		s.logger.Error("", slog.String("error", err.Error()))
		return "", err
	}
	err = s.storage.AddUser(models.User{Login: login, Password: hex.EncodeToString(hash)})
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return "", err
	}
	return s.authService.CreateToken(login)
}

func (s *GophermartService) AddOrder(login string, orderID int64) error {
	order, err := s.loyaltyService.LoyaltyAccrual(orderID)
	if err != nil {
		s.logger.Error("loyalty service", slog.String("error", err.Error()))
		return err
	}
	if order == nil {
		order = &models.Order{
			Number: orderID,
			Status: models.New,
		}
	}
	order.User = models.User{Login: login}
	err = s.storage.AddOrder(*order)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *GophermartService) Orders(login string) ([]models.Order, error) {
	orders, err := s.storage.OrdersByUser(login)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return nil, err
	}
	return orders, nil
}

func (s *GophermartService) Withdrawals(login string) ([]models.Transaction, error) {
	transactions, err := s.storage.TransactionsByUser(login)
	return transactions, err
}

func (s *GophermartService) MakeWithdrawal(login string, order int64, amount float64) error {
	user, err := s.storage.UserByLogin(login)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return nil
	}
	if user.Balance < amount {
		return &models.NotEnoughFunds{}
	}
	_, err = s.storage.AddTransaction(models.Transaction{User: login, Amount: amount, ID: order, Date: time.Now()})
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *GophermartService) GetUser(login string) (*models.User, error) {
	user, err := s.storage.UserByLogin(login)
	return user, err
}

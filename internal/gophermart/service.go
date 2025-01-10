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

	"github.com/gin-gonic/gin"
)

type Gophermart interface {
	Login(ctx *gin.Context, login string, password string) (string, error)
	Register(ctx *gin.Context, login string, password string) (string, error)
	AddOrder(ctx *gin.Context, login string, orderID int64) error
	Orders(ctx *gin.Context, login string) ([]models.Order, error)
	Withdrawals(ctx *gin.Context, login string) ([]models.Transaction, error)
	MakeWithdrawal(ctx *gin.Context, login string, order int64, amount float64) error
	GetUser(ctx *gin.Context, login string) (*models.User, error)
}

type GophermartService struct {
	storage        db.GeneralRepository
	authService    auth.Authentication
	loyaltyService loyalty.Loyalty
	logger         *slog.Logger
	config         *config.Config
}

func NewGophermartService(config *config.Config, logger slog.Handler, storage db.GeneralRepository, authSvc auth.Authentication, loyaltySvc loyalty.Loyalty) *GophermartService {
	return &GophermartService{
		logger:         slog.New(logger),
		storage:        storage,
		config:         config,
		authService:    authSvc,
		loyaltyService: loyaltySvc,
	}
}

func (s *GophermartService) Login(ctx *gin.Context, login string, password string) (string, error) {
	user, err := s.storage.UserByLogin(ctx, login)
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

func (s *GophermartService) Register(ctx *gin.Context, login string, password string) (string, error) {
	hash, err := utils.HashBytes([]byte(password), s.config.HashKey)
	if err != nil {
		s.logger.Error("", slog.String("error", err.Error()))
		return "", err
	}
	user := models.User{Login: login, Password: hex.EncodeToString(hash)}
	err = s.storage.AddUser(ctx, user)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return "", err
	}
	return s.authService.CreateToken(login)
}

func (s *GophermartService) AddOrder(ctx *gin.Context, login string, orderID int64) error {
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
	err = s.storage.AddOrder(ctx, *order)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *GophermartService) Orders(ctx *gin.Context, login string) ([]models.Order, error) {
	orders, err := s.storage.OrdersByUser(ctx, login)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return nil, err
	}
	return orders, nil
}

func (s *GophermartService) Withdrawals(ctx *gin.Context, login string) ([]models.Transaction, error) {
	transactions, err := s.storage.TransactionsByUser(ctx, login)
	return transactions, err
}

func (s *GophermartService) MakeWithdrawal(ctx *gin.Context, login string, order int64, amount float64) error {
	user, err := s.storage.UserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return nil
	}
	if user.Balance < amount {
		return &models.NotEnoughFunds{}
	}
	transactions := models.Transaction{User: login, Amount: amount, ID: order, Date: time.Now()}
	_, err = s.storage.AddTransaction(ctx, transactions)
	if err != nil {
		s.logger.Error("storage request failed", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *GophermartService) GetUser(ctx *gin.Context, login string) (*models.User, error) {
	user, err := s.storage.UserByLogin(ctx, login)
	return user, err
}

package loyalty

import (
	"encoding/json"
	"log/slog"
	"musthave_tpl/config"
	"musthave_tpl/internal/loyalty/dto"
	"musthave_tpl/internal/models"
	"net/http"
	"strconv"
)

type Loyalty interface {
	LoyaltyAccrual(orderID int64) (*models.Order, error)
}

type LoyaltyService struct {
	logger *slog.Logger
	config *config.Config
	client *http.Client
}

func NewLoyaltyService(config *config.Config, logger slog.Handler) *LoyaltyService {
	return &LoyaltyService{
		logger: slog.New(logger),
		config: config,
		client: http.DefaultClient,
	}
}

func (s *LoyaltyService) LoyaltyAccrual(orderID int64) (*models.Order, error) {
	param := s.config.AccrualAddress + "/api/orders/" + strconv.FormatInt(orderID, 10)
	var order dto.Order
	res, err := s.client.Get(param)
	if err != nil {
		s.logger.Info("exited with error:", slog.Any("response data: ", res))
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	err = json.NewDecoder(res.Body).Decode(&order)
	if err != nil {
		return nil, err
	}
	converted, err := order.ToInternal()
	return &converted, err
}

package loyalty

import (
	"encoding/json"
	"log/slog"
	"musthave_tpl/config"
	"musthave_tpl/internal"
	"musthave_tpl/internal/loyalty/dto"
	"net/http"
	"strconv"
)

type LoyaltyService struct {
	logger *slog.Logger
	config *config.GophermartConfig
	client *http.Client
}

func NewLoyaltyService(config *config.GophermartConfig, logger slog.Handler) *LoyaltyService {
	return &LoyaltyService{
		logger: slog.New(logger),
		config: config,
		client: http.DefaultClient,
	}
}

func (s *LoyaltyService) LoyaltyAccrual(orderId int64) (*internal.Order, error) {
	param := "http://" + s.config.AccrualAddress + "/api/orders/" + strconv.FormatInt(orderId, 10)
	var order dto.Order
	res, err := s.client.Get(param)
	if err != nil {
		return nil, err
	}
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

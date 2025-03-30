// Module loyalty provides access to loyalty service. Loyalty service runs as a standalone HTTP server.
package loyalty

import (
	"encoding/json"
	"log/slog"
	"musthave_tpl/config"
	"musthave_tpl/internal/loyalty/dto"
	"musthave_tpl/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Loyalty is a client that interacts with loyalty HTTP server.
type Loyalty interface {
	LoyaltyAccrual(ctx *gin.Context, orderID int64) (*models.Order, error)
}

// LoyaltyService is implementation of Loyalty interface.
type LoyaltyService struct {
	logger *slog.Logger
	config *config.Config
	client *http.Client
}

// NewLoyaltyService creates new LoyaltyService instance along with HTTP client.
func NewLoyaltyService(config *config.Config, logger *slog.Logger) *LoyaltyService {
	return &LoyaltyService{
		logger: logger,
		config: config,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

// LoyaltyAccrual requests and returns Order information by given orderID.
func (s *LoyaltyService) LoyaltyAccrual(ctx *gin.Context, orderID int64) (*models.Order, error) {
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

package routes

import (
	"musthave_tpl/internal"
	"musthave_tpl/internal/handlers"
	"musthave_tpl/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(svc *internal.LoyaltyService) *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/api/user")
	{
		apiGroup.POST("/register", handlers.Register(svc), middleware.Authenticate(svc))
		apiGroup.POST("/login", handlers.Login(svc), middleware.Authenticate(svc))

		apiGroup.GET("/withdrawals", handlers.Withdrawals(svc), middleware.Authorize(svc))
		gOrders := apiGroup.Group("/orders", middleware.Authorize(svc))
		{
			gOrders.POST("", handlers.PostOrder(svc))
			gOrders.GET("", handlers.Orders(svc))
		}
		gBalance := apiGroup.Group("/balance", middleware.Authorize(svc))
		{
			gBalance.GET("", handlers.Balance(svc))
			gBalance.POST("/withdraw", handlers.MakeWithdrawal(svc))
		}
	}
	return r
}

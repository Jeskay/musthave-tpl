package routes

import (
	"musthave_tpl/internal"
	"musthave_tpl/internal/handlers"
	"musthave_tpl/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(svc *internal.LoyaltyService) *gin.Engine {
	r := gin.Default()

	r.POST("api/user/register", handlers.Register(svc)).Use(middleware.Authenticate(svc))

	authorized := r.Group("/api/user")
	{

		authorized.POST("/login", handlers.Login(svc))
		authorized.GET("/withdrawals", handlers.Withdrawals(svc))
		gOrders := authorized.Group("/orders")
		{
			gOrders.POST("", handlers.PostOrder(svc))
			gOrders.GET("", handlers.Orders(svc))
		}
		gBalance := authorized.Group("/balance")
		{
			gBalance.GET("", handlers.Balance(svc))
			gBalance.POST("/withdraw", handlers.MakeWithdrawal(svc))
		}
	}
	authorized.Use(middleware.Authorize(svc))
	return r
}

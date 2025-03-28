// Module internal contains business logic of gophermart HTTP server.
package internal

import (
	"errors"
	"io"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/dto"
	"musthave_tpl/internal/middleware"
	"musthave_tpl/internal/models"
	"musthave_tpl/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// App stores instances of Gophermart and Middleware services. Handles incoming HTTP requests and acts as HTTP server.
type App struct {
	gophermartSvc gophermart.Gophermart
	middlewareSvc middleware.Middleware
}

// NewApp creates new instance of App.
func NewApp(gophermartSvc gophermart.Gophermart, middleware middleware.Middleware) App {
	return App{
		gophermartSvc: gophermartSvc,
		middlewareSvc: middleware,
	}
}

// Router registers server endpoints and returns gin router.
func (a *App) Router() *gin.Engine {
	r := gin.Default()
	apiGroup := r.Group("/api/user")
	{
		apiGroup.POST("/register", a.Register, a.middlewareSvc.Authenticate)
		apiGroup.POST("/login", a.Login, a.middlewareSvc.Authenticate)

		gWithdrawals := apiGroup.Group("/withdrawals", a.middlewareSvc.Authorize)
		{
			gWithdrawals.GET("", a.Withdrawals)
		}
		gOrders := apiGroup.Group("/orders", a.middlewareSvc.Authorize)
		{
			gOrders.POST("", a.CreateOrder)
			gOrders.GET("", a.Orders)
		}
		gBalance := apiGroup.Group("/balance", a.middlewareSvc.Authorize)
		{
			gBalance.GET("", a.Balance)
			gBalance.POST("/withdraw", a.MakeWithdrawal)
		}
	}
	return r
}

// Login handles login requests.
func (a *App) Login(ctx *gin.Context) {
	var user dto.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := a.gophermartSvc.Login(ctx, user.Login, user.Password)
	if err != nil {
		if errors.Is(err, models.ErrIncorrectPassword) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Set("Token", token)
	ctx.Status(http.StatusOK)
}

// Register handles register requests.
func (a *App) Register(ctx *gin.Context) {
	var user dto.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := a.gophermartSvc.Register(ctx, user.Login, user.Password)
	if err != nil {
		if errors.Is(err, models.ErrUsedLogin) {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Set("Token", token)
	ctx.Status(http.StatusOK)
}

// Balance handles balance requests. Responds with client's balance information.
func (a *App) Balance(ctx *gin.Context) {
	login := ctx.GetString("Login")
	user, err := a.gophermartSvc.GetUser(ctx, login)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, dto.Balance{Balance: float64(user.Balance), Withdrawn: user.Withdrawn})
}

// MakeWithdrawal handles withdrawal requests.
func (a *App) MakeWithdrawal(ctx *gin.Context) {
	var withdrawal dto.Withdrawal
	login := ctx.GetString("Login")
	err := ctx.ShouldBindJSON(&withdrawal)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(withdrawal.OrderID, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = a.gophermartSvc.MakeWithdrawal(ctx, login, id, withdrawal.Sum)
	if err != nil {
		if errors.Is(err, models.ErrNotEnoughFunds) {
			ctx.AbortWithStatus(http.StatusPaymentRequired)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

// Withdrawals handles withdrawals requests. Responds with list of withdrawal transactions made by client.
func (a *App) Withdrawals(ctx *gin.Context) {
	login := ctx.GetString("Login")
	withdrawals, err := a.gophermartSvc.Withdrawals(ctx, login)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, dto.NewWithdrawals(withdrawals))
}

// Orders handles orders requests. Responds with list of orders associated with client.
func (a *App) Orders(ctx *gin.Context) {
	login := ctx.GetString("Login")
	if login == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	orders, err := a.gophermartSvc.Orders(ctx, login)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.JSON(http.StatusOK, dto.NewOrders(orders))
}

// CreateOrder handles create order requests.
func (a *App) CreateOrder(ctx *gin.Context) {
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	login := ctx.GetString("Login")
	orderID, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !utils.LuhnAlgorithm(string(data)) {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	if err := a.gophermartSvc.AddOrder(ctx, login, orderID); err != nil {
		if errors.Is(err, models.ErrOrderExists) {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		if errors.Is(err, models.ErrOrderUsed) {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusAccepted)
}

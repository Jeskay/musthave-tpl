package internal

import (
	"io"
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/dto"
	"musthave_tpl/internal/loyalty"
	"musthave_tpl/internal/middleware"
	"musthave_tpl/internal/models"
	"musthave_tpl/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type App struct {
	gophermartSvc *gophermart.GophermartService
	middlewareSvc *middleware.MiddlewareService
}

func NewApp(authSvc *auth.AuthService, gophermartSvc *gophermart.GophermartService, loyatySvc loyalty.LoyaltyService) App {
	return App{
		gophermartSvc: gophermartSvc,
		middlewareSvc: middleware.NewMiddleware(authSvc, gophermartSvc),
	}
}

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

func (a *App) Login(ctx *gin.Context) {
	var user dto.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := a.gophermartSvc.Login(user.Login, user.Password)
	if err != nil {
		if _, ok := err.(*models.IncorrectPassword); ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Set("Token", token)
	ctx.Status(http.StatusOK)
}

func (a *App) Register(ctx *gin.Context) {
	var user dto.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	token, err := a.gophermartSvc.Register(user.Login, user.Password)
	if err != nil {
		if _, ok := err.(*models.UsedLoginError); ok {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Set("Token", token)
	ctx.Status(http.StatusOK)
}

func (a *App) Balance(ctx *gin.Context) {
	login := ctx.GetString("Login")
	user, err := a.gophermartSvc.GetUser(login)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, dto.Balance{Balance: float64(user.Balance), Withdrawn: user.Withdrawn})
}

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
	err = a.gophermartSvc.MakeWithdrawal(login, id, withdrawal.Sum)
	if err != nil {
		if _, ok := err.(*models.NotEnoughFunds); ok {
			ctx.AbortWithStatus(http.StatusPaymentRequired)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (a *App) Withdrawals(ctx *gin.Context) {
	login := ctx.GetString("Login")
	withdrawals, err := a.gophermartSvc.Withdrawals(login)
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

func (a *App) Orders(ctx *gin.Context) {
	login := ctx.GetString("Login")
	if login == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	orders, err := a.gophermartSvc.Orders(login)
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
	if err := a.gophermartSvc.AddOrder(login, orderID); err != nil {
		if _, ok := err.(*models.OrderExists); ok {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		if _, ok := err.(*models.OrderUsed); ok {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusAccepted)
}

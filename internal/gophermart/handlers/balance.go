package handlers

import (
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Balance(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login := ctx.GetString("Login")
		user, err := svc.GetUser(login)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, dto.Balance{Balance: float64(user.Balance), Withdrawn: user.Withdrawn})
	}
}

func MakeWithdrawal(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func Withdrawals(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login := ctx.GetString("Login")
		withdrawals, err := svc.Withdrawals(login)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, dto.NewWithdrawals(withdrawals))
	}
}

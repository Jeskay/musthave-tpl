package handlers

import (
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/dto"
	"musthave_tpl/internal/models"
	"net/http"
	"strconv"

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
		err = svc.MakeWithdrawal(login, id, withdrawal.Sum)
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
}

func Withdrawals(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login := ctx.GetString("Login")
		withdrawals, err := svc.Withdrawals(login)
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
}

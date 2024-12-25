package handlers

import (
	"io"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/dto"
	"musthave_tpl/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Orders(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login := ctx.GetString("Login")
		if login == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		orders, err := svc.Orders(login)
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
}

func PostOrder(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
		if err := svc.AddOrder(login, orderID); err != nil {
			if _, ok := err.(*gophermart.OrderExists); ok {
				ctx.AbortWithStatus(http.StatusOK)
				return
			}
			if _, ok := err.(*gophermart.OrderUsed); ok {
				ctx.AbortWithStatus(http.StatusConflict)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusAccepted)
	}
}

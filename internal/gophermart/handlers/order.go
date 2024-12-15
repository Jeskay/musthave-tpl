package handlers

import (
	"musthave_tpl/internal/gophermart"
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
		ctx.JSON(http.StatusOK, orders)
	}
}

func PostOrder(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		order := ctx.Param("order")
		login := ctx.GetString("Login")
		orderId, err := strconv.ParseInt(order, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := svc.AddOrder(login, orderId); err != nil {
			if _, ok := err.(*gophermart.OrderExists); ok {
				ctx.AbortWithStatus(http.StatusOK)
				return
			}
			if _, ok := err.(*gophermart.OrderUsed); ok {
				ctx.AbortWithStatus(http.StatusConflict)
				return
			}
			if _, ok := err.(*gophermart.InvalidOrderNumber); ok {
				ctx.AbortWithStatus(http.StatusUnprocessableEntity)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusCreated)
	}
}

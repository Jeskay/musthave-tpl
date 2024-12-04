package handlers

import (
	"musthave_tpl/internal"
	"musthave_tpl/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(svc *internal.LoyaltyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user dto.User
		if err := ctx.ShouldBindJSON(user); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		token, err := svc.Login(user.Login, user.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("Login", token)
		ctx.Status(http.StatusOK)
	}
}

func Register(svc *internal.LoyaltyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user dto.User
		if err := ctx.ShouldBindJSON(user); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err := svc.Register(user.Login, user.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

package handlers

import (
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user dto.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		token, err := svc.Login(user.Login, user.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("Token", token)
		ctx.Status(http.StatusOK)
	}
}

func Register(svc *gophermart.GophermartService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user dto.User
		if err := ctx.ShouldBindJSON(&user); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		token, err := svc.Register(user.Login, user.Password)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.Set("Token", token)
		ctx.Status(http.StatusOK)
	}
}

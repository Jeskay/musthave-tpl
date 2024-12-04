package middleware

import (
	"musthave_tpl/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authorize(svc *internal.LoyaltyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("auth_token")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		user, err := svc.Authenticate(token)
		if err != nil || user == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Set("Login", user.Login)
		ctx.Next()
	}
}

// TODO change cookie parameters to the relevant ones
func Authenticate(svc *internal.LoyaltyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if v, ok := ctx.Get("Token"); ok {
			token, ok := v.(string)
			if ok {
				ctx.SetCookie("auth_token", token, 3600, "/", "localhost", false, true)
			}
		}
	}
}

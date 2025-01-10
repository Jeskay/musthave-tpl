package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *MiddlewareService) Authorize(ctx *gin.Context) {
	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	login, err := m.authSvc.VerifyToken(token)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := m.dbSvc.UserByLogin(login)
	if err != nil || user == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set("Login", user.Login)
	ctx.Next()
}

func (m *MiddlewareService) Authenticate(ctx *gin.Context) {
	ctx.Next()
	if v, ok := ctx.Get("Token"); ok {
		token, ok := v.(string)
		if ok {
			ctx.SetCookie("auth_token", token, 3600, "/", "localhost", false, true)
		}
	}
}

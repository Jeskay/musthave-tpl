package middleware

import (
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/gophermart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MiddlewareService struct {
	authSvc       *auth.AuthService
	gophermartSvc *gophermart.GophermartService
}

func NewMiddleware(authSvc *auth.AuthService, gophermartSvc *gophermart.GophermartService) *MiddlewareService {
	return &MiddlewareService{
		authSvc:       authSvc,
		gophermartSvc: gophermartSvc,
	}
}

func (m *MiddlewareService) Authorize(ctx *gin.Context) {
	token, err := ctx.Cookie("auth_token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := m.gophermartSvc.Authenticate(token)
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

package middleware

import (
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/gophermart/db"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Authorize(ctx *gin.Context)
	Authenticate(ctx *gin.Context)
}

type MiddlewareService struct {
	authSvc auth.Authentication
	dbSvc   db.UserRepository
}

func NewMiddlewareService(authSvc auth.Authentication, dbSvc db.UserRepository) *MiddlewareService {
	return &MiddlewareService{
		authSvc: authSvc,
		dbSvc:   dbSvc,
	}
}

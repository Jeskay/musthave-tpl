package middleware

import (
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/gophermart/db"

	"github.com/gin-gonic/gin"
)

// Middleware provides methods for client authentication and authorization.
type Middleware interface {
	Authorize(ctx *gin.Context)
	Authenticate(ctx *gin.Context)
}

// MiddlewareService is an implementation of Middleware interface.
type MiddlewareService struct {
	authSvc auth.Authentication
	dbSvc   db.UserRepository
}

// NewMiddlewareService creates new instance of MiddlewareService.
func NewMiddlewareService(authSvc auth.Authentication, dbSvc db.UserRepository) *MiddlewareService {
	return &MiddlewareService{
		authSvc: authSvc,
		dbSvc:   dbSvc,
	}
}

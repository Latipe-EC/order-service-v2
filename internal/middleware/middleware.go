package middleware

import (
	"latipe-order-service-v2/internal/middleware/auth"
)

type Middleware struct {
	Authentication *auth.AuthenticationMiddleware
}

func NewMiddleware(auth *auth.AuthenticationMiddleware) *Middleware {
	return &Middleware{Authentication: auth}
}

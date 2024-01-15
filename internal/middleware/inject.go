package middleware

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/middleware/auth"
)

var Set = wire.NewSet(
	NewMiddleware,
	auth.NewAuthMiddleware,
)

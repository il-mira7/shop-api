package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"shopAPI/internal/shop/domain"

	"github.com/gin-gonic/gin"
)

type TokenValidator interface {
	ValidateToken(ctx context.Context, token string) (string, error)
}
type errorResponse struct {
	Message string `json:"error"`
}

const authTimeout = 5 * time.Second

func Auth(tv TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Message: "Authorization header is required",
			})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Message: "Authorization header must be Bearer {token}",
			})
			return
		}

		token := parts[1]

		ctx, cancel := context.WithTimeout(c.Request.Context(), authTimeout)
		defer cancel()

		userID, err := tv.ValidateToken(ctx, token)
		if err != nil {
			if errors.Is(err, domain.ErrAuthServiceUnavailable) {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, errorResponse{
					Message: domain.ErrAuthServiceUnavailable.Error(),
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{
				Message: "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

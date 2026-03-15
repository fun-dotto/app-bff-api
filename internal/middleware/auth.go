package middleware

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	userEmailKey contextKey = "email"
)

// AuthMiddleware returns a Gin middleware that verifies Firebase ID tokens
// and stores the user ID in the context.
func AuthMiddleware(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		idToken := strings.TrimPrefix(authorization, "Bearer ")

		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid ID token"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), userIDKey, token.UID)

		if email, ok := token.Claims["email"].(string); ok {
			ctx = context.WithValue(ctx, userEmailKey, email)
		}

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// UserIDFromContext extracts the user ID from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// UserEmailFromContext extracts the user email from context.
func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(userEmailKey).(string)
	return email, ok
}

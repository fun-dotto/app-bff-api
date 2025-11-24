package middleware

import (
	"net/http"

	"firebase.google.com/go/v4/appcheck"
	"github.com/gin-gonic/gin"
)

// AppCheckMiddleware returns a Gin middleware that verifies Firebase App Check tokens.
func AppCheckMiddleware(appCheckClient *appcheck.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		appCheckToken := c.GetHeader("X-Firebase-AppCheck")
		if appCheckToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		_, err := appCheckClient.VerifyToken(appCheckToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

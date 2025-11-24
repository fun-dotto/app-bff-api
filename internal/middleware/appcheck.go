package middleware

import (
	"net/http"
	"strings"

	"firebase.google.com/go/v4/appcheck"
	"github.com/gin-gonic/gin"
)

// AppCheckMiddleware returns a Gin middleware that verifies Firebase App Check tokens.
func AppCheckMiddleware(appCheckClient *appcheck.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Firebase-AppCheck")
		if token == "" {
			message := "X-Firebase-AppCheck header is required"
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
			return
		}

		_, err := appCheckClient.VerifyToken(strings.Replace(token, "Bearer ", "", 1))
		if err != nil {
			message := "X-Firebase-AppCheck header value is invalid"
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
			return
		}

		c.Next()
	}
}

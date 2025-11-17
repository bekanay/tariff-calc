package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleChecker interface {
	HasRole(accessToken, role string) (bool, error)
}

// RequireRole ensures the caller has the given role before allowing the request
func RequireRole(roleChecker RoleChecker, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "role is not configured"})
			return
		}

		accessToken := strings.TrimSpace(c.GetHeader("Authorization"))
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		hasRole, err := roleChecker.HasRole(accessToken, role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			return
		}

		c.Next()
	}
}

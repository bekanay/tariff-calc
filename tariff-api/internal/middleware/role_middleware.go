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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "role_not_configured",
				"message": "Роль не настроена на сервере",
			})
			return
		}

		accessToken := strings.TrimSpace(c.GetHeader("Authorization"))
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "missing_token",
				"message": "Требуется заголовок Authorization",
			})
			return
		}

		hasRole, err := roleChecker.HasRole(accessToken, role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   "role_check_failed",
				"message": "Не удалось проверить права доступа",
			})
			return
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "insufficient_role",
				"message": "Недостаточно прав",
			})
			return
		}

		c.Next()
	}
}

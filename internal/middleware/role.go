package middleware

import (
	"net/http"

	"github.com/riparuk/meet-book-api/internal/model"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has the required role
func RequireRole(requiredRole model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		if userRole.(model.UserRole) != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}

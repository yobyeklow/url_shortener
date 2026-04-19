package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequirePermission(allowRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error": "User role not found",
			})
		}

		role := userRole.(string)

		isAuthorized := slices.Contains(allowRoles, role)
		if !isAuthorized {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			return
		}
		ctx.Next()
	}
}
